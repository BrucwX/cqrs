import { del, get, idField, post } from './http'
import type { CourseSlot, CourseSlotSet } from './types'

export interface ListParams {
    page?: number
    page_size?: number
}

/**
 * 课次（CourseSlot）没有 Save 接口，只有三个 Assign + Delete。
 * 三个 Assign 的 HTTP 规则都是 `body: "*"`，整个 body 直接绑定到请求消息。
 */

export async function listCourseSlots(params: ListParams = {}): Promise<CourseSlot[]> {
    const set = await get<CourseSlotSet>('/v1/course-slots', {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
    })
    return set.course_slots ?? []
}

export async function assignClassroom(slotIds: string[], classroomId: string): Promise<void> {
    await post('/v1/course-slots/assign-classroom', {
        slot_ids: slotIds,
        classroom_id: classroomId,
    })
}

export async function assignCourse(slotIds: string[], courseId: string): Promise<void> {
    await post('/v1/course-slots/assign-course', {
        slot_ids: slotIds,
        course_id: courseId,
    })
}

export async function assignTeacher(slotIds: string[], teacherId: string): Promise<void> {
    await post('/v1/course-slots/assign-teacher', {
        slot_ids: slotIds,
        teacher_id: idField(teacherId),
    })
}

export async function deleteCourseSlot(id: string): Promise<void> {
    await del(`/v1/course-slots/${encodeURIComponent(id)}`)
}
