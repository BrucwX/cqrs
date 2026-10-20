package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	classroomv1 "courseScheduling/api/v1/course_scheduling/classroom"
	classroomcmd "courseScheduling/internal/app/command/classroom"
	classroomqry "courseScheduling/internal/app/query/classroom"
	classroomdo "courseScheduling/internal/domain/aggregate/classroom"
)

// ClassroomService 教室服务。
type ClassroomService struct {
	classroomv1.UnimplementedClassroomServiceServer

	cmd *classroomcmd.Handler
	qry *classroomqry.Handler
}

// NewClassroomService 创建教室服务。
func NewClassroomService(cmd *classroomcmd.Handler, qry *classroomqry.Handler) *ClassroomService {
	return &ClassroomService{cmd: cmd, qry: qry}
}

// SaveClassroom 新增或更新教室：没带 ID 走新增（ID 由聚合生成），带了 ID 走更新。
func (s *ClassroomService) SaveClassroom(ctx context.Context, req *classroomv1.SaveClassroomRequest) (*emptypb.Empty, error) {
	var (
		locPtr *classroomdo.Location
		capPtr *int
		stPtr  *classroomdo.Status
	)
	if c := req.GetClassroom(); c != nil {
		if c.Location != nil {
			loc, err := classroomdo.NewLocation(
				c.Location.Building,
				int(c.Location.Floor),
				c.Location.Room,
			)
			if err != nil {
				return nil, err
			}
			locPtr = &loc
		}
		if c.Capacity != 0 {
			capacity := int(c.Capacity)
			capPtr = &capacity
		}
		if c.Status != classroomv1.ClassroomStatus_CLASSROOM_STATUS_UNSPECIFIED {
			status := classroomdo.Status(c.Status)
			stPtr = &status
		}
	}

	if id := req.GetClassroom().GetId(); id != "" {
		err := s.cmd.UpdateClassroom(ctx, classroomcmd.UpdateClassroomInput{
			ID:       id,
			Location: locPtr,
			Capacity: capPtr,
			Status:   stPtr,
		})
		if err != nil {
			return nil, err
		}
		return &emptypb.Empty{}, nil
	}

	if err := s.cmd.CreateClassroom(ctx, classroomcmd.CreateClassroomInput{
		Location: locPtr,
		Capacity: capPtr,
		Status:   stPtr,
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteClassroom 删除教室。
func (s *ClassroomService) DeleteClassroom(ctx context.Context, req *classroomv1.DeleteClassroomRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteClassroom(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListClassrooms 分页查询教室。
func (s *ClassroomService) ListClassrooms(ctx context.Context, req *classroomv1.ListClassroomsRequest) (*classroomv1.ClassroomSet, error) {
	items, err := s.qry.PageClassrooms(ctx, classroomqry.PageClassrooms{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	set := &classroomv1.ClassroomSet{Classrooms: make([]*classroomv1.Classroom, 0, len(items))}
	for _, do := range items {
		set.Classrooms = append(set.Classrooms, toClassroom(do))
	}
	return set, nil
}

// toClassroom 领域对象 → proto 消息。教室没有时间戳。
func toClassroom(do *classroomdo.Classroom) *classroomv1.Classroom {
	loc := do.Location()
	return &classroomv1.Classroom{
		Id: do.ID(),
		Location: &classroomv1.Location{
			Building: loc.Building(),
			Floor:    int32(loc.Floor()),
			Room:     loc.Room(),
		},
		Capacity:  int32(do.Capacity()),
		Allocated: int32(do.AllocatedSeats()),
		Status:    classroomv1.ClassroomStatus(do.Status()),
	}
}
