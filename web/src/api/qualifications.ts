import { del, get, idField, post } from './http'
import type { CourseType, CourseTypeSet, Qualification, QualificationSet } from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

export interface QualifyForm {
    teacher_id: string
    course_type_id: string
}

export async function listQualifications(
    params: ListParams = {},
): Promise<Qualification[]> {
    const set = await get<QualificationSet>('/v1/qualifications', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.qualifications ?? []
}

/** Qualify —— 只创建。服务端只读入 teacher_id + course_type_id。 */
export async function qualify(form: QualifyForm): Promise<void> {
    await post('/v1/qualifications', {
        teacher_id: idField(form.teacher_id),
        course_type_id: form.course_type_id,
    })
}

export async function deleteQualification(id: string): Promise<void> {
    await del(`/v1/qualifications/${encodeURIComponent(id)}`)
}

/** 某讲师有资质的课程类型（跨聚合：返回课程类型）。 */
export async function courseTypesByTeacher(teacherId: string): Promise<CourseType[]> {
    const set = await get<CourseTypeSet>(
        `/v1/teachers/${encodeURIComponent(teacherId)}/qualifications`,
    )
    return set.course_types ?? []
}
