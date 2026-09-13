import { del, get, idField, post } from './http'
import type { CourseEnrollment, CourseEnrollmentSet } from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

export interface EnrollForm {
    student_id: string
    course_id: string
}

export async function listEnrollments(params: ListParams = {}): Promise<CourseEnrollment[]> {
    const set = await get<CourseEnrollmentSet>('/v1/enrollments', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.enrollments ?? []
}

/**
 * EnrollStudent —— 只创建。服务端只读入 student_id 和 course_id，
 * 状态、时间都由领域层决定，所以表单只需要这两个字段。
 */
export async function enrollStudent(form: EnrollForm): Promise<void> {
    await post('/v1/enrollments', {
        student_id: idField(form.student_id),
        course_id: form.course_id,
    })
}

export async function deleteEnrollment(id: string): Promise<void> {
    await del(`/v1/enrollments/${encodeURIComponent(id)}`)
}
