package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	coursetypev1 "cqrs/api/v1/course_scheduling/course_type"
	qualificationv1 "cqrs/api/v1/course_scheduling/qualification"
	teacherv1 "cqrs/api/v1/course_scheduling/teacher"
	qualificationcmd "cqrs/internal/core/course_scheduling/app/command/qualification"
	qualificationqry "cqrs/internal/core/course_scheduling/app/query/qualification"
	coursetypedo "cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	qualificationdo "cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

// QualificationService 授课资质服务。
type QualificationService struct {
	qualificationv1.UnimplementedQualificationServiceServer

	cmd *qualificationcmd.Handler
	qry *qualificationqry.Handler
}

// NewQualificationService 创建授课资质服务。
func NewQualificationService(cmd *qualificationcmd.Handler, qry *qualificationqry.Handler) *QualificationService {
	return &QualificationService{cmd: cmd, qry: qry}
}

// Qualify 授予授课资质。
func (s *QualificationService) Qualify(ctx context.Context, req *qualificationv1.QualifyRequest) (*emptypb.Empty, error) {
	in := qualificationcmd.QualifyInput{}
	if req.Qualification != nil {
		in.TeacherID = req.Qualification.TeacherId
		in.CourseTypeID = req.Qualification.CourseTypeId
	}
	if err := s.cmd.Qualify(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteQualification 删除授课资质。
func (s *QualificationService) DeleteQualification(ctx context.Context, req *qualificationv1.DeleteQualificationRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteQualification(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListQualifications 分页查询授课资质。
func (s *QualificationService) ListQualifications(ctx context.Context, req *qualificationv1.ListQualificationsRequest) (*qualificationv1.QualificationSet, error) {
	items, err := s.qry.PageQualifications(ctx, qualificationqry.PageQualifications{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	set := &qualificationv1.QualificationSet{Qualifications: make([]*qualificationv1.Qualification, 0, len(items))}
	for _, do := range items {
		set.Qualifications = append(set.Qualifications, toQualification(do))
	}
	return set, nil
}

// ListQualificationsByTeacherID 讲师有资质的课程类型（跨聚合：返回课程类型）。
func (s *QualificationService) ListQualificationsByTeacherID(ctx context.Context, req *qualificationv1.ListQualificationsByTeacherIDRequest) (*coursetypev1.CourseTypeSet, error) {
	items, err := s.qry.CourseTypesByTeacherID(ctx, qualificationqry.CourseTypesByTeacherID{TeacherID: req.TeacherId})
	if err != nil {
		return nil, err
	}
	set := &coursetypev1.CourseTypeSet{CourseTypes: make([]*coursetypev1.CourseType, 0, len(items))}
	for _, do := range items {
		set.CourseTypes = append(set.CourseTypes, toCourseType(do))
	}
	return set, nil
}

// ListTeachersQualifiedForCourseType 有某课程类型资质的讲师（跨聚合：返回讲师）。
func (s *QualificationService) ListTeachersQualifiedForCourseType(ctx context.Context, req *qualificationv1.ListTeachersQualifiedForCourseTypeRequest) (*teacherv1.TeacherSet, error) {
	items, err := s.qry.TeachersByCourseTypeID(ctx, qualificationqry.TeachersByCourseTypeID{CourseTypeID: req.CourseTypeId})
	if err != nil {
		return nil, err
	}
	set := &teacherv1.TeacherSet{Teachers: make([]*teacherv1.Teacher, 0, len(items))}
	for _, do := range items {
		set.Teachers = append(set.Teachers, toTeacher(do))
	}
	return set, nil
}

// CheckTeacherQualifiedForCourse 讲师是否有教某门课的资质。
func (s *QualificationService) CheckTeacherQualifiedForCourse(ctx context.Context, req *qualificationv1.CheckTeacherQualifiedForCourseRequest) (*qualificationv1.CheckTeacherQualifiedForCourseResponse, error) {
	ok, err := s.qry.IsTeacherQualifiedForCourse(ctx, qualificationqry.IsTeacherQualifiedForCourse{
		TeacherID: req.TeacherId,
		CourseID:  req.CourseId,
	})
	if err != nil {
		return nil, err
	}
	return &qualificationv1.CheckTeacherQualifiedForCourseResponse{Qualified: ok}, nil
}

// toQualification 领域对象 → proto 消息。资质没有创建时间，只有认证时间。
func toQualification(do *qualificationdo.Qualification) *qualificationv1.Qualification {
	return &qualificationv1.Qualification{
		Id:           do.ID(),
		TeacherId:    do.TeacherID(),
		CourseTypeId: do.CourseTypeID(),
		CertifiedAt:  toTimestamp(do.CertifiedAt()),
		Status:       qualificationv1.QualificationStatus(do.Status()),
		UpdatedAt:    toTimestamp(do.UpdatedAt()),
	}
}

// toCourseType 领域对象 → proto 消息。
func toCourseType(do *coursetypedo.CourseType) *coursetypev1.CourseType {
	return &coursetypev1.CourseType{
		Id:          do.ID(),
		Name:        do.Name(),
		Description: do.Description(),
		CreatedAt:   toTimestamp(do.CreatedAt()),
		UpdatedAt:   toTimestamp(do.UpdatedAt()),
	}
}
