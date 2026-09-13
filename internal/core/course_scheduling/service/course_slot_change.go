package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	courseslotchangev1 "cqrs/api/v1/course_scheduling/course_slot_change"
	courseslotchangecmd "cqrs/internal/core/course_scheduling/app/command/courseSlotChange"
	courseslotchangeqry "cqrs/internal/core/course_scheduling/app/query/courseSlotChange"
	courseslotchangedo "cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// CourseSlotChangeService 临时换课服务。
type CourseSlotChangeService struct {
	courseslotchangev1.UnimplementedCourseSlotChangeServiceServer

	cmd *courseslotchangecmd.Handler
	qry *courseslotchangeqry.Handler
}

// NewCourseSlotChangeService 创建临时换课服务。
func NewCourseSlotChangeService(cmd *courseslotchangecmd.Handler, qry *courseslotchangeqry.Handler) *CourseSlotChangeService {
	return &CourseSlotChangeService{cmd: cmd, qry: qry}
}

// ChangeCourseSlot 发起临时换课。
func (s *CourseSlotChangeService) ChangeCourseSlot(ctx context.Context, req *courseslotchangev1.ChangeCourseSlotRequest) (*emptypb.Empty, error) {
	c := req.CourseSlotChange
	if c == nil {
		return nil, courseslotchangedo.ErrSlotChangeRequired
	}
	in := courseslotchangecmd.ChangeCourseSlot{
		CourseID:    c.CourseId,
		ApplicantID: c.ApplicantId,
		ChangeType:  courseslotchangedo.ChangeType(c.ChangeType),
		Reason:      c.Reason,
	}
	if c.OriginalPlan != nil {
		in.Original = courseslotchangedo.NewOriginalPlan(
			c.OriginalPlan.SlotId,
			tsTime(c.OriginalPlan.Date),
			c.OriginalPlan.TeacherId,
			c.OriginalPlan.ClassroomId,
			c.OriginalPlan.StartTime,
			c.OriginalPlan.EndTime,
		)
	}
	if c.TargetPlan != nil {
		target, err := courseslotchangedo.NewTargetPlan(
			tsTime(c.TargetPlan.TargetStartAt),
			tsTime(c.TargetPlan.TargetEndAt),
			c.TargetPlan.TeacherId,
			c.TargetPlan.ClassroomId,
		)
		if err != nil {
			return nil, err
		}
		in.Target = target
	}
	if _, err := s.cmd.ChangeCourseSlot(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteCourseSlotChange 删除换课记录。
func (s *CourseSlotChangeService) DeleteCourseSlotChange(ctx context.Context, req *courseslotchangev1.DeleteCourseSlotChangeRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteCourseSlotChange(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListCourseSlotChanges 分页查询换课记录。
func (s *CourseSlotChangeService) ListCourseSlotChanges(ctx context.Context, req *courseslotchangev1.ListCourseSlotChangesRequest) (*courseslotchangev1.CourseSlotChangeSet, error) {
	items, err := s.qry.PageCourseSlotChanges(ctx, courseslotchangeqry.PageCourseSlotChanges{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	return newCourseSlotChangeSet(items), nil
}

// ListCourseSlotChangesByCourseID 某课程的全部换课记录。
func (s *CourseSlotChangeService) ListCourseSlotChangesByCourseID(ctx context.Context, req *courseslotchangev1.ListCourseSlotChangesByCourseIDRequest) (*courseslotchangev1.CourseSlotChangeSet, error) {
	items, err := s.qry.ListByCourseID(ctx, courseslotchangeqry.ListByCourseID{CourseID: req.CourseId})
	if err != nil {
		return nil, err
	}
	return newCourseSlotChangeSet(items), nil
}

// newCourseSlotChangeSet []*CourseSlotChange → CourseSlotChangeSet。
func newCourseSlotChangeSet(items []*courseslotchangedo.CourseSlotChange) *courseslotchangev1.CourseSlotChangeSet {
	set := &courseslotchangev1.CourseSlotChangeSet{CourseSlotChanges: make([]*courseslotchangev1.CourseSlotChange, 0, len(items))}
	for _, do := range items {
		set.CourseSlotChanges = append(set.CourseSlotChanges, toCourseSlotChange(do))
	}
	return set
}

// toCourseSlotChange 领域对象 → proto 消息。
func toCourseSlotChange(do *courseslotchangedo.CourseSlotChange) *courseslotchangev1.CourseSlotChange {
	original := do.OriginalPlan()
	target := do.TargetPlan()
	return &courseslotchangev1.CourseSlotChange{
		Id:          do.ID(),
		CourseId:    do.CourseID(),
		ApplicantId: do.ApplicantID(),
		ChangeType:  courseslotchangev1.ChangeType(do.ChangeType()),
		OriginalPlan: &courseslotchangev1.OriginalPlan{
			SlotId:      original.SlotID(),
			Date:        toTimestamp(original.Date()),
			TeacherId:   original.TeacherID(),
			ClassroomId: original.ClassroomID(),
			StartTime:   original.StartTimeStr(),
			EndTime:     original.EndTimeStr(),
		},
		TargetPlan: &courseslotchangev1.TargetPlan{
			TargetStartAt: toTimestamp(target.TargetStartAt()),
			TargetEndAt:   toTimestamp(target.TargetEndAt()),
			TeacherId:     target.TeacherID(),
			ClassroomId:   target.ClassroomID(),
		},
		Reason:    do.Reason(),
		CreatedAt: toTimestamp(do.CreatedAt()),
		UpdatedAt: toTimestamp(do.UpdatedAt()),
	}
}
