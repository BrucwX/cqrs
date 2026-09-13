package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	absencev1 "cqrs/api/v1/course_scheduling/absence"
	absencecmd "cqrs/internal/core/course_scheduling/app/command/absence"
	absenceqry "cqrs/internal/core/course_scheduling/app/query/absence"
	absencedo "cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// AbsenceService 缺勤记录服务。
type AbsenceService struct {
	absencev1.UnimplementedAbsenceServiceServer

	cmd *absencecmd.Handler
	qry *absenceqry.Handler
}

// NewAbsenceService 创建缺勤记录服务。
func NewAbsenceService(cmd *absencecmd.Handler, qry *absenceqry.Handler) *AbsenceService {
	return &AbsenceService{cmd: cmd, qry: qry}
}

// RecordAbsence 记录缺勤。
func (s *AbsenceService) RecordAbsence(ctx context.Context, req *absencev1.RecordAbsenceRequest) (*emptypb.Empty, error) {
	in := absencecmd.RecordAbsence{}
	if req.Absence != nil {
		a := req.Absence
		in.StudentID = a.StudentId
		in.CourseID = a.CourseId
		in.CourseSlotID = a.CourseSlotId
		in.ScheduleDate = tsTime(a.ScheduleDate)
		in.MissedHours = int(a.MissedHours)
		in.AbsenceType = absencedo.AbsenceType(a.AbsenceType)
		in.Reason = a.Reason
	}
	if _, err := s.cmd.RecordAbsence(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteAbsence 删除缺勤记录。
func (s *AbsenceService) DeleteAbsence(ctx context.Context, req *absencev1.DeleteAbsenceRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteAbsence(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListAbsences 分页查询缺勤记录。
func (s *AbsenceService) ListAbsences(ctx context.Context, req *absencev1.ListAbsencesRequest) (*absencev1.AbsenceRecordSet, error) {
	items, err := s.qry.PageAbsences(ctx, absenceqry.PageAbsences{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	return newAbsenceSet(items), nil
}

// ListAbsencesByStudentID 某学员的全部缺勤记录。
func (s *AbsenceService) ListAbsencesByStudentID(ctx context.Context, req *absencev1.ListAbsencesByStudentIDRequest) (*absencev1.AbsenceRecordSet, error) {
	items, err := s.qry.ListByStudentID(ctx, absenceqry.ListByStudentID{StudentID: req.StudentId})
	if err != nil {
		return nil, err
	}
	return newAbsenceSet(items), nil
}

// ListAbsencesByCourseID 某课程的全部缺勤记录。
func (s *AbsenceService) ListAbsencesByCourseID(ctx context.Context, req *absencev1.ListAbsencesByCourseIDRequest) (*absencev1.AbsenceRecordSet, error) {
	items, err := s.qry.ListByCourseID(ctx, absenceqry.ListByCourseID{CourseID: req.CourseId})
	if err != nil {
		return nil, err
	}
	return newAbsenceSet(items), nil
}

// newAbsenceSet []*AbsenceRecord → AbsenceRecordSet。
func newAbsenceSet(items []*absencedo.AbsenceRecord) *absencev1.AbsenceRecordSet {
	set := &absencev1.AbsenceRecordSet{Absences: make([]*absencev1.AbsenceRecord, 0, len(items))}
	for _, do := range items {
		set.Absences = append(set.Absences, toAbsence(do))
	}
	return set
}

// toAbsence 领域对象 → proto 消息。
func toAbsence(do *absencedo.AbsenceRecord) *absencev1.AbsenceRecord {
	return &absencev1.AbsenceRecord{
		Id:           do.ID(),
		StudentId:    do.StudentID(),
		CourseId:     do.CourseID(),
		CourseSlotId: do.CourseSlotID(),
		ScheduleDate: toTimestamp(do.ScheduleDate()),
		MissedHours:  int32(do.MissedHours()),
		AbsenceType:  absencev1.AbsenceType(do.AbsenceType()),
		Reason:       do.Reason(),
		CreatedAt:    toTimestamp(do.CreatedAt()),
		UpdatedAt:    toTimestamp(do.UpdatedAt()),
	}
}
