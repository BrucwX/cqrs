package classroom

import (
	"context"
	"errors"
	"testing"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

func intPtr(v int) *int { return &v }

func strPtr(v string) *string { return &v }

// TestSaveClassroomCreates ID 为 nil 时新建，ID 由服务端生成。
func TestSaveClassroomCreates(t *testing.T) {
	h, d := newHandler(t)

	loc, err := classroom.NewLocation("A", 1, "R101")
	if err != nil {
		t.Fatalf("new location: %v", err)
	}

	if err := h.SaveClassroom(context.Background(), ClassroomInput{
		Location: &loc,
		Capacity: intPtr(30),
	}); err != nil {
		t.Fatalf("SaveClassroom: %v", err)
	}

	rooms := d.Classrooms()
	if len(rooms) != 1 {
		t.Fatalf("教室数量 = %d, want 1", len(rooms))
	}
	got := rooms[0]
	if got.ID() == "" {
		t.Error("新建教室应当自动分配 ID")
	}
	if got.Capacity() != 30 {
		t.Errorf("capacity = %d, want 30", got.Capacity())
	}
	if got.Status() != classroom.StatusAvailable {
		t.Errorf("status = %v, want %v", got.Status(), classroom.StatusAvailable)
	}
}

// TestSaveClassroomUpdatesPartially 更新时只改命令里给出的字段。
func TestSaveClassroomUpdatesPartially(t *testing.T) {
	h, d := newHandler(t)
	d.SeedClassroom(newTestClassroom(t, "R101"))

	// 只给容量：位置保持原值
	if err := h.SaveClassroom(context.Background(), ClassroomInput{
		ID:       strPtr("R101"),
		Capacity: intPtr(50),
	}); err != nil {
		t.Fatalf("SaveClassroom: %v", err)
	}

	got, _ := d.ClassroomByID("R101")
	if got.Capacity() != 50 {
		t.Errorf("capacity = %d, want 50", got.Capacity())
	}
	if got.Location().Room() != "R101" {
		t.Errorf("location room = %q, want R101（未给出的字段应保持原值）", got.Location().Room())
	}
}

// TestSaveClassroomUpdatesLocationAndStatus 位置与维护状态都能改。
func TestSaveClassroomUpdatesLocationAndStatus(t *testing.T) {
	h, d := newHandler(t)
	d.SeedClassroom(newTestClassroom(t, "R101"))

	loc, err := classroom.NewLocation("B", 2, "B201")
	if err != nil {
		t.Fatalf("new location: %v", err)
	}
	underMaintenance := classroom.StatusUnderMaintenance

	if err := h.SaveClassroom(context.Background(), ClassroomInput{
		ID:       strPtr("R101"),
		Location: &loc,
		Status:   &underMaintenance,
	}); err != nil {
		t.Fatalf("SaveClassroom: %v", err)
	}

	got, _ := d.ClassroomByID("R101")
	if got.Location().Building() != "B" || got.Location().Room() != "B201" {
		t.Errorf("location = %v, want B-2F-B201", got.Location().FullName())
	}
	if got.Status() != classroom.StatusUnderMaintenance {
		t.Errorf("status = %v, want %v", got.Status(), classroom.StatusUnderMaintenance)
	}

	// 再切回可用
	available := classroom.StatusAvailable
	if err := h.SaveClassroom(context.Background(), ClassroomInput{
		ID:     strPtr("R101"),
		Status: &available,
	}); err != nil {
		t.Fatalf("SaveClassroom: %v", err)
	}
	got, _ = d.ClassroomByID("R101")
	if got.Status() != classroom.StatusAvailable {
		t.Errorf("status = %v, want %v", got.Status(), classroom.StatusAvailable)
	}
}

// TestSaveClassroomUpdateMissing 更新一间不存在的教室报 not found，不会顺手新建。
func TestSaveClassroomUpdateMissing(t *testing.T) {
	h, d := newHandler(t)

	err := h.SaveClassroom(context.Background(), ClassroomInput{
		ID:       strPtr("R999"),
		Capacity: intPtr(30),
	})
	if !errors.Is(err, repo.ErrClassroomNotFound) {
		t.Errorf("err = %v, want %v", err, repo.ErrClassroomNotFound)
	}
	if rooms := d.Classrooms(); len(rooms) != 0 {
		t.Errorf("被拒绝时不应新建，教室数量 = %d", len(rooms))
	}
}

// TestSaveClassroomInputRejected 新建必填项与非法状态都被拒绝。
func TestSaveClassroomInputRejected(t *testing.T) {
	loc, err := classroom.NewLocation("A", 1, "R101")
	if err != nil {
		t.Fatalf("new location: %v", err)
	}
	decommissioned := classroom.StatusDecommissioned

	cases := []struct {
		name string
		cmd  ClassroomInput
	}{
		{name: "新建缺位置", cmd: ClassroomInput{Capacity: intPtr(30)}},
		{name: "新建缺容量", cmd: ClassroomInput{Location: &loc}},
		{name: "新建容量非正", cmd: ClassroomInput{Location: &loc, Capacity: intPtr(0)}},
		{name: "新建不支持报废状态", cmd: ClassroomInput{Location: &loc, Capacity: intPtr(30), Status: &decommissioned}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// 每个用例一份干净存储，否则前面漏写的数据会污染后面的断言
			h, d := newHandler(t)

			if err := h.SaveClassroom(context.Background(), tc.cmd); err == nil {
				t.Error("期望被拒绝，实际成功")
			}
			if rooms := d.Classrooms(); len(rooms) != 0 {
				t.Errorf("被拒绝时不应写入，教室数量 = %d", len(rooms))
			}
		})
	}
}

// TestSaveClassroomUpdateRejected 更新已存在教室时，非法值同样被聚合挡下。
func TestSaveClassroomUpdateRejected(t *testing.T) {
	h, d := newHandler(t)
	d.SeedClassroom(newTestClassroom(t, "R101"))

	if err := h.SaveClassroom(context.Background(), ClassroomInput{
		ID:       strPtr("R101"),
		Capacity: intPtr(0),
	}); !errors.Is(err, classroom.ErrInvalidCapacity) {
		t.Errorf("err = %v, want %v", err, classroom.ErrInvalidCapacity)
	}

	got, _ := d.ClassroomByID("R101")
	if got.Capacity() != 30 {
		t.Errorf("被拒绝时不应改动，capacity = %d, want 30", got.Capacity())
	}
}
