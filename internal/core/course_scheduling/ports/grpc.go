package ports

import (
	absencev1 "cqrs/api/v1/course_scheduling/absence"
	classroomv1 "cqrs/api/v1/course_scheduling/classroom"
	coursev1 "cqrs/api/v1/course_scheduling/course"
	courseslotv1 "cqrs/api/v1/course_scheduling/course_slot"
	courseslotchangev1 "cqrs/api/v1/course_scheduling/course_slot_change"
	enrollmentv1 "cqrs/api/v1/course_scheduling/enrollment"
	makeupv1 "cqrs/api/v1/course_scheduling/makeup"
	qualificationv1 "cqrs/api/v1/course_scheduling/qualification"
	studentv1 "cqrs/api/v1/course_scheduling/student"
	teacherv1 "cqrs/api/v1/course_scheduling/teacher"
	"cqrs/internal/conf"
	"cqrs/internal/core/course_scheduling/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
)

// GRPCServer is a wrapper for teaching gRPC server.
type GRPCServer struct {
	*grpc.Server
}

// NewGRPCServer creates a new gRPC server for teaching context.
func NewGRPCServer(c *conf.Server,
	student *service.StudentService,
	teacher *service.TeacherService,
	classroom *service.ClassroomService,
	course *service.CourseService,
	courseSlot *service.CourseSlotService,
	courseSlotChange *service.CourseSlotChangeService,
	enrollment *service.EnrollmentService,
	absence *service.AbsenceService,
	makeup *service.MakeupService,
	qualification *service.QualificationService,
) *GRPCServer {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	studentv1.RegisterStudentServiceServer(srv, student)
	teacherv1.RegisterTeacherServiceServer(srv, teacher)
	classroomv1.RegisterClassroomServiceServer(srv, classroom)
	coursev1.RegisterCourseServiceServer(srv, course)
	courseslotv1.RegisterCourseSlotServiceServer(srv, courseSlot)
	courseslotchangev1.RegisterCourseSlotChangeServiceServer(srv, courseSlotChange)
	enrollmentv1.RegisterEnrollmentServiceServer(srv, enrollment)
	absencev1.RegisterAbsenceServiceServer(srv, absence)
	makeupv1.RegisterMakeupServiceServer(srv, makeup)
	qualificationv1.RegisterQualificationServiceServer(srv, qualification)
	return &GRPCServer{Server: srv}
}
