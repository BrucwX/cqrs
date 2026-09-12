package student

import (
	"context"
	"errors"
	"testing"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

func strPtr(v string) *string { return &v }

func int64Ptr(v int64) *int64 { return &v }

func newContact(t *testing.T) student.ContactInfo {
	t.Helper()

	contact, err := student.NewContactInfo("13800000000", "")
	if err != nil {
		t.Fatalf("new contact: %v", err)
	}
	return contact
}

// TestSaveStudentCreates ID 为 nil 时新建，默认正常状态。
func TestSaveStudentCreates(t *testing.T) {
	h, d := newHandler(t)

	contact := newContact(t)
	external := student.TypeExternal
	if err := h.SaveStudent(context.Background(), StudentInput{
		Name:        strPtr("陈晨"),
		StudentType: &external,
		Contact:     &contact,
	}); err != nil {
		t.Fatalf("SaveStudent: %v", err)
	}

	students := d.Students()
	if len(students) != 1 {
		t.Fatalf("学员数量 = %d, want 1", len(students))
	}
	got := students[0]
	if got.ID() == 0 {
		t.Error("新建学员应当自动分配 ID")
	}
	if got.Name() != "陈晨" || got.StudentType() != student.TypeExternal {
		t.Errorf("student = (%q, %v), want (陈晨, %v)", got.Name(), got.StudentType(), student.TypeExternal)
	}
	if got.Status() != student.StatusActive {
		t.Errorf("status = %v, want %v", got.Status(), student.StatusActive)
	}
}

// TestSaveStudentUpdatesContact 更新联系方式，其他字段保持原值。
func TestSaveStudentUpdatesContact(t *testing.T) {
	h, d := newHandler(t)

	external := student.TypeExternal
	if err := h.SaveStudent(context.Background(), StudentInput{
		Name:        strPtr("陈晨"),
		StudentType: &external,
	}); err != nil {
		t.Fatalf("SaveStudent: %v", err)
	}
	id := d.Students()[0].ID()

	contact := newContact(t)
	if err := h.SaveStudent(context.Background(), StudentInput{ID: &id, Contact: &contact}); err != nil {
		t.Fatalf("SaveStudent: %v", err)
	}

	got, _ := d.StudentByID(id)
	if got.Contact().Phone() != "13800000000" {
		t.Errorf("phone = %q, want 13800000000", got.Contact().Phone())
	}
	if got.Name() != "陈晨" {
		t.Errorf("name = %q, want 陈晨（未给出的字段应保持原值）", got.Name())
	}
}

// TestSaveStudentChangesStatus 封禁 / 解封 / 注销。
func TestSaveStudentChangesStatus(t *testing.T) {
	h, d := newHandler(t)

	if err := h.SaveStudent(context.Background(), StudentInput{Name: strPtr("陈晨")}); err != nil {
		t.Fatalf("SaveStudent: %v", err)
	}
	id := d.Students()[0].ID()

	cases := []struct {
		name string
		set  student.Status
	}{
		{name: "封禁", set: student.StatusBlocked},
		{name: "解封", set: student.StatusActive},
		{name: "注销", set: student.StatusCancelled},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			set := tc.set
			if err := h.SaveStudent(context.Background(), StudentInput{ID: &id, Status: &set}); err != nil {
				t.Fatalf("SaveStudent: %v", err)
			}
			if got, _ := d.StudentByID(id); got.Status() != tc.set {
				t.Errorf("status = %v, want %v", got.Status(), tc.set)
			}
		})
	}
}

// TestSaveStudentUpdateMissing 更新不存在的学员报 not found。
func TestSaveStudentUpdateMissing(t *testing.T) {
	h, d := newHandler(t)

	err := h.SaveStudent(context.Background(), StudentInput{ID: int64Ptr(999)})
	if !errors.Is(err, student.ErrStudentNotFound) {
		t.Errorf("err = %v, want %v", err, student.ErrStudentNotFound)
	}
	if students := d.Students(); len(students) != 0 {
		t.Errorf("被拒绝时不应新建，学员数量 = %d", len(students))
	}
}

// TestSaveStudentRejected 新建缺名字 / 非法状态被拒绝。
func TestSaveStudentRejected(t *testing.T) {
	h, _ := newHandler(t)
	unknown := student.Status(99)

	if err := h.SaveStudent(context.Background(), StudentInput{}); err == nil {
		t.Error("缺名字应当被拒绝")
	}
	if err := h.SaveStudent(context.Background(), StudentInput{
		Name:   strPtr("陈晨"),
		Status: &unknown,
	}); !errors.Is(err, student.ErrUnknownStatus) {
		t.Errorf("err = %v, want %v", err, student.ErrUnknownStatus)
	}
}
