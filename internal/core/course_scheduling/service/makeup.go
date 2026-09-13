package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	coursev1 "cqrs/api/v1/course_scheduling/course"
	makeupv1 "cqrs/api/v1/course_scheduling/makeup"
	studentv1 "cqrs/api/v1/course_scheduling/student"
	makeupcmd "cqrs/internal/core/course_scheduling/app/command/makeup"
	makeupqry "cqrs/internal/core/course_scheduling/app/query/makeup"
	makeupdo "cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

// MakeupService 补课服务。
type MakeupService struct {
	makeupv1.UnimplementedMakeupServiceServer

	cmd *makeupcmd.Handler
	qry *makeupqry.Handler
}

// NewMakeupService 创建补课服务。
func NewMakeupService(cmd *makeupcmd.Handler, qry *makeupqry.Handler) *MakeupService {
	return &MakeupService{cmd: cmd, qry: qry}
}

// RecordMakeup 记录补课。
func (s *MakeupService) RecordMakeup(ctx context.Context, req *makeupv1.RecordMakeupRequest) (*emptypb.Empty, error) {
	in := makeupcmd.RecordMakeup{}
	if req.Makeup != nil {
		m := req.Makeup
		in.StudentID = m.StudentId
		in.CourseID = m.CourseId
		in.OriginalSlotID = m.OriginalSlotId
		in.OriginalDate = tsTime(m.OriginalDate)
		in.TargetSlotID = m.TargetSlotId
		in.TargetDate = tsTime(m.TargetDate)
		in.MakeupHours = int(m.MakeupHours)
	}
	if _, err := s.cmd.RecordMakeup(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteMakeup 删除补课记录。
func (s *MakeupService) DeleteMakeup(ctx context.Context, req *makeupv1.DeleteMakeupRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteMakeup(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListMakeups 分页查询补课记录。
func (s *MakeupService) ListMakeups(ctx context.Context, req *makeupv1.ListMakeupsRequest) (*makeupv1.StudentMakeupSet, error) {
	items, err := s.qry.PageMakeups(ctx, makeupqry.PageMakeups{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	set := &makeupv1.StudentMakeupSet{Makeups: make([]*makeupv1.StudentMakeup, 0, len(items))}
	for _, do := range items {
		set.Makeups = append(set.Makeups, toMakeup(do))
	}
	return set, nil
}

// ListMakeupsByStudentID 学员补课涉及的课程（跨聚合：返回课程）。
func (s *MakeupService) ListMakeupsByStudentID(ctx context.Context, req *makeupv1.ListMakeupsByStudentIDRequest) (*coursev1.CourseSet, error) {
	items, err := s.qry.CoursesByStudentID(ctx, makeupqry.CoursesByStudentID{StudentID: req.StudentId})
	if err != nil {
		return nil, err
	}
	return newCourseSet(items), nil
}

// ListMakeupsByCourseID 课程补课涉及的学员（跨聚合：返回学员）。
func (s *MakeupService) ListMakeupsByCourseID(ctx context.Context, req *makeupv1.ListMakeupsByCourseIDRequest) (*studentv1.StudentSet, error) {
	items, err := s.qry.StudentsByCourseID(ctx, makeupqry.StudentsByCourseID{CourseID: req.CourseId})
	if err != nil {
		return nil, err
	}
	return newStudentSet(items), nil
}

// toMakeup 领域对象 → proto 消息。
func toMakeup(do *makeupdo.StudentMakeup) *makeupv1.StudentMakeup {
	return &makeupv1.StudentMakeup{
		Id:             do.ID(),
		StudentId:      do.StudentID(),
		CourseId:       do.CourseID(),
		OriginalSlotId: do.OriginalSlotID(),
		OriginalDate:   toTimestamp(do.OriginalDate()),
		TargetSlotId:   do.TargetSlotID(),
		TargetDate:     toTimestamp(do.TargetDate()),
		MakeupHours:    int32(do.MakeupHours()),
		Status:         makeupv1.MakeupStatus(do.Status()),
		CompletedAt:    toTimestampOrNil(do.CompletedAt()),
		CreatedAt:      toTimestamp(do.CreatedAt()),
		UpdatedAt:      toTimestamp(do.UpdatedAt()),
	}
}
