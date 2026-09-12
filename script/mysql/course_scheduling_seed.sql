-- =============================================================================
-- course_scheduling 上下文 · MySQL 演示数据脚本
--
-- 与 internal/core/course_scheduling/adapters/memory/seed.go 的 SeedDemo 对齐：
-- 同一批 ID、同一批名字、同一批时间，方便内存实现与 MySQL 实现对照验证。
--
-- 用法：
--   mysql -uroot -p < script/mysql/course_scheduling.sql        # 先建库建表
--   mysql -uroot -p < script/mysql/course_scheduling_seed.sql   # 再灌演示数据
--
-- 可重复执行：脚本开头会清空这 11 张表再插入，
-- ⚠️ 表里原有的数据会被删掉（与建表脚本的先删后建保持同一种「重置」语义）。
-- lock_version 不显式写，统一由 DEFAULT 0 给出。
--
-- 时间全部钉死 2026 年秋季学期，保证每次执行结果完全一致（分页、冲突判定可复现）：
--   2026-09-07  学期开始 / 课程类型与资质的时间戳
--   2026-09-01 10:00:00  演示「当前时刻」（created_at / updated_at 用它）
--   2026-08-01 ~ 2026-09-06  选课窗口    2026-09-20  退课截止
--   2026-09-07 ~ 2027-01-15  教学周期
--
-- 期望行数：course_type 2 / classroom 2 / student 3 / teacher 3 / course 4 /
--           course_slot 6 / course_slot_change 2 / course_enrollment 4 /
--           absence_record 2 / student_makeup 2 / qualification 5  = 35 行
--
-- 刻意保留的两处「模型对不上」的数据（见建表脚本头部注释）：
--   a. absence_record.course_slot_id、student_makeup.original/target_slot_id、
--      course_slot_change.original_slot_id 存的是 1/2/3 这种 int64，
--      而 course_slot.id 是 uuid 字符串 —— 这里按 Go 模型如实灌。
--   b. course_slot_change 的原计划时间是 "09:00" 字符串，目标时间是 datetime。
--
-- 数据刻意包含时间冲突与不冲突两类课程，方便直接观察查询结果差异：
--   C001 周一/周三 09:00-11:00 张伟 R101   C003 周三/周五 09:00-11:00 王强 R102
--   → 周三 09:00-11:00 这段是热点：C001 与 C003 互相冲突。
-- =============================================================================

USE course_scheduling;

-- 演示「当前时刻」：与 seed.go 的 demoNow 一致
SET @demo_now = '2026-09-01 10:00:00';
-- 学期开始：课程类型 / 资质的 certify 时间戳
SET @term_start = '2026-09-07 00:00:00';

-- -----------------------------------------------------------------------------
-- 零、清空（无外键，顺序随意；可重复执行）
-- -----------------------------------------------------------------------------
DELETE FROM course_slot_change;
DELETE FROM student_makeup;
DELETE FROM absence_record;
DELETE FROM course_enrollment;
DELETE FROM qualification;
DELETE FROM course_slot;
DELETE FROM course;
DELETE FROM course_type;
DELETE FROM classroom;
DELETE FROM student;
DELETE FROM teacher;

-- -----------------------------------------------------------------------------
-- 一、基础档案
-- -----------------------------------------------------------------------------

-- 课程类型：ct-0001 少儿编程（C001/C002）/ ct-0002 成人英语（C003/C004）
INSERT INTO course_type (id, name, description, created_at, updated_at) VALUES
  ('ct-0001', '少儿编程', '面向 6-12 岁的图形化编程入门', @term_start, @term_start),
  ('ct-0002', '成人英语', '成人口语与商务英语',             @term_start, @term_start);

-- 教室：R101 30 座 / R102 20 座，都在 A 座 3 楼
INSERT INTO classroom (id, building, floor, room, capacity, allocated, status) VALUES
  ('R101', 'A座', 3, '301', 30, 0, 1),
  ('R102', 'A座', 3, '302', 20, 0, 1);

-- 学员：101 陈晨是内部员工学员（本身可能也是讲师），102/103 是外部客户学员
INSERT INTO student (id, name, student_type, phone, email, status, created_at, updated_at) VALUES
  (101, '陈晨', 2, '13900000101', 'chenchen@example.com', 1, @demo_now, @demo_now),
  (102, '刘洋', 1, '13900000102', 'liuyang@example.com',  1, @demo_now, @demo_now),
  (103, '赵敏', 1, '13900000103', 'zhaomin@example.com',  1, @demo_now, @demo_now);

