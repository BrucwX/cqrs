import { del, get, post } from './http'
import {
    dateToSeconds,
    type AbsenceRecord,
    type AbsenceRecordSet,
    type Course,
    type CourseSet,
    type CourseSlot,
    type CourseSlotChange,
    type CourseSlotChangeSet,
    type CourseSlotSet,
    type Student,
    type StudentMakeup,
    type StudentMakeupSet,
    type StudentSet,
} from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

export interface CourseForm {
    id?: string
    course_type_id: string
    capacity_max: number | undefined
    capacity_enrolled: number | undefined
    enroll_start_at: string
    enroll_end_at: string
    drop_deadline: string
    period_start_at: string
    period_end_at: string
    total_hours: number | undefined
    completed_hours: number | undefined
}

export async function listCourses(params: ListParams = {}): Promise<Course[]> {
    const set = await get<CourseSet>('/v1/courses', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.courses ?? []
}

/** 新增（id 空）或更新（带 id）。课程主键是 UUID 字符串，不走 idField。 */
export async function saveCourse(form: CourseForm): Promise<void> {
    const body: Record<string, unknown> = {
        course_type_id: form.course_type_id,
        capacity: {
            max: form.capacity_max ?? 0,
            enrolled: form.capacity_enrolled ?? 0,
        },
        enrollment: {
            start_at: { seconds: dateToSeconds(form.enroll_start_at) },
            end_at: { seconds: dateToSeconds(form.enroll_end_at) },
            drop_deadline: { seconds: dateToSeconds(form.drop_deadline) },
        },
        period: {
            start_at: { seconds: dateToSeconds(form.period_start_at) },
            end_at: { seconds: dateToSeconds(form.period_end_at) },
            total_hours: form.total_hours ?? 0,
            completed_hours: form.completed_hours ?? 0,
        },
    }
    if (form.id) body.id = form.id
    await post('/v1/courses', body)
}

export async function deleteCourse(id: string): Promise<void> {
    await del(`/v1/courses/${encodeURIComponent(id)}`)
}

// ---------------------------------------------------------- 课程关联查询
// 两段路径的接口没被 GetCourse 的单段路径挡住，都可用。

export async function courseStudents(courseId: string): Promise<Student[]> {
    const set = await get<StudentSet>(`/v1/courses/${encodeURIComponent(courseId)}/students`)
    return set.students ?? []
}

export async function courseSlots(courseId: string): Promise<CourseSlot[]> {
    const set = await get<CourseSlotSet>(`/v1/courses/${encodeURIComponent(courseId)}/slots`)
    return set.course_slots ?? []
}

export async function courseAbsences(courseId: string): Promise<AbsenceRecord[]> {
    const set = await get<AbsenceRecordSet>(`/v1/courses/${encodeURIComponent(courseId)}/absences`)
    return set.absences ?? []
}

export async function courseMakeups(courseId: string): Promise<StudentMakeup[]> {
    const set = await get<StudentMakeupSet>(`/v1/courses/${encodeURIComponent(courseId)}/makeups`)
    return set.makeups ?? []
}

export async function courseSlotChanges(courseId: string): Promise<CourseSlotChange[]> {
    const set = await get<CourseSlotChangeSet>(
        `/v1/courses/${encodeURIComponent(courseId)}/slot-changes`,
    )
    return set.course_slot_changes ?? []
}
