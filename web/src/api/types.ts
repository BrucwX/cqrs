/**
 * 与后端 proto 一一对应的类型。
 *
 * 命名故意保持 snake_case，和实际 JSON 一致 —— 中间不再做一层 camelCase 转换，
 * 这样对着 proto 或 curl 输出排查问题时不会有认知负担。
 */

/** google.protobuf.Timestamp 被序列化成 { seconds, nanos }，不是 RFC3339。 */
export interface PbTimestamp {
    seconds?: number
    nanos?: number
}

/** 秒级时间戳 → 本地时间。后端没给 nanos，精确到分就够了。 */
export function formatTimestamp(ts?: PbTimestamp): string {
    if (!ts?.seconds) return '—'
    const d = new Date(ts.seconds * 1000)
    const p = (n: number) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

export interface ContactInfo {
    phone?: string
    email?: string
}

/** 枚举下拉/标签的渲染元数据。 */
export interface EnumOption {
    value: number
    label: string
    tag: 'primary' | 'success' | 'info' | 'warning' | 'danger'
}

/**
 * 取值口径直接抄 script/mysql/course_scheduling.sql 里的列注释，
 * 免得自己再起一套中文名。
 */
export function enumLabel(options: EnumOption[], value?: number): string {
    return options.find((o) => o.value === value)?.label ?? '—'
}

export function enumTag(
    options: EnumOption[],
    value?: number,
): 'primary' | 'success' | 'info' | 'warning' | 'danger' {
    return options.find((o) => o.value === value)?.tag ?? 'info'
}

// ---------------------------------------------------------------- Student

export interface Student {
    /** int64。API 层已转成字符串，避免超出 JS 安全整数范围时丢精度。 */
    id: string
    name: string
    /** 1 外部客户学员 / 2 内部员工学员。0 会被 omitempty 吃掉，所以是可选的。 */
    student_type?: number
    contact?: ContactInfo
    /** 1 正常 / 2 封禁 / 3 已注销 */
    status?: number
    created_at?: PbTimestamp
    updated_at?: PbTimestamp
}

export interface StudentSet {
    students?: Student[]
}

export const STUDENT_TYPES: EnumOption[] = [
    { value: 1, label: '外部客户学员', tag: 'info' },
    { value: 2, label: '内部员工学员', tag: 'primary' },
]

export const STUDENT_STATUSES: EnumOption[] = [
    { value: 1, label: '正常', tag: 'success' },
    { value: 2, label: '封禁', tag: 'danger' },
    { value: 3, label: '已注销', tag: 'info' },
]

// ---------------------------------------------------------------- Teacher

export interface Teacher {
    id: string
    /** 关联的学员系统 ID（讲师本人也是学员时）。 */
    student_id?: number
    name: string
    title?: string
    contact?: ContactInfo
    /** 1 在职 / 2 休假 / 3 已离职 */
    status?: number
    created_at?: PbTimestamp
    updated_at?: PbTimestamp
}

export interface TeacherSet {
    teachers?: Teacher[]
}

export const TEACHER_STATUSES: EnumOption[] = [
    { value: 1, label: '在职', tag: 'success' },
    { value: 2, label: '休假', tag: 'warning' },
    { value: 3, label: '已离职', tag: 'info' },
]

// ---------------------------------------------------------------- Classroom

export interface Location {
    building?: string
    floor?: number
    room?: string
}

export interface Classroom {
    /** 字符串主键（"R101" 或 UUID），不是 int64。 */
    id: string
    location?: Location
    capacity?: number
    /** 已分配座位；为 0 时后端不会返回这个字段。 */
    allocated?: number
    /** 1 可用 / 2 维护中 / 3 已报废 */
    status?: number
}

export interface ClassroomSet {
    classrooms?: Classroom[]
}

export const CLASSROOM_STATUSES: EnumOption[] = [
    { value: 1, label: '可用', tag: 'success' },
    { value: 2, label: '维护中', tag: 'warning' },
    { value: 3, label: '已报废', tag: 'info' },
]

/** 统一的手机号校验，跟领域层的规则保持一致。 */
export const PHONE_PATTERN = /^1[3-9]\d{9}$/

/** 秒级时间戳 → 本地 'YYYY-MM-DD'，给日期选择器回填用。 */
export function tsToDateStr(ts?: PbTimestamp): string {
    if (!ts?.seconds) return ''
    const d = new Date(ts.seconds * 1000)
    const p = (n: number) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}`
}

/** 'YYYY-MM-DD' → 本地当天 00:00 的秒，给日期选择器出站用。 */
export function dateToSeconds(dateStr: string): number {
    if (!dateStr) return 0
    return Math.floor(new Date(`${dateStr}T00:00:00`).getTime() / 1000)
}

// ---------------------------------------------------------------- Course

export interface Capacity {
    /** 最大人数。为 0 时后端不返回。 */
    max?: number
    /** 已报名人数。为 0 时后端不返回。 */
    enrolled?: number
}

export interface EnrollmentWindow {
    start_at?: PbTimestamp
    end_at?: PbTimestamp
    drop_deadline?: PbTimestamp
}

export interface CoursePeriod {
    start_at?: PbTimestamp
    end_at?: PbTimestamp
    total_hours?: number
    completed_hours?: number
}

export interface Course {
    /** 字符串主键（UUID）。 */
    id: string
    course_type_id?: string
    capacity?: Capacity
    enrollment?: EnrollmentWindow
    period?: CoursePeriod
}

export interface CourseSet {
    courses?: Course[]
}

// ---------------------------------------------------------------- CourseSlot

export interface DayTime {
    hour?: number
    minute?: number
}

export interface DayTimeRange {
    start?: DayTime
    end?: DayTime
}

export interface CourseSlot {
    id: string
    course_id?: string
    /** proto 1=周日 … 7=周六（domain 0-6 再 +1）。 */
    weekday?: number
    time_range?: DayTimeRange
    /** -1 = 未安排。 */
    teacher_id?: number | string
    /** 空 = 未安排。 */
    classroom_id?: string
    created_at?: PbTimestamp
    updated_at?: PbTimestamp
}

export interface CourseSlotSet {
    course_slots?: CourseSlot[]
}

/** 周几：proto 口径，1=周日 … 7=周六。 */
export const WEEKDAYS: EnumOption[] = [
    { value: 1, label: '周日', tag: 'info' },
    { value: 2, label: '周一', tag: 'info' },
    { value: 3, label: '周二', tag: 'info' },
    { value: 4, label: '周三', tag: 'info' },
    { value: 5, label: '周四', tag: 'info' },
    { value: 6, label: '周五', tag: 'info' },
    { value: 7, label: '周六', tag: 'info' },
]

/** 槽位时间范围 → "09:00–11:00"。minute 为 0 时后端不返回，按 0 处理。 */
export function formatTimeRange(range?: DayTimeRange): string {
    if (!range?.start) return '—'
    const fmt = (t?: DayTime): string => {
        if (!t) return '??:??'
        return `${String(t.hour ?? 0).padStart(2, '0')}:${String(t.minute ?? 0).padStart(2, '0')}`
    }
    return `${fmt(range.start)}–${fmt(range.end)}`
}

// ---------------------------------------------------------------- CourseSlotChange

export interface OriginalPlan {
    slot_id?: string
    date?: PbTimestamp
    teacher_id?: number | string
    classroom_id?: string
    /** "HH:MM"。 */
    start_time?: string
    end_time?: string
}

export interface TargetPlan {
    target_start_at?: PbTimestamp
    target_end_at?: PbTimestamp
    teacher_id?: number | string
    classroom_id?: string
}

export interface CourseSlotChange {
    /** int64，API 层已转字符串。 */
    id: string
    course_id?: string
    applicant_id?: number | string
    change_type?: number
    original_plan?: OriginalPlan
    target_plan?: TargetPlan
    reason?: string
    created_at?: PbTimestamp
    updated_at?: PbTimestamp
}

export interface CourseSlotChangeSet {
    course_slot_changes?: CourseSlotChange[]
}

export const CHANGE_TYPES: EnumOption[] = [
    { value: 1, label: '改期', tag: 'primary' },
    { value: 2, label: '代课', tag: 'warning' },
    { value: 3, label: '换教室', tag: 'success' },
    { value: 4, label: '复合变动', tag: 'danger' },
]

// ---------------------------------------------------------------- Enrollment

export interface CourseEnrollment {
    id: string
    student_id?: number | string
    course_id?: string
    /** 1 在读 / 2 已结业 / 3 已退课 / 4 未选课。 */
    status?: number
    enrolled_at?: PbTimestamp
    completed_at?: PbTimestamp
    dropped_at?: PbTimestamp
    updated_at?: PbTimestamp
}

export interface CourseEnrollmentSet {
    enrollments?: CourseEnrollment[]
}

export const ENROLLMENT_STATUSES: EnumOption[] = [
    { value: 1, label: '在读', tag: 'success' },
    { value: 2, label: '已结业', tag: 'info' },
    { value: 3, label: '已退课', tag: 'danger' },
    { value: 4, label: '未选课', tag: 'warning' },
]

// ---------------------------------------------------------------- Absence

export interface AbsenceRecord {
    id: string
    student_id?: number | string
    course_id?: string
    course_slot_id?: string
    schedule_date?: PbTimestamp
    missed_hours?: number
    /** 1 事假 / 2 公假 / 3 旷课。 */
    absence_type?: number
    reason?: string
    created_at?: PbTimestamp
    updated_at?: PbTimestamp
}

export interface AbsenceRecordSet {
    absences?: AbsenceRecord[]
}

export const ABSENCE_TYPES: EnumOption[] = [
    { value: 1, label: '事假', tag: 'info' },
    { value: 2, label: '公假', tag: 'primary' },
    { value: 3, label: '旷课', tag: 'danger' },
]

// ---------------------------------------------------------------- Makeup

export interface StudentMakeup {
    id: string
    student_id?: number | string
    course_id?: string
    original_slot_id?: string
    original_date?: PbTimestamp
    target_slot_id?: string
    target_date?: PbTimestamp
    makeup_hours?: number
    /** 1 已预约 / 2 已补课 / 3 已取消。 */
    status?: number
    completed_at?: PbTimestamp
    created_at?: PbTimestamp
    updated_at?: PbTimestamp
}

export interface StudentMakeupSet {
    makeups?: StudentMakeup[]
}

export const MAKEUP_STATUSES: EnumOption[] = [
    { value: 1, label: '已预约', tag: 'primary' },
    { value: 2, label: '已补课', tag: 'success' },
    { value: 3, label: '已取消', tag: 'info' },
]

// ---------------------------------------------------------------- Qualification

export interface Qualification {
    id: string
    teacher_id?: number | string
    course_type_id?: string
    certified_at?: PbTimestamp
    /** 1 生效 / 2 已吊销 / 3 已过期。 */
    status?: number
    updated_at?: PbTimestamp
}

export interface QualificationSet {
    qualifications?: Qualification[]
}

export const QUALIFICATION_STATUSES: EnumOption[] = [
    { value: 1, label: '生效', tag: 'success' },
    { value: 2, label: '已吊销', tag: 'danger' },
    { value: 3, label: '已过期', tag: 'info' },
]

// ---------------------------------------------------------------- CourseType（只读）

export interface CourseType {
    id: string
    name?: string
    description?: string
    created_at?: PbTimestamp
    updated_at?: PbTimestamp
}

export interface CourseTypeSet {
    course_types?: CourseType[]
}