-- 讲师：student_id 是 TA 作为学员上课时的 ID（2xx 段，与讲师 ID 区分开）
INSERT INTO teacher (id, student_id, name, title, phone, email, status, created_at, updated_at) VALUES
  (1, 201, '张伟', '金牌讲师',   '13800000001', 'zhangwei@example.com',  1, @demo_now, @demo_now),
  (2, 202, '李娜', '特级培训师', '13800000002', 'lina@example.com',      1, @demo_now, @demo_now),
  (3, 203, '王强', '讲师',       '13800000003', 'wangqiang@example.com', 1, @demo_now, @demo_now);

-- -----------------------------------------------------------------------------
-- 二、教学安排
-- -----------------------------------------------------------------------------

-- 课程：C001/C002 少儿编程，C003/C004 成人英语
-- completed_hours = capacity_enrolled * 4（与 seed.go 的 demoCourse 一致）
INSERT INTO course (
  id, course_type_id, capacity_max, capacity_enrolled,
  enroll_start_at, enroll_end_at, drop_deadline,
  period_start_at, period_end_at, total_hours, completed_hours
) VALUES
  ('C001', 'ct-0001', 30, 2, '2026-08-01 00:00:00', '2026-09-06 00:00:00', '2026-09-20 00:00:00',
                        '2026-09-07 00:00:00', '2027-01-15 00:00:00', 48, 8),
  ('C002', 'ct-0001', 20, 1, '2026-08-01 00:00:00', '2026-09-06 00:00:00', '2026-09-20 00:00:00',
                        '2026-09-07 00:00:00', '2027-01-15 00:00:00', 48, 4),
  ('C003', 'ct-0002', 15, 0, '2026-08-01 00:00:00', '2026-09-06 00:00:00', '2026-09-20 00:00:00',
                        '2026-09-07 00:00:00', '2027-01-15 00:00:00', 48, 0),
  ('C004', 'ct-0002', 12, 0, '2026-08-01 00:00:00', '2026-09-06 00:00:00', '2026-09-20 00:00:00',
                        '2026-09-07 00:00:00', '2027-01-15 00:00:00', 48, 0);

-- 周排期模板：weekday 1=周一 2=周二 3=周三 5=周五
INSERT INTO course_slot (
  id, course_id, weekday, start_time, end_time, teacher_id, classroom_id, created_at, updated_at
) VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'C001', 1, '09:00', '11:00', 1, 'R101', @demo_now, @demo_now),
  ('550e8400-e29b-41d4-a716-446655440002', 'C001', 3, '09:00', '11:00', 1, 'R101', @demo_now, @demo_now),
  ('550e8400-e29b-41d4-a716-446655440003', 'C002', 1, '14:00', '16:00', 2, 'R102', @demo_now, @demo_now),
  ('550e8400-e29b-41d4-a716-446655440004', 'C003', 3, '09:00', '11:00', 3, 'R102', @demo_now, @demo_now),
  ('550e8400-e29b-41d4-a716-446655440005', 'C003', 5, '09:00', '11:00', 3, 'R102', @demo_now, @demo_now),
  ('550e8400-e29b-41d4-a716-446655440006', 'C004', 2, '09:00', '11:00', 3, 'R102', @demo_now, @demo_now);

-- -----------------------------------------------------------------------------
-- 三、过程记录
-- -----------------------------------------------------------------------------

-- 临时换课：1 条改期（C001 周二下午）+ 1 条代课（C001 由李娜代张伟）
-- ⚠️ original_slot_id 是 int64 的 1/2，指不到 uuid 的 course_slot.id（见文件头注释 a）
INSERT INTO course_slot_change (
  id, course_id, applicant_id, change_type,
  original_slot_id, original_date, original_teacher_id, original_classroom_id,
  original_start_time, original_end_time,
  target_start_at, target_end_at, target_teacher_id, target_classroom_id,
  reason, created_at, updated_at
) VALUES
  (1, 'C001', 1, 1, 1, '2026-09-14', 1, 'R101', '09:00', '11:00',
     '2026-09-15 14:00:00', '2026-09-15 16:00:00', 1, 'R102',
     '场地检修，临时调至周二下午', @demo_now, @demo_now),
  (2, 'C001', 1, 2, 2, '2026-09-16', 1, 'R101', '09:00', '11:00',
     '2026-09-16 09:00:00', '2026-09-16 11:00:00', 2, 'R101',
     '讲师出差，由李娜代课', @demo_now, @demo_now);

