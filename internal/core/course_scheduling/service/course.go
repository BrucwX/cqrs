package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	coursev1 "cqrs/api/v1/course_scheduling/course"
	coursecmd "cqrs/internal/core/course_scheduling/app/command/course"
	courseqry "cqrs/internal/core/course_scheduling/app/query/course"
	coursedo "cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// CourseService 课程服务。
type CourseService struct {
	coursev1.UnimplementedCourseServiceServer

	cmd *coursecmd.Handler
	qry *courseqry.Handler
}

// NewCourseService 创建课程服务。
func NewCourseService(cmd *coursecmd.Handler, qry *courseqry.Handler) *CourseService {
	return &CourseService{cmd: cmd, qry: qry}
}

// SaveCourse 新增或更新课程。
func (s *CourseService) SaveCourse(ctx context.Context, req *coursev1.SaveCourseRequest) (*emptypb.Empty, error) {
	in := coursecmd.CourseInput{}
	if req.Course != nil {
		c := req.Course
		if c.Id != "" {
			id := c.Id
			in.ID = &id
		}
		if c.CourseTypeId != "" {
			tid := c.CourseTypeId
			in.CourseTypeID = &tid
		}
		if c.Capacity != nil {
			cap, err := coursedo.NewCapacity(int(c.Capacity.Max), int(c.Capacity.Enrolled))
			if err != nil {
				return nil, err
			}
			in.Capacity = &cap
		}
		if c.Enrollment != nil {
			w := coursedo.NewEnrollmentWindow(
				tsTime(c.Enrollment.StartAt),
				tsTime(c.Enrollment.EndAt),
				tsTime(c.Enrollment.DropDeadline),
			)
			in.Enrollment = &w
		}
		if c.Period != nil {
			p, err := coursedo.NewCoursePeriod(
				tsTime(c.Period.StartAt),
				tsTime(c.Period.EndAt),
				int(c.Period.TotalHours),
				int(c.Period.CompletedHours),
			)
			if err != nil {
				return nil, err
			}
			in.Period = &p
		}
	}
	if err := s.cmd.SaveCourse(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteCourse 删除课程。
func (s *CourseService) DeleteCourse(ctx context.Context, req *coursev1.DeleteCourseRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteCourse(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListCourses 分页查询课程。
func (s *CourseService) ListCourses(ctx context.Context, req *coursev1.ListCoursesRequest) (*coursev1.CourseSet, error) {
	items, err := s.qry.PageCourses(ctx, courseqry.PageCourses{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	return newCourseSet(items), nil
}

// ListCoursesAvailableForClassroom 与教室现有排课不冲突的课程。
func (s *CourseService) ListCoursesAvailableForClassroom(ctx context.Context, req *coursev1.ListCoursesAvailableForClassroomRequest) (*coursev1.CourseSet, error) {
	items, err := s.qry.AvailableForClassroom(ctx, courseqry.AvailableForClassroom{ClassroomID: req.ClassroomId})
	if err != nil {
		return nil, err
	}
	return newCourseSet(items), nil
}

// ListCoursesAvailableForStudent 与学员当前选课不冲突的课程。
func (s *CourseService) ListCoursesAvailableForStudent(ctx context.Context, req *coursev1.ListCoursesAvailableForStudentRequest) (*coursev1.CourseSet, error) {
	items, err := s.qry.AvailableForStudent(ctx, courseqry.AvailableForStudent{StudentID: req.StudentId})
	if err != nil {
		return nil, err
	}
	return newCourseSet(items), nil
}

// ListCoursesAvailableForTeacher 与讲师现有排课不冲突的课程。
func (s *CourseService) ListCoursesAvailableForTeacher(ctx context.Context, req *coursev1.ListCoursesAvailableForTeacherRequest) (*coursev1.CourseSet, error) {
	items, err := s.qry.AvailableForTeacher(ctx, courseqry.AvailableForTeacher{TeacherID: req.TeacherId})
	if err != nil {
		return nil, err
	}
	return newCourseSet(items), nil
}

// newCourseSet []*course.Course → CourseSet。
func newCourseSet(items []*coursedo.Course) *coursev1.CourseSet {
	set := &coursev1.CourseSet{Courses: make([]*coursev1.Course, 0, len(items))}
	for _, do := range items {
		set.Courses = append(set.Courses, toCourse(do))
	}
	return set
}

// toCourse 领域对象 → proto 消息。课程没有时间戳。
func toCourse(do *coursedo.Course) *coursev1.Course {
	capacity := do.Capacity()
	window := do.Enrollment()
	period := do.Period()
	return &coursev1.Course{
		Id:           do.ID(),
		CourseTypeId: do.CourseTypeID(),
		Capacity: &coursev1.Capacity{
			Max:      int32(capacity.Max()),
			Enrolled: int32(capacity.Enrolled()),
		},
		Enrollment: &coursev1.EnrollmentWindow{
			StartAt:      toTimestamp(window.StartAt()),
			EndAt:        toTimestamp(window.EndAt()),
			DropDeadline: toTimestamp(window.DropDeadline()),
		},
		Period: &coursev1.CoursePeriod{
			StartAt:        toTimestamp(period.StartAt()),
			EndAt:          toTimestamp(period.EndAt()),
			TotalHours:     int32(period.TotalHours()),
			CompletedHours: int32(period.CompletedHours()),
		},
	}
}
