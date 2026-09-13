import { del, get, post } from './http'
import type { Classroom, ClassroomSet } from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

export interface ClassroomForm {
    /** 空 = 新增（后端生成 UUID）；有值 = 更新。 */
    id: string
    building: string
    floor: number | undefined
    room: string
    capacity: number | undefined
    status: number
}

export async function listClassrooms(params: ListParams = {}): Promise<Classroom[]> {
    const set = await get<ClassroomSet>('/v1/classrooms', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.classrooms ?? []
}

/**
 * SaveClassroom。
 *
 * 注意这里**不需要** idField()：教室主键是字符串（"R101" / UUID），
 * 不存在 int64 精度问题，原样发字符串即可。
 */
export async function saveClassroom(form: ClassroomForm): Promise<void> {
    const body: Record<string, unknown> = {
        location: {
            building: form.building,
            floor: form.floor ?? 0,
            room: form.room,
        },
        capacity: form.capacity ?? 0,
        status: form.status,
    }
    if (form.id) body.id = form.id
    await post('/v1/classrooms', body)
}

export async function deleteClassroom(id: string): Promise<void> {
    await del(`/v1/classrooms/${encodeURIComponent(id)}`)
}
