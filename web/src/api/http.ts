import axios, { type AxiosError, type AxiosResponse } from 'axios'

/**
 * 后端不是 protojson，而是 encoding/json 直接打在 pb 结构体上（实测），所以：
 *   - 字段名是 snake_case（student_type / created_at），不是 camelCase
 *   - 枚举是数字（status: 1），不是 "STUDENT_STATUS_ACTIVE"
 *   - Timestamp 是 { seconds } 而不是 RFC3339 字符串
 *   - 生成 struct 带 omitempty，零值字段直接不出现在响应里
 *
 * 最要命的一条：Student / Teacher 的 id 由服务端 time.Now().UnixNano() 生成，
 * 形如 1789298836897654925 —— 超过 Number.MAX_SAFE_INTEGER(9007199254740991)，
 * JSON.parse 会静默把它四舍五入成 1789298836897655000，
 * 再拿这个 id 去 DELETE 就删不到东西，而且不报错。
 *
 * 所以下面在 JSON.parse 之前把大整数包成字符串，出站时再还原成裸数字。
 */

/** 任意字段下 ≥16 位的裸整数。16 位以内都还在安全区，不用管。 */
const READ_BIG_INT = /("(?:[^"\\]|\\.)*"\s*:\s*)(-?\d{16,})(?=\s*[,}])/g
/** id 一律读成字符串，免得出现「同一列有时 number 有时 string」。 */
const READ_ID = /("id"\s*:\s*)(-?\d+)(?=\s*[,}])/g
/** 出站：把 idField() 打过标记的超大 id 还原成裸数字。 */
const WRITE_BIG_INT = /"__BIGINT__(-?\d+)"/g

const BIG_INT_PREFIX = '__BIGINT__'

function quoteBigInts(raw: string): string {
    return raw.replace(READ_BIG_INT, '$1"$2"').replace(READ_ID, '$1"$2"')
}

function unquoteBigInts(raw: string): string {
    return raw.replace(WRITE_BIG_INT, '$1')
}

/**
 * 把界面上的字符串 id 转成请求体里该有的形态。
 *
 * 小 id 直接给数字；超出安全区间的加前缀，由 transformRequest 还原成裸数字
 * （后端 encoding/json 不接受把字符串塞进 int64 字段）。
 * 非纯数字 id（Classroom 的 "R101" / UUID）原样输出。
 */
export function idField(id: string): number | string {
    if (!/^-?\d+$/.test(id)) return id
    const n = Number(id)
    return Number.isSafeInteger(n) ? n : BIG_INT_PREFIX + id
}

/** 归一化后的接口错误：界面只需要拿到 message 去弹提示。 */
export class ApiError extends Error {
    readonly status?: number
    readonly code?: number

    constructor(message: string, status?: number, code?: number) {
        super(message)
        this.name = 'ApiError'
        this.status = status
        this.code = code
    }
}

interface ErrorBody {
    code?: number
    message?: string
}

const http = axios.create({
    baseURL: '/',
    timeout: 15000,
    headers: { 'Content-Type': 'application/json' },
    transformResponse: [
        (data: unknown) => {
            if (typeof data !== 'string' || data.length === 0) return data
            try {
                return JSON.parse(quoteBigInts(data))
            } catch {
                // 不是 JSON（比如网关返回的 HTML 错误页），原样交回
                return data
            }
        },
    ],
    transformRequest: [
        (data: unknown) => {
            if (data === undefined || data === null) return data
            // 自己 stringify 是为了能把 __BIGINT__ 标记还原成裸数字
            return unquoteBigInts(JSON.stringify(data))
        },
    ],
})

http.interceptors.response.use(
    (res: AxiosResponse) => res,
    (error: AxiosError<ErrorBody>) => {
        const body = error.response?.data
        let message = body?.message
        if (!message) {
            if (error.code === 'ECONNABORTED') {
                message = '请求超时'
            } else if (!error.response) {
                message = '连不上后端服务，确认 :8000 已经起来了'
            } else {
                message = error.message || '请求失败'
            }
        }
        // 注意：后端把领域校验失败也返回成 HTTP 500，别按状态码分支，直接显示 message
        return Promise.reject(new ApiError(message, error.response?.status, body?.code))
    },
)

export async function get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
    const res = await http.get<T>(url, { params })
    return res.data
}

export async function post<T>(url: string, body: unknown): Promise<T> {
    const res = await http.post<T>(url, body)
    return res.data
}

export async function del<T>(url: string): Promise<T> {
    const res = await http.delete<T>(url)
    return res.data
}
