import { createRouter, createWebHistory } from 'vue-router'

/**
 * 只挂已经实现了的聚合。后端那 10 个 `GET /v1/xxx/{id}` 单条详情
 * 目前返回 501（Unimplemented），所以没有详情路由，列表就是全部数据源。
 *
 * 注意：`/v1/courses/available-for-*` 和 `/v1/qualifications/check`
 * 目前被 `/v1/courses/{id}` / `/v1/qualifications/{id}` 的单段路径遮挡
 * （HTTP 路由把静态段也喂给了 {id}），所以这些接口前端没有接。
 */
const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: '/', redirect: '/students' },
        {
            path: '/students',
            name: 'students',
            component: () => import('@/views/StudentsView.vue'),
        },
        {
            path: '/teachers',
            name: 'teachers',
            component: () => import('@/views/TeachersView.vue'),
        },
        {
            path: '/classrooms',
            name: 'classrooms',
            component: () => import('@/views/ClassroomsView.vue'),
        },
        {
            path: '/courses',
            name: 'courses',
            component: () => import('@/views/CoursesView.vue'),
        },
        {
            path: '/course-slots',
            name: 'course-slots',
            component: () => import('@/views/CourseSlotsView.vue'),
        },
        {
            path: '/course-slot-changes',
            name: 'course-slot-changes',
            component: () => import('@/views/CourseSlotChangesView.vue'),
        },
        {
            path: '/enrollments',
            name: 'enrollments',
            component: () => import('@/views/EnrollmentsView.vue'),
        },
        {
            path: '/absences',
            name: 'absences',
            component: () => import('@/views/AbsencesView.vue'),
        },
        {
            path: '/makeups',
            name: 'makeups',
            component: () => import('@/views/MakeupsView.vue'),
        },
        {
            path: '/qualifications',
            name: 'qualifications',
            component: () => import('@/views/QualificationsView.vue'),
        },
        { path: '/:pathMatch(.*)*', redirect: '/students' },
    ],
})

export default router
