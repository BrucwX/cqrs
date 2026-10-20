package classroom

import (
	"context"
	"testing"

	"courseScheduling/internal/domain/aggregate/classroom"
)

// intPtr 给可选字段造指针 —— 同包测试共用（updateClassroom_test.go 也用）。
func intPtr(v int) *int { return &v }

// TestCreateClassroom 新建：ID 由聚合生成，初始状态是可用。
func TestCreateClassroom(t *testing.T) {
	h, d := newHandler(t)

	loc, err := classroom.NewLocation("A", 1, "R101")
	if err != nil {
		t.Fatalf("new location: %v", err)
	}

	if err := h.CreateClassroom(context.Background(), CreateClassroomInput{
		Location: &loc,
		Capacity: intPtr(30),
	}); err != nil {
		t.Fatalf("CreateClassroom: %v", err)
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

// TestCreateClassroomInputRejected 必填项缺失、非法容量与非法状态都被聚合拒绝，且不落库。
func TestCreateClassroomInputRejected(t *testing.T) {
	loc, err := classroom.NewLocation("A", 1, "R101")
	if err != nil {
		t.Fatalf("new location: %v", err)
	}
	decommissioned := classroom.StatusDecommissioned

	cases := []struct {
		name string
		cmd  CreateClassroomInput
	}{
		{name: "缺位置", cmd: CreateClassroomInput{Capacity: intPtr(30)}},
		{name: "缺容量", cmd: CreateClassroomInput{Location: &loc}},
		{name: "容量非正", cmd: CreateClassroomInput{Location: &loc, Capacity: intPtr(0)}},
		{name: "不支持报废状态", cmd: CreateClassroomInput{Location: &loc, Capacity: intPtr(30), Status: &decommissioned}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// 每个用例一份干净存储，否则前面漏写的数据会污染后面的断言
			h, d := newHandler(t)

			if err := h.CreateClassroom(context.Background(), tc.cmd); err == nil {
				t.Error("期望被拒绝，实际成功")
			}
			if rooms := d.Classrooms(); len(rooms) != 0 {
				t.Errorf("被拒绝时不应写入，教室数量 = %d", len(rooms))
			}
		})
	}
}
