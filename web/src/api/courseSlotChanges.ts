import { del, get, idField, post } from './http'
import { dateToSeconds, type CourseSlotChange, type CourseSlotChangeSet } from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

export interface ChangeForm {
    course_id: string
    applicant_id: string
    change_type: number
    reason: string
    original_slot_id: string
    original_date: string
    original_teacher_id: string
    original_classroom_id: string
    original_start_time: string
    original_end_time: string
    /** 'YYYY-MM-DDTHH:mm'，本地时间。 */
    target_start_at: string
    target_end_at: string
    target_teacher_id: string
    target_classroom_id: string
}

function dtToSeconds(value: string): number {
    if (!value) return 0
    return Math.floor(new Date(value).getTime() / 1000)
}

export async function listCourseSlotChanges(
    params: ListParams = {},
): Promise<CourseSlotChange[]> {
    const set = await get<CourseSlotChangeSet>('/v1/course-slot-changes', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.course_slot_changes ?? []
}

/**
 * ChangeCourseSlot —— 只创建（id 由服务端生成，入参里的 id 会被忽略）。
 * body: "course_slot_change" → 直接发对象。
 */
export async function changeCourseSlot(form: ChangeForm): Promise<void> {
    const body: Record<string, unknown> = {
        course_id: form.course_id,
        applicant_id: idField(form.applicant_id),
        change_type: form.change_type,
        reason: form.reason,
        original_plan: {
            slot_id: form.original_slot_id,
            date: { seconds: dateToSeconds(form.original_date) },
            teacher_id: form.original_teacher_id ? idField(form.original_teacher_id) : 0,
            classroom_id: form.original_classroom_id,
            start_time: form.original_start_time,
            end_time: form.original_end_time,
        },
        target_plan: {
            target_start_at: { seconds: dtToSeconds(form.target_start_at) },
            target_end_at: { seconds: dtToSeconds(form.target_end_at) },
            teacher_id: idField(form.target_teacher_id),
            classroom_id: form.target_classroom_id,
        },
    }
    await post('/v1/course-slot-changes', body)
}

export async function deleteCourseSlotChange(id: string): Promise<void> {
    await del(`/v1/course-slot-changes/${encodeURIComponent(id)}`)
}
