package classroom

import (
	"context"
	"errors"
	"testing"

	"courseScheduling/internal/domain/aggregate/classroom"
)

// TestUpdateClassroomPartially 更新时只改命令里给出的字段。
func TestUpdateClassroomPartially(t *testing.T) {
	h, d := newHandler(t)
	d.SeedClassroom(newTestClassroom(t, "R101"))

	// 只给容量：位置保持原值
	if err := h.UpdateClassroom(context.Background(), UpdateClassroomInput{
		ID:       "R101",
		Capacity: intPtr(50),
	}); err != nil {
		t.Fatalf("UpdateClassroom: %v", err)
	}

	got, _ := d.ClassroomByID("R101")
	if got.Capacity() != 50 {
		t.Errorf("capacity = %d, want 50", got.Capacity())
	}
	if got.Location().Room() != "R101" {
		t.Errorf("location room = %q, want R101（未给出的字段应保持原值）", got.Location().Room())
	}
}

// TestUpdateClassroomLocationAndStatus 位置与维护状态都能改。
func TestUpdateClassroomLocationAndStatus(t *testing.T) {
	h, d := newHandler(t)
	d.SeedClassroom(newTestClassroom(t, "R101"))

	loc, err := classroom.NewLocation("B", 2, "B201")
	if err != nil {
		t.Fatalf("new location: %v", err)
	}
	underMaintenance := classroom.StatusUnderMaintenance

	if err := h.UpdateClassroom(context.Background(), UpdateClassroomInput{
		ID:       "R101",
		Location: &loc,
		Status:   &underMaintenance,
	}); err != nil {
		t.Fatalf("UpdateClassroom: %v", err)
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
	if err := h.UpdateClassroom(context.Background(), UpdateClassroomInput{
		ID:     "R101",
		Status: &available,
	}); err != nil {
		t.Fatalf("UpdateClassroom: %v", err)
	}
	got, _ = d.ClassroomByID("R101")
	if got.Status() != classroom.StatusAvailable {
		t.Errorf("status = %v, want %v", got.Status(), classroom.StatusAvailable)
	}
}

// TestUpdateClassroomMissing 更新一间不存在的教室报 not found，不会顺手新建。
func TestUpdateClassroomMissing(t *testing.T) {
	h, d := newHandler(t)

	err := h.UpdateClassroom(context.Background(), UpdateClassroomInput{
		ID:       "R999",
		Capacity: intPtr(30),
	})
	if !errors.Is(err, classroom.ErrClassroomNotFound) {
		t.Errorf("err = %v, want %v", err, classroom.ErrClassroomNotFound)
	}
	if rooms := d.Classrooms(); len(rooms) != 0 {
		t.Errorf("被拒绝时不应新建，教室数量 = %d", len(rooms))
	}
}

// TestUpdateClassroomRejected 更新已存在教室时，非法值同样被聚合挡下，且不改动原值。
func TestUpdateClassroomRejected(t *testing.T) {
	h, d := newHandler(t)
	d.SeedClassroom(newTestClassroom(t, "R101"))

	if err := h.UpdateClassroom(context.Background(), UpdateClassroomInput{
		ID:       "R101",
		Capacity: intPtr(0),
	}); !errors.Is(err, classroom.ErrInvalidCapacity) {
		t.Errorf("err = %v, want %v", err, classroom.ErrInvalidCapacity)
	}

	got, _ := d.ClassroomByID("R101")
	if got.Capacity() != 30 {
		t.Errorf("被拒绝时不应改动，capacity = %d, want 30", got.Capacity())
	}
}
