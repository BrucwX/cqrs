-- =============================================================================
-- course_scheduling 上下文 · MySQL 建库建表脚本
--
-- 数据库名 = 限界上下文名：course_scheduling
-- 脚本自带建库与 USE，直接 `mysql -uroot -p < course_scheduling.sql` 即可。
-- 每次执行都会 DROP DATABASE IF EXISTS 再重建，结果是干净的一套表（数据不保留）。
--
-- 一个聚合根一张表，表名 = 聚合根类型名的 snake_case：
--   CourseType -> course_type          CourseEnrollment -> course_enrollment
--   Course     -> course               AbsenceRecord    -> absence_record
--   Classroom  -> classroom            StudentMakeup    -> student_makeup
--   CourseSlot -> course_slot          Qualification    -> qualification
--   CourseSlotChange -> course_slot_change
--   Student -> student                 Teacher -> teacher
--
-- 三条对应规则：
--   1. 值对象（Location / Capacity / EnrollmentWindow / CoursePeriod /
--      ContactInfo / DayTimeRange / OriginalPlan / TargetPlan）就地打平，
--      列名用「值对象名_字段名」，不再单开表。
--   2. 枚举（Status / AbsenceType / ChangeType / StudentType）都是 `iota + 1`，
--      存 tinyint unsigned，取值范围见列注释。
--   3. 主键就用聚合自己发的 ID：string 型是 uuid -> varchar(36)；
--      int64 型是 time.Now().UnixNano() -> bigint。--  4. 每张表都带 lock_version（bigint unsigned，默认 0）做乐观锁：
--     读时取出，写时 `SET lock_version = lock_version + 1 WHERE lock_version = ?`，
--     影响行数为 0 即版本冲突。它由存储层维护，不是领域字段。--
-- 两个「模型里就对不上」的地方，这里按 Go 类型如实建，先不改：
--   a. course_slot.id 是 varchar(36)，但 absence_record.course_slot_id、
--      student_makeup.original_slot_id / target_slot_id、
--      course_slot_change.original_slot_id 都是 int64，指不到它上面去。
--   b. course_slot 的时间是值对象（time），course_slot_change 的原计划快照
--      却是 "16:00" 这种字符串（varchar(5)）。
--
-- 不加外键（引用列只建普通索引）：ID 全由应用生成、写入顺序自由，
-- 删除也由应用决定，让 MySQL 少一份约束来源。
-- =============================================================================

-- 建库：先删后建，保证每次执行都是从零开始的一套干净表
-- ⚠️ 破坏性操作：库若已存在，其中所有数据会被删除
DROP DATABASE IF EXISTS course_scheduling;

CREATE DATABASE course_scheduling
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE course_scheduling;

SET NAMES utf8mb4;

-- -----------------------------------------------------------------------------
-- 一、基础档案
-- -----------------------------------------------------------------------------

