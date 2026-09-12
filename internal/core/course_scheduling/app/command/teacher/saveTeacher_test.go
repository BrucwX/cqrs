package teacher

import (
	"context"
	"errors"
	"testing"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

func strPtr(v string) *string { return &v }

func int64Ptr(v int64) *int64 { return &v }

func newContact(t *testing.T) teacher.ContactInfo {
	t.Helper()

	contact, err := teacher.NewContactInfo("13800000000", "")
	if err != nil {
		t.Fatalf("new contact: %v", err)
	}
	return contact
}

// TestSaveTeacherCreates ID 为 nil 时新建，默认在职。
func TestSaveTeacherCreates(t *testing.T) {
	h, d := newHandler(t)

	contact := newContact(t)
	if err := h.SaveTeacher(context.Background(), TeacherInput{
		Name:    strPtr("李娜"),
		Title:   strPtr("讲师"),
		Contact: &contact,
	}); err != nil {
		t.Fatalf("SaveTeacher: %v", err)
	}

	teachers := d.Teachers()
	if len(teachers) != 1 {
		t.Fatalf("讲师数量 = %d, want 1", len(teachers))
	}
	got := teachers[0]
	if got.ID() == 0 {
		t.Error("新建讲师应当自动分配 ID")
	}
	if got.Name() != "李娜" || got.Title() != "讲师" {
		t.Errorf("teacher = (%q, %q), want (李娜, 讲师)", got.Name(), got.Title())
	}
	if !got.IsActive() {
		t.Errorf("status = %v, want %v", got.Status(), teacher.StatusActive)
	}
	if got.StudentID() == 0 {
		t.Error("没指定学员 ID 时，聚合应当自己发一个")
	}
}

// TestSaveTeacherCreatesWithStudentID 新增时可以指定「作为学员上课」的 ID。
func TestSaveTeacherCreatesWithStudentID(t *testing.T) {
	h, d := newHandler(t)

	contact := newContact(t)
	if err := h.SaveTeacher(context.Background(), TeacherInput{
		StudentID: int64Ptr(203),
		Name:      strPtr("李娜"),
		Contact:   &contact,
	}); err != nil {
		t.Fatalf("SaveTeacher: %v", err)
	}

	if got := d.Teachers()[0].StudentID(); got != 203 {
		t.Errorf("studentID = %d, want 203", got)
	}
}

// TestSaveTeacherKeepsStudentID 学员 ID 是身份，更新时改不动。
func TestSaveTeacherKeepsStudentID(t *testing.T) {
	h, d := newHandler(t)

	contact := newContact(t)
	if err := h.SaveTeacher(context.Background(), TeacherInput{
		StudentID: int64Ptr(203),
		Name:      strPtr("李娜"),
		Contact:   &contact,
	}); err != nil {
		t.Fatalf("SaveTeacher: %v", err)
	}
	id := d.Teachers()[0].ID()

	if err := h.SaveTeacher(context.Background(), TeacherInput{
		ID:        &id,
		StudentID: int64Ptr(999),
		Name:      strPtr("李娜"),
		Contact:   &contact,
	}); err != nil {
		t.Fatalf("SaveTeacher: %v", err)
	}

	if got := d.Teachers()[0].StudentID(); got != 203 {
		t.Errorf("studentID = %d, want 203（更新不该动它）", got)
	}
}

// TestSaveTeacherCreatesOnLeave 新建时可以直接指定休假状态。
func TestSaveTeacherCreatesOnLeave(t *testing.T) {
	h, d := newHandler(t)

	contact := newContact(t)
	onLeave := teacher.StatusOnLeave
	if err := h.SaveTeacher(context.Background(), TeacherInput{
		Name:    strPtr("李娜"),
		Contact: &contact,
		Status:  &onLeave,
	}); err != nil {
		t.Fatalf("SaveTeacher: %v", err)
	}

	got := d.Teachers()[0]
	if !got.IsOnLeave() {
		t.Errorf("status = %v, want %v", got.Status(), teacher.StatusOnLeave)
	}
}

// TestSaveTeacherUpdatesPartially 更新时只改命令里给出的字段。
func TestSaveTeacherUpdatesPartially(t *testing.T) {
	h, d := newHandler(t)

	contact := newContact(t)
	if err := h.SaveTeacher(context.Background(), TeacherInput{
		Name:    strPtr("李娜"),
		Title:   strPtr("讲师"),
		Contact: &contact,
	}); err != nil {
		t.Fatalf("SaveTeacher: %v", err)
	}
	id := d.Teachers()[0].ID()

	newTitle := "特级培训师"
	if err := h.SaveTeacher(context.Background(), TeacherInput{
		ID:    &id,
		Title: &newTitle,
	}); err != nil {
		t.Fatalf("SaveTeacher: %v", err)
	}

	got, _ := d.TeacherByID(id)
	if got.Title() != newTitle {
		t.Errorf("title = %q, want %q", got.Title(), newTitle)
	}
	if got.Name() != "李娜" {
		t.Errorf("name = %q, want 李娜（未给出的字段应保持原值）", got.Name())
	}
}

// TestSaveTeacherChangesStatus 状态切换（含幂等）。
func TestSaveTeacherChangesStatus(t *testing.T) {
	h, d := newHandler(t)

	contact := newContact(t)
	if err := h.SaveTeacher(context.Background(), TeacherInput{
		Name:    strPtr("李娜"),
		Contact: &contact,
	}); err != nil {
		t.Fatalf("SaveTeacher: %v", err)
	}
	id := d.Teachers()[0].ID()

	onLeave := teacher.StatusOnLeave
	if err := h.SaveTeacher(context.Background(), TeacherInput{ID: &id, Status: &onLeave}); err != nil {
		t.Fatalf("SaveTeacher(onLeave): %v", err)
	}
	if got, _ := d.TeacherByID(id); !got.IsOnLeave() {
		t.Errorf("status = %v, want %v", got.Status(), teacher.StatusOnLeave)
	}

	// 重复设成同一个状态应当幂等，而不是报 ErrAlreadyOnLeave
	if err := h.SaveTeacher(context.Background(), TeacherInput{ID: &id, Status: &onLeave}); err != nil {
		t.Errorf("重复设置同一状态应当幂等，实际 %v", err)
	}

	active := teacher.StatusActive
	if err := h.SaveTeacher(context.Background(), TeacherInput{ID: &id, Status: &active}); err != nil {
		t.Fatalf("SaveTeacher(active): %v", err)
	}
	if got, _ := d.TeacherByID(id); !got.IsActive() {
		t.Errorf("status = %v, want %v", got.Status(), teacher.StatusActive)
	}
}

// TestSaveTeacherUpdateMissing 更新不存在的讲师报 not found。
func TestSaveTeacherUpdateMissing(t *testing.T) {
	h, d := newHandler(t)

	err := h.SaveTeacher(context.Background(), TeacherInput{ID: int64Ptr(999)})
	if !errors.Is(err, repo.ErrTeacherNotFound) {
		t.Errorf("err = %v, want %v", err, repo.ErrTeacherNotFound)
	}
	if teachers := d.Teachers(); len(teachers) != 0 {
		t.Errorf("被拒绝时不应新建，讲师数量 = %d", len(teachers))
	}
}

// TestSaveTeacherRejected 新建缺名字 / 非法状态被拒绝。
func TestSaveTeacherRejected(t *testing.T) {
	h, _ := newHandler(t)
	unknown := teacher.Status(99)

	if err := h.SaveTeacher(context.Background(), TeacherInput{}); err == nil {
		t.Error("缺名字应当被拒绝")
	}
	if err := h.SaveTeacher(context.Background(), TeacherInput{
		Name:   strPtr("李娜"),
		Status: &unknown,
	}); !errors.Is(err, teacher.ErrUnknownStatus) {
		t.Errorf("err = %v, want %v", err, teacher.ErrUnknownStatus)
	}
}
