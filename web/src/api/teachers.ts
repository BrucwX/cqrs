import { del, get, idField, post } from './http'
import type { Teacher, TeacherSet } from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

export interface TeacherForm {
    /** 空 = 新增；有值 = 更新。 */
    id?: string
    student_id: string
    name: string
    title: string
    phone: string
    email: string
    status: number
}

export async function listTeachers(params: ListParams = {}): Promise<Teacher[]> {
    const set = await get<TeacherSet>('/v1/teachers', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.teachers ?? []
}

/** SaveTeacher —— `body: "teacher"`，直接发讲师对象。 */
export async function saveTeacher(form: TeacherForm): Promise<void> {
    const body: Record<string, unknown> = {
        name: form.name,
        title: form.title,
        contact: { phone: form.phone, email: form.email },
        status: form.status,
    }
    if (form.id) body.id = idField(form.id)
    // student_id == 0 表示没关联，后端会跳过（service 层判的是 != 0）
    const sid = Number(form.student_id)
    if (Number.isInteger(sid) && sid > 0) body.student_id = sid
    await post('/v1/teachers', body)
}

export async function deleteTeacher(id: string): Promise<void> {
    await del(`/v1/teachers/${encodeURIComponent(id)}`)
}