-- 选课报名：103 的 C002 那条是「已结业」，用来验证查询侧不过滤状态
INSERT INTO course_enrollment (
  id, student_id, course_id, status, enrolled_at, completed_at, dropped_at, updated_at
) VALUES
  (1, 101, 'C001', 1, '2026-08-20 09:00:00', NULL,       NULL, @demo_now),
  (2, 102, 'C002', 1, '2026-08-20 09:00:00', NULL,       NULL, @demo_now),
  (3, 103, 'C001', 1, '2026-08-21 09:00:00', NULL,       NULL, @demo_now),
  (4, 103, 'C002', 2, '2026-08-21 09:00:00', @demo_now,   NULL, @demo_now);

-- 缺勤记录：1 条事假 + 1 条旷课（absence_type 1 事假 / 2 公假 / 3 旷课）
-- ⚠️ course_slot_id 存 int64 的 2/3，同样指不到 uuid 的 course_slot.id
INSERT INTO absence_record (
  id, student_id, course_id, course_slot_id, schedule_date,
  missed_hours, absence_type, reason, created_at, updated_at
) VALUES
  (1, 101, 'C001', 2, '2026-09-16', 2, 1, '家中急事',   @demo_now, @demo_now),
  (2, 102, 'C002', 3, '2026-09-14', 2, 3, '未请假缺席', @demo_now, @demo_now);

-- 补课预约：1 条已补课（现场核销）+ 1 条已预约（status 1 已预约 / 2 已补课 / 3 已取消）
INSERT INTO student_makeup (
  id, student_id, course_id,
  original_slot_id, original_date, target_slot_id, target_date,
  makeup_hours, status, completed_at, created_at, updated_at
) VALUES
  (1, 101, 'C001', 1, '2026-09-14', 2, '2026-09-16', 2, 2, @demo_now, @demo_now, @demo_now),
  (2, 102, 'C002', 3, '2026-09-14', 3, '2026-09-21', 2, 1, NULL,      @demo_now, @demo_now);

-- -----------------------------------------------------------------------------
-- 四、授课资质
-- -----------------------------------------------------------------------------

-- 资质绑课程类型不绑具体课程；5 号「王强 · 少儿编程」已吊销，
-- 用来验证查询侧不过滤状态（status 1 生效 / 2 已吊销 / 3 已过期）
INSERT INTO qualification (
  id, teacher_id, course_type_id, certified_at, status, updated_at
) VALUES
  (1, 1, 'ct-0001', @term_start, 1, @demo_now),  -- 张伟 · 少儿编程
  (2, 2, 'ct-0001', @term_start, 1, @demo_now),  -- 李娜 · 少儿编程
  (3, 3, 'ct-0002', @term_start, 1, @demo_now),  -- 王强 · 成人英语
  (4, 1, 'ct-0002', @term_start, 1, @demo_now),  -- 张伟 · 成人英语
  (5, 3, 'ct-0001', @term_start, 2, @demo_now);  -- 王强 · 少儿编程（已吊销）

-- -----------------------------------------------------------------------------
-- 五、自检：打印每张表的行数与期望值，跑完一眼能看出有没有漏
-- -----------------------------------------------------------------------------
SELECT 'course_type'        AS table_name, COUNT(*) AS rows_now, 2 AS rows_expected FROM course_type
UNION ALL SELECT 'classroom',          COUNT(*), 2 FROM classroom
UNION ALL SELECT 'student',            COUNT(*), 3 FROM student
UNION ALL SELECT 'teacher',            COUNT(*), 3 FROM teacher
UNION ALL SELECT 'course',             COUNT(*), 4 FROM course
UNION ALL SELECT 'course_slot',        COUNT(*), 6 FROM course_slot
UNION ALL SELECT 'course_slot_change', COUNT(*), 2 FROM course_slot_change
UNION ALL SELECT 'course_enrollment',  COUNT(*), 4 FROM course_enrollment
UNION ALL SELECT 'absence_record',     COUNT(*), 2 FROM absence_record
UNION ALL SELECT 'student_makeup',     COUNT(*), 2 FROM student_makeup
UNION ALL SELECT 'qualification',      COUNT(*), 5 FROM qualification;
