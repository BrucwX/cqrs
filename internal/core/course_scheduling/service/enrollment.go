package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	coursev1 "cqrs/api/v1/course_scheduling/course"
	enrollmentv1 "cqrs/api/v1/course_scheduling/enrollment"
	studentv1 "cqrs/api/v1/course_scheduling/student"
	enrollmentcmd "cqrs/internal/core/course_scheduling/app/command/enrollment"
	enrollmentqry "cqrs/internal/core/course_scheduling/app/query/enrollment"
	enrollmentdo "cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// EnrollmentService 报名服务。
type EnrollmentService struct {
	enrollmentv1.UnimplementedEnrollmentServiceServer

	cmd *enrollmentcmd.Handler
	qry *enrollmentqry.Handler
}

// NewEnrollmentService 创建报名服务。
func NewEnrollmentService(cmd *enrollmentcmd.Handler, qry *enrollmentqry.Handler) *EnrollmentService {
	return &EnrollmentService{cmd: cmd, qry: qry}
}

// EnrollStudent 学生选课。
func (s *EnrollmentService) EnrollStudent(ctx context.Context, req *enrollmentv1.EnrollStudentRequest) (*emptypb.Empty, error) {
	in := enrollmentcmd.StudentEnroll{}
	if req.Enrollment != nil {
		in.StudentID = req.Enrollment.StudentId
		in.CourseID = req.Enrollment.CourseId
	}
	if _, err := s.cmd.StudentEnroll(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteEnrollment 删除报名记录。
func (s *EnrollmentService) DeleteEnrollment(ctx context.Context, req *enrollmentv1.DeleteEnrollmentRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteEnrollment(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListEnrollments 分页查询报名记录。
func (s *EnrollmentService) ListEnrollments(ctx context.Context, req *enrollmentv1.ListEnrollmentsRequest) (*enrollmentv1.CourseEnrollmentSet, error) {
	items, err := s.qry.PageEnrollments(ctx, enrollmentqry.PageEnrollments{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	set := &enrollmentv1.CourseEnrollmentSet{Enrollments: make([]*enrollmentv1.CourseEnrollment, 0, len(items))}
	for _, do := range items {
		set.Enrollments = append(set.Enrollments, toEnrollment(do))
	}
	return set, nil
}

// ListEnrolledCoursesByStudentID 学员已报名的课程（跨聚合：返回课程）。
func (s *EnrollmentService) ListEnrolledCoursesByStudentID(ctx context.Context, req *enrollmentv1.ListEnrolledCoursesByStudentIDRequest) (*coursev1.CourseSet, error) {
	items, err := s.qry.CoursesByStudentID(ctx, enrollmentqry.CoursesByStudentID{StudentID: req.StudentId})
	if err != nil {
		return nil, err
	}
	return newCourseSet(items), nil
}

// ListEnrolledStudentsByCourseID 课程已报名的学员（跨聚合：返回学员）。
func (s *EnrollmentService) ListEnrolledStudentsByCourseID(ctx context.Context, req *enrollmentv1.ListEnrolledStudentsByCourseIDRequest) (*studentv1.StudentSet, error) {
	items, err := s.qry.StudentsByCourseID(ctx, enrollmentqry.StudentsByCourseID{CourseID: req.CourseId})
	if err != nil {
		return nil, err
	}
	return newStudentSet(items), nil
}

// toEnrollment 领域对象 → proto 消息。
func toEnrollment(do *enrollmentdo.CourseEnrollment) *enrollmentv1.CourseEnrollment {
	return &enrollmentv1.CourseEnrollment{
		Id:          do.ID(),
		StudentId:   do.StudentID(),
		CourseId:    do.CourseID(),
		Status:      enrollmentv1.EnrollmentStatus(do.Status()),
		EnrolledAt:  toTimestamp(do.EnrolledAt()),
		CompletedAt: toTimestampOrNil(do.CompletedAt()),
		DroppedAt:   toTimestampOrNil(do.DroppedAt()),
		UpdatedAt:   toTimestamp(do.UpdatedAt()),
	}
}
