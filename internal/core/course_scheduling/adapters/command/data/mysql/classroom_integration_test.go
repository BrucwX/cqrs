package mysql

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"cqrs/internal/conf"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// 命令侧集成测试连一个独立库，与读侧集成测试（course_scheduling 种子库）
// 隔开 —— 两边在 go test ./... 里会并行跑，共用同一个库会让对方看到测试中途
// 的临时行。本测试自建 uuid 教室，不依赖种子数据；连不上就跳过。
const cmdTestDB = "course_scheduling_cmd_test"

func cmdTestDSN() string {
	if dsn := os.Getenv("COURSE_SCHEDULING_CMD_TEST_DSN"); dsn != "" {
		return dsn
	}
	return "root@unix(/var/run/mysqld/mysqld.sock)/" + cmdTestDB
}

// bootstrapCmdTestDB 建好独立的测试库和 classroom 表（幂等）。
func bootstrapCmdTestDB(t *testing.T) {
	t.Helper()

	db, err := sql.Open("mysql", "root@unix(/var/run/mysqld/mysqld.sock)/")
	if err != nil {
		t.Skipf("跳过集成测试：连不上 MySQL（%v）", err)
	}
	defer func() { _ = db.Close() }()

	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + cmdTestDB + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		t.Skipf("跳过集成测试：无法建测试库（%v）", err)
	}

	// 按主库（course_scheduling，由建表脚本建好并灌入种子）的结构复制空表。
	tables := []string{
		"course_type", "classroom", "student", "teacher", "course",
		"course_slot", "course_slot_change", "course_enrollment",
		"absence_record", "student_makeup", "qualification",
	}
	for _, tbl := range tables {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS " + cmdTestDB + "." + tbl + " LIKE course_scheduling." + tbl); err != nil {
			t.Skipf("跳过集成测试：无法建表 %s（%v）", tbl, err)
		}
	}
}

func newCommandData(t *testing.T) *MysqlData {
	t.Helper()

	bootstrapCmdTestDB(t)

	data, cleanup, err := NewMysqlData(&conf.Data{
		Database: &conf.Data_Database{Driver: "mysql", Source: cmdTestDSN()},
	})
	if err != nil {
		t.Skipf("跳过集成测试：连不上 MySQL（%v）", err)
	}
	t.Cleanup(cleanup)
	return data
}

func newTestClassroom(t *testing.T) *classroom.Classroom {
	t.Helper()
	loc, err := classroom.NewLocation("测试楼", 9, "901")
	if err != nil {
		t.Fatal(err)
	}
	cap := 30
	cl, err := classroom.NewClassroom(&loc, &cap, nil)
	if err != nil {
		t.Fatal(err)
	}
	return cl
}

// createAndTrack 建一间 uuid 教室并注册清理，测试结束时删掉，避免污染共享库。
func createAndTrack(t *testing.T, c repo.ClassroomCommand) *classroom.Classroom {
	t.Helper()
	cl := newTestClassroom(t)
	if err := c.CreateClassroom(context.Background(), cl); err != nil {
		t.Fatalf("Create: %v", err)
	}
	t.Cleanup(func() {
		if err := c.DeleteClassroom(context.Background(), cl.ID()); err != nil && !errors.Is(err, classroom.ErrClassroomNotFound) {
			t.Errorf("cleanup delete: %v", err)
		}
	})
	return cl
}

// TestClassroomImpCreateAndGet 新增后能原样读回。
func TestClassroomImpCreateAndGet(t *testing.T) {
	c := newCommandData(t)
	cl := createAndTrack(t, c)

	got, err := c.MustGetClassroom(context.Background(), cl.ID())
	if err != nil {
		t.Fatalf("MustGet: %v", err)
	}
	if got.ID() != cl.ID() {
		t.Errorf("ID = %q, want %q", got.ID(), cl.ID())
	}
	if got.Location().FullName() != cl.Location().FullName() {
		t.Errorf("Location = %q, want %q", got.Location().FullName(), cl.Location().FullName())
	}
	if got.Capacity() != 30 {
		t.Errorf("Capacity = %d, want 30", got.Capacity())
	}
	if got.Status() != classroom.StatusAvailable {
		t.Errorf("Status = %v, want available", got.Status())
	}
}

// TestClassroomImpUpdate 更新后字段生效。
func TestClassroomImpUpdate(t *testing.T) {
	c := newCommandData(t)
	cl := createAndTrack(t, c)

	newCap := 25
	err := c.UpdateClassroom(context.Background(), cl.ID(), func(_ context.Context, cur *classroom.Classroom) (*classroom.Classroom, error) {
		if err := cur.Update(nil, &newCap, nil); err != nil {
			return nil, err
		}
		return cur, nil
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := c.MustGetClassroom(context.Background(), cl.ID())
	if err != nil {
		t.Fatalf("MustGet: %v", err)
	}
	if got.Capacity() != 25 {
		t.Errorf("Capacity = %d, want 25", got.Capacity())
	}
}

// TestClassroomImpDelete 删除后取不到，重复删报 NotFound。
func TestClassroomImpDelete(t *testing.T) {
	c := newCommandData(t)
	cl := createAndTrack(t, c)

	if err := c.DeleteClassroom(context.Background(), cl.ID()); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := c.MustGetClassroom(context.Background(), cl.ID()); !errors.Is(err, classroom.ErrClassroomNotFound) {
		t.Fatalf("MustGet(after delete) = %v, want ErrClassroomNotFound", err)
	}
	if err := c.DeleteClassroom(context.Background(), cl.ID()); !errors.Is(err, classroom.ErrClassroomNotFound) {
		t.Fatalf("Delete(twice) = %v, want ErrClassroomNotFound", err)
	}
}

// TestClassroomImpTransaction 事务提交可见、回滚不可见。
func TestClassroomImpTransaction(t *testing.T) {
	data := newCommandData(t)
	c := data
	cl := newTestClassroom(t)
	t.Cleanup(func() {
		_ = c.DeleteClassroom(context.Background(), cl.ID())
	})

	// 回滚：End 收到非 nil err，Create 的写入应被撤销。
	ctx, err := data.Begin(context.Background())
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if err := c.CreateClassroom(ctx, cl); err != nil {
		t.Fatalf("Create(rollback): %v", err)
	}
	if err := data.End(ctx, errors.New("boom")); err == nil {
		t.Fatal("End should surface the injected error")
	}
	if _, err := c.MustGetClassroom(context.Background(), cl.ID()); !errors.Is(err, classroom.ErrClassroomNotFound) {
		t.Fatalf("MustGet(after rollback) = %v, want ErrClassroomNotFound", err)
	}

	// 提交：End 收到 nil err，Create 的写入应可见。
	ctx, err = data.Begin(context.Background())
	if err != nil {
		t.Fatalf("Begin(commit): %v", err)
	}
	if err := c.CreateClassroom(ctx, cl); err != nil {
		t.Fatalf("Create(commit): %v", err)
	}
	if err := data.End(ctx, nil); err != nil {
		t.Fatalf("End(commit): %v", err)
	}
	if _, err := c.MustGetClassroom(context.Background(), cl.ID()); err != nil {
		t.Fatalf("MustGet(after commit) = %v, want nil", err)
	}
}
