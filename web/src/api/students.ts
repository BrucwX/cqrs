import { del, get, idField, post } from './http'
import type { Student, StudentSet } from './types'

/** 后端分页参数。默认页大小 20，无上限。 */
export interface ListParams {
    page?: number
    page_size?: number
}

/** 表单模型 —— 界面上的形状，和请求体不完全一样（id 是字符串）。 */
export interface StudentForm {
    /** 空 = 新增；有值 = 更新。（proto 里那句 "Set id to create" 的注释是反的，以实现为准） */
    id?: string
    name: string
    student_type: number
    phone: string
    email: string
    status: number
}

/**
 * ListStudents。
 *
 * 注意后端不返回 total，只返回当前页，所以调用方没法算总页数。
 */
export async function listStudents(params: ListParams = {}): Promise<Student[]> {
    const set = await get<StudentSet>('/v1/students', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    // 空列表时后端返回 {}，字段直接不存在
    return set.students ?? []
}

/**
 * SaveStudent —— 新增和更新都打这一个接口（后端是 upsert）。
 *
 * HTTP 规则是 `body: "student"`，意思是整个 body 绑定到 student 字段，
 * 所以要直接发 Student 对象，不能包一层 {"student": {...}}。
 */
export async function saveStudent(form: StudentForm): Promise<void> {
    const body: Record<string, unknown> = {
        name: form.name,
        student_type: form.student_type,
        contact: { phone: form.phone, email: form.email },
        status: form.status,
    }
    if (form.id) body.id = idField(form.id)
    await post('/v1/students', body)
}

export async function deleteStudent(id: string): Promise<void> {
    await del(`/v1/students/${encodeURIComponent(id)}`)
}
