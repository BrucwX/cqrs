import { del, get, idField, post } from './http'
import { dateToSeconds, type AbsenceRecord, type AbsenceRecordSet } from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

export interface AbsenceForm {
    student_id: string
    course_id: string
    course_slot_id: string
    schedule_date: string
    missed_hours: number | undefined
    absence_type: number
    reason: string
}

export async function listAbsences(params: ListParams = {}): Promise<AbsenceRecord[]> {
    const set = await get<AbsenceRecordSet>('/v1/absences', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.absences ?? []
}

/** RecordAbsence —— 只创建。 */
export async function recordAbsence(form: AbsenceForm): Promise<void> {
    await post('/v1/absences', {
        student_id: idField(form.student_id),
        course_id: form.course_id,
        course_slot_id: form.course_slot_id,
        schedule_date: { seconds: dateToSeconds(form.schedule_date) },
        missed_hours: form.missed_hours ?? 0,
        absence_type: form.absence_type,
        reason: form.reason,
    })
}

export async function deleteAbsence(id: string): Promise<void> {
    await del(`/v1/absences/${encodeURIComponent(id)}`)
}
