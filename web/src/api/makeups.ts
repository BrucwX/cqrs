import { del, get, idField, post } from './http'
import { dateToSeconds, type StudentMakeup, type StudentMakeupSet } from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

export interface MakeupForm {
    student_id: string
    course_id: string
    original_slot_id: string
    original_date: string
    target_slot_id: string
    target_date: string
    makeup_hours: number | undefined
}

export async function listMakeups(params: ListParams = {}): Promise<StudentMakeup[]> {
    const set = await get<StudentMakeupSet>('/v1/makeups', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.makeups ?? []
}

/** RecordMakeup —— 只创建。状态由领域层决定（新记录默认「已预约」）。 */
export async function recordMakeup(form: MakeupForm): Promise<void> {
    await post('/v1/makeups', {
        student_id: idField(form.student_id),
        course_id: form.course_id,
        original_slot_id: form.original_slot_id,
        original_date: { seconds: dateToSeconds(form.original_date) },
        target_slot_id: form.target_slot_id,
        target_date: { seconds: dateToSeconds(form.target_date) },
        makeup_hours: form.makeup_hours ?? 0,
    })
}

export async function deleteMakeup(id: string): Promise<void> {
    await del(`/v1/makeups/${encodeURIComponent(id)}`)
}