-- 课程类型：描述「教什么」（少儿编程 / 成人英语），讲师资质绑在这一层
CREATE TABLE IF NOT EXISTS course_type (
  id          varchar(36)  NOT NULL                COMMENT '课程类型 ID（uuid）',
  name        varchar(64)  NOT NULL                COMMENT '类型名称',
  description varchar(255) NOT NULL DEFAULT ''     COMMENT '类型描述',
  created_at  datetime     NOT NULL                COMMENT '创建时间',
  updated_at  datetime     NOT NULL                COMMENT '更新时间',
  lock_version bigint unsigned NOT NULL DEFAULT 0  COMMENT '乐观锁版本号',
  PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '课程类型';

-- 教室：Location 打平为 building / floor / room
CREATE TABLE IF NOT EXISTS classroom (
  id        varchar(36)    NOT NULL              COMMENT '教室 ID（uuid）',
  building  varchar(64)    NOT NULL              COMMENT '楼栋',
  floor     int            NOT NULL              COMMENT '楼层',
  room      varchar(32)    NOT NULL              COMMENT '房间号',
  capacity  int            NOT NULL              COMMENT '总容量（座位数）',
  allocated int            NOT NULL DEFAULT 0    COMMENT '已分配座位数',
  status    tinyint unsigned NOT NULL            COMMENT '1 可用 / 2 维护中 / 3 已报废',
  lock_version bigint unsigned NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_classroom_location (building, floor, room)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '教室';

-- 学员：id 是全局统一的人员 ID，内部员工学员可能与讲师是同一个人
CREATE TABLE IF NOT EXISTS student (
  id           bigint           NOT NULL           COMMENT '学员 ID（人员/用户 ID）',
  name         varchar(64)      NOT NULL           COMMENT '姓名',
  student_type tinyint unsigned NOT NULL           COMMENT '1 外部客户学员 / 2 内部员工学员',
  phone        varchar(20)      NOT NULL DEFAULT '' COMMENT '手机号（ContactInfo）',
  email        varchar(128)     NOT NULL DEFAULT '' COMMENT '邮箱（ContactInfo，可空）',
  status       tinyint unsigned NOT NULL           COMMENT '1 正常 / 2 封禁 / 3 已注销',
  created_at   datetime         NOT NULL           COMMENT '创建时间',
  updated_at   datetime         NOT NULL           COMMENT '更新时间',
  lock_version bigint unsigned  NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_student_phone (phone)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '学员档案';

-- 讲师：student_id 是 TA 作为学员上课时的 ID，报名/缺勤记录挂在它上面
CREATE TABLE IF NOT EXISTS teacher (
  id         bigint           NOT NULL           COMMENT '讲师 ID',
  student_id bigint           NOT NULL           COMMENT '作为学员上课时的 ID（报名、缺勤记录用）',
  name       varchar(64)      NOT NULL           COMMENT '姓名',
  title      varchar(64)      NOT NULL DEFAULT '' COMMENT '职衔（金牌讲师 / 特级培训师）',
  phone      varchar(20)      NOT NULL DEFAULT '' COMMENT '手机号（ContactInfo）',
  email      varchar(128)     NOT NULL DEFAULT '' COMMENT '邮箱（ContactInfo，可空）',
  status     tinyint unsigned NOT NULL           COMMENT '1 在职 / 2 休假 / 3 已离职',
  created_at datetime         NOT NULL           COMMENT '创建时间',
  updated_at datetime         NOT NULL           COMMENT '更新时间',
  lock_version bigint unsigned NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_teacher_student (student_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '讲师';

-- -----------------------------------------------------------------------------
-- 二、教学安排
-- -----------------------------------------------------------------------------

-- 课程：Capacity / EnrollmentWindow / CoursePeriod 三个值对象打平
CREATE TABLE IF NOT EXISTS course (
  id                varchar(36) NOT NULL           COMMENT '课程 ID（uuid）',
  course_type_id    varchar(36) NOT NULL           COMMENT '所属课程类型（course_type.id）',
  capacity_max      int         NOT NULL           COMMENT '总容量',
  capacity_enrolled int         NOT NULL DEFAULT 0 COMMENT '已报名人数',
  enroll_start_at   datetime    NOT NULL           COMMENT '选课开始时间',
  enroll_end_at     datetime    NOT NULL           COMMENT '选课结束时间',
  drop_deadline     datetime    NOT NULL           COMMENT '退课截止时间',
  period_start_at   datetime    NOT NULL           COMMENT '教学周期开始',
  period_end_at     datetime    NOT NULL           COMMENT '教学周期结束',
  total_hours       int         NOT NULL           COMMENT '总课时',
  completed_hours   int         NOT NULL DEFAULT 0 COMMENT '已完成课时',
  lock_version      bigint unsigned NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_course_type (course_type_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '课程（具体的班）';

-- 课表模板：每周重复的排课槽位，teacher_id / classroom_id 支持「待分配」
CREATE TABLE IF NOT EXISTS course_slot (
  id           varchar(36)      NOT NULL            COMMENT '槽位 ID（uuid）',
  course_id    varchar(36)      NOT NULL            COMMENT '所属课程（course.id）',
  weekday      tinyint unsigned NOT NULL            COMMENT '星期几（0=周日 .. 6=周六）',
  start_time   time             NOT NULL            COMMENT '当天上课开始时间',
  end_time     time             NOT NULL            COMMENT '当天上课结束时间',
  teacher_id   bigint           NOT NULL DEFAULT -1 COMMENT '默认讲师 ID（-1 = 待分配）',
  classroom_id varchar(36)      NOT NULL DEFAULT '' COMMENT '默认教室 ID（空 = 待分配）',
  created_at   datetime         NOT NULL            COMMENT '创建时间',
  updated_at   datetime         NOT NULL            COMMENT '更新时间',
  lock_version bigint unsigned  NOT NULL DEFAULT 0  COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_slot_course (course_id),
  KEY idx_slot_teacher (teacher_id),
  KEY idx_slot_classroom (classroom_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '课表模板（每周重复）';

-- -----------------------------------------------------------------------------
-- 三、过程记录
-- -----------------------------------------------------------------------------

-- 临时换课：登记即生效，没有审批；原计划是一份快照，所以字段全打平
CREATE TABLE IF NOT EXISTS course_slot_change (
  id                    bigint           NOT NULL COMMENT '变更单 ID',
  course_id             varchar(36)      NOT NULL COMMENT '关联课程（course.id）',
  applicant_id          bigint           NOT NULL COMMENT '发起申请人 ID',
  change_type           tinyint unsigned NOT NULL COMMENT '1 改期 / 2 代课 / 3 换教室 / 4 复合变动',
  original_slot_id      bigint           NOT NULL COMMENT '原课表模板 ID（OriginalPlan 快照）',
  original_date         date             NOT NULL COMMENT '原定上课日期',
  original_teacher_id   bigint           NOT NULL COMMENT '原讲师 ID',
  original_classroom_id varchar(36)      NOT NULL COMMENT '原教室 ID',
  original_start_time   varchar(5)       NOT NULL COMMENT '原开始时间（HH:MM）',
  original_end_time     varchar(5)       NOT NULL COMMENT '原结束时间（HH:MM）',
  target_start_at       datetime         NOT NULL COMMENT '调整后开始时间',
  target_end_at         datetime         NOT NULL COMMENT '调整后结束时间（与开始同日）',
  target_teacher_id     bigint           NOT NULL COMMENT '实际授课讲师 ID',
  target_classroom_id   varchar(36)      NOT NULL COMMENT '实际使用教室 ID',
  reason                varchar(255)     NOT NULL COMMENT '调课/代课事由',
  created_at            datetime         NOT NULL COMMENT '创建时间',
  updated_at            datetime         NOT NULL COMMENT '更新时间',
  lock_version          bigint unsigned  NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_change_course (course_id),
  KEY idx_change_target_time (target_start_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '临时换课记录';

-- 选课报名：同一个学员可以重复报同一门课（退课后重选），所以不加唯一键
CREATE TABLE IF NOT EXISTS course_enrollment (
  id           bigint           NOT NULL           COMMENT '报名记录 ID',
  student_id   bigint           NOT NULL           COMMENT '学员 ID',
  course_id    varchar(36)      NOT NULL           COMMENT '课程 ID（course.id）',
  status       tinyint unsigned NOT NULL           COMMENT '1 在读 / 2 已结业 / 3 已退课',
  enrolled_at  datetime         NOT NULL           COMMENT '报名时间',
  completed_at datetime         NULL DEFAULT NULL  COMMENT '结业时间（未结业为 NULL）',
  dropped_at   datetime         NULL DEFAULT NULL  COMMENT '退课时间（未退课为 NULL）',
  updated_at   datetime         NOT NULL           COMMENT '更新时间',
  lock_version bigint unsigned  NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_enrollment_student (student_id),
  KEY idx_enrollment_course (course_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '选课报名';

-- 缺勤记录：一条事实，没有审批、也没有「已补卡」这类状态
CREATE TABLE IF NOT EXISTS absence_record (
  id             bigint           NOT NULL           COMMENT '缺勤记录 ID',
  student_id     bigint           NOT NULL           COMMENT '学员/员工 ID',
  course_id      varchar(36)      NOT NULL           COMMENT '关联课程 ID',
  course_slot_id bigint           NOT NULL           COMMENT '关联的排课槽位 ID',
  schedule_date  date             NOT NULL           COMMENT '具体上课日期',
  missed_hours   int              NOT NULL           COMMENT '缺席课时数',
  absence_type   tinyint unsigned NOT NULL           COMMENT '1 事假 / 2 公假 / 3 旷课',
  reason         varchar(255)     NOT NULL DEFAULT '' COMMENT '缺勤事由',
  created_at     datetime         NOT NULL           COMMENT '创建时间',
  updated_at     datetime         NOT NULL           COMMENT '更新时间',
  lock_version   bigint unsigned  NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_absence_student (student_id),
  KEY idx_absence_course (course_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '缺勤记录';

-- 补课预约：预约即生效，补课不跨课程
CREATE TABLE IF NOT EXISTS student_makeup (
  id               bigint           NOT NULL           COMMENT '补课记录 ID',
  student_id       bigint           NOT NULL           COMMENT '学员 ID',
  course_id        varchar(36)      NOT NULL           COMMENT '课程 ID（补课不跨课程）',
  original_slot_id bigint           NOT NULL           COMMENT '原排课槽位 ID',
  original_date    date             NOT NULL           COMMENT '原缺课日期',
  target_slot_id   bigint           NOT NULL           COMMENT '目标补课槽位 ID',
  target_date      date             NOT NULL           COMMENT '目标补课日期',
  makeup_hours     int              NOT NULL           COMMENT '补课课时数',
  status           tinyint unsigned NOT NULL           COMMENT '1 已预约 / 2 已补课 / 3 已取消',
  completed_at     datetime         NULL DEFAULT NULL  COMMENT '现场核销时间（未核销为 NULL）',
  created_at       datetime         NOT NULL           COMMENT '创建时间',
  updated_at       datetime         NOT NULL           COMMENT '更新时间',
  lock_version     bigint unsigned  NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_makeup_student (student_id),
  KEY idx_makeup_course (course_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '补课预约';

-- -----------------------------------------------------------------------------
-- 四、授课资质
-- -----------------------------------------------------------------------------

-- 资质绑「课程类型」不绑具体课程；能不能发由领域规则判定（讲师得先修完该类型的课）
CREATE TABLE IF NOT EXISTS qualification (
  id             bigint           NOT NULL COMMENT '资质 ID',
  teacher_id     bigint           NOT NULL COMMENT '讲师 ID',
  course_type_id varchar(36)      NOT NULL COMMENT '课程类型 ID（course_type.id）',
  certified_at   datetime         NOT NULL COMMENT '认证/试讲通过时间',
  status         tinyint unsigned NOT NULL COMMENT '1 生效 / 2 已吊销 / 3 已过期（模型已去掉有效期，当前不会出现）',
  updated_at     datetime         NOT NULL COMMENT '更新时间',
  lock_version   bigint unsigned  NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  PRIMARY KEY (id),
  KEY idx_qualification_teacher (teacher_id),
  KEY idx_qualification_course_type (course_type_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '授课资质';
