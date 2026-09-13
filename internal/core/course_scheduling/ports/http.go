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
	"github.com/go-kratos/kratos/v3/transport/http"
)

// HTTPServer is a wrapper for teaching HTTP server.
type HTTPServer struct {
	*http.Server
}

// NewHTTPServer creates a new HTTP server for teaching context.
func NewHTTPServer(c *conf.Server,
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
) *HTTPServer {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	studentv1.RegisterStudentServiceHTTPServer(srv, student)
	teacherv1.RegisterTeacherServiceHTTPServer(srv, teacher)
	classroomv1.RegisterClassroomServiceHTTPServer(srv, classroom)
	coursev1.RegisterCourseServiceHTTPServer(srv, course)
	courseslotv1.RegisterCourseSlotServiceHTTPServer(srv, courseSlot)
	courseslotchangev1.RegisterCourseSlotChangeServiceHTTPServer(srv, courseSlotChange)
	enrollmentv1.RegisterEnrollmentServiceHTTPServer(srv, enrollment)
	absencev1.RegisterAbsenceServiceHTTPServer(srv, absence)
	makeupv1.RegisterMakeupServiceHTTPServer(srv, makeup)
	qualificationv1.RegisterQualificationServiceHTTPServer(srv, qualification)
	return &HTTPServer{Server: srv}
}
