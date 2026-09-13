package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	courseslotv1 "cqrs/api/v1/course_scheduling/course_slot"
	courseslotcmd "cqrs/internal/core/course_scheduling/app/command/courseSlot"
	courseslotqry "cqrs/internal/core/course_scheduling/app/query/courseSlot"
	courseslotdo "cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// CourseSlotService 课表槽位服务。
type CourseSlotService struct {
	courseslotv1.UnimplementedCourseSlotServiceServer

	cmd *courseslotcmd.Handler
	qry *courseslotqry.Handler
}

// NewCourseSlotService 创建课表槽位服务。
func NewCourseSlotService(cmd *courseslotcmd.Handler, qry *courseslotqry.Handler) *CourseSlotService {
	return &CourseSlotService{cmd: cmd, qry: qry}
}

// AssignClassroomToCourseSlot 给槽位安排教室。
func (s *CourseSlotService) AssignClassroomToCourseSlot(ctx context.Context, req *courseslotv1.AssignClassroomToCourseSlotRequest) (*emptypb.Empty, error) {
	if err := s.cmd.AssignClassroom(ctx, courseslotcmd.AssignClassroomInput{
		SlotIDs:     req.SlotIds,
		ClassroomID: req.ClassroomId,
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// AssignCourseToCourseSlot 给槽位设置课程。
func (s *CourseSlotService) AssignCourseToCourseSlot(ctx context.Context, req *courseslotv1.AssignCourseToCourseSlotRequest) (*emptypb.Empty, error) {
	if err := s.cmd.AssignCourse(ctx, courseslotcmd.AssignCourseInput{
		SlotIDs:  req.SlotIds,
		CourseID: req.CourseId,
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// AssignTeacherToCourseSlot 给槽位安排讲师。
func (s *CourseSlotService) AssignTeacherToCourseSlot(ctx context.Context, req *courseslotv1.AssignTeacherToCourseSlotRequest) (*emptypb.Empty, error) {
	if err := s.cmd.AssignTeacher(ctx, courseslotcmd.AssignTeacherInput{
		SlotIDs:   req.SlotIds,
		TeacherID: req.TeacherId,
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteCourseSlot 删除槽位。
func (s *CourseSlotService) DeleteCourseSlot(ctx context.Context, req *courseslotv1.DeleteCourseSlotRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteCourseSlot(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListCourseSlots 分页查询槽位。
func (s *CourseSlotService) ListCourseSlots(ctx context.Context, req *courseslotv1.ListCourseSlotsRequest) (*courseslotv1.CourseSlotSet, error) {
	items, err := s.qry.PageCourseSlots(ctx, courseslotqry.PageCourseSlots{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	return newCourseSlotSet(items), nil
}

// ListCourseSlotsByCourseID 某课程的全部槽位。
func (s *CourseSlotService) ListCourseSlotsByCourseID(ctx context.Context, req *courseslotv1.ListCourseSlotsByCourseIDRequest) (*courseslotv1.CourseSlotSet, error) {
	items, err := s.qry.ListByCourseID(ctx, courseslotqry.ListByCourseID{CourseID: req.CourseId})
	if err != nil {
		return nil, err
	}
	return newCourseSlotSet(items), nil
}

// newCourseSlotSet []*courseSlot.CourseSlot → CourseSlotSet。
func newCourseSlotSet(items []*courseslotdo.CourseSlot) *courseslotv1.CourseSlotSet {
	set := &courseslotv1.CourseSlotSet{CourseSlots: make([]*courseslotv1.CourseSlot, 0, len(items))}
	for _, do := range items {
		set.CourseSlots = append(set.CourseSlots, toCourseSlot(do))
	}
	return set
}

// toCourseSlot 领域对象 → proto 消息。
//
// 注意 Weekday：domain 用 time.Weekday（Sunday=0 … Saturday=6），
// proto 枚举是 SUNDAY=1 … SATURDAY=7，所以要 +1。
func toCourseSlot(do *courseslotdo.CourseSlot) *courseslotv1.CourseSlot {
	tr := do.TimeRange()
	start, end := tr.Start(), tr.End()
	return &courseslotv1.CourseSlot{
		Id:       do.ID(),
		CourseId: do.CourseID(),
		Weekday:  courseslotv1.Weekday(do.Weekday() + 1),
		TimeRange: &courseslotv1.DayTimeRange{
			Start: &courseslotv1.DayTime{Hour: int32(start.TotalMinutes() / 60), Minute: int32(start.TotalMinutes() % 60)},
			End:   &courseslotv1.DayTime{Hour: int32(end.TotalMinutes() / 60), Minute: int32(end.TotalMinutes() % 60)},
		},
		TeacherId:   do.TeacherID(),
		ClassroomId: do.ClassroomID(),
		CreatedAt:   toTimestamp(do.CreatedAt()),
		UpdatedAt:   toTimestamp(do.UpdatedAt()),
	}
}
