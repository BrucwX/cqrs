package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	courseschedulingv1 "cqrs/api/v1/course_scheduling"
	teacherv1 "cqrs/api/v1/course_scheduling/teacher"
	teachercmd "cqrs/internal/core/course_scheduling/app/command/teacher"
	teacherqry "cqrs/internal/core/course_scheduling/app/query/teacher"
	teacherdo "cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// TeacherService 讲师服务。
type TeacherService struct {
	teacherv1.UnimplementedTeacherServiceServer

	cmd *teachercmd.Handler
	qry *teacherqry.Handler
}

// NewTeacherService 创建讲师服务。
func NewTeacherService(cmd *teachercmd.Handler, qry *teacherqry.Handler) *TeacherService {
	return &TeacherService{cmd: cmd, qry: qry}
}

// SaveTeacher 新增或更新讲师。
func (s *TeacherService) SaveTeacher(ctx context.Context, req *teacherv1.SaveTeacherRequest) (*emptypb.Empty, error) {
	in := teachercmd.TeacherInput{}
	if req.Teacher != nil {
		if req.Teacher.Id != 0 {
			id := req.Teacher.Id
			in.ID = &id
		}
		if req.Teacher.StudentId != 0 {
			sid := req.Teacher.StudentId
			in.StudentID = &sid
		}
		in.Name = &req.Teacher.Name
		in.Title = &req.Teacher.Title
		if req.Teacher.Status != teacherv1.TeacherStatus_TEACHER_STATUS_UNSPECIFIED {
			st := teacherdo.Status(req.Teacher.Status)
			in.Status = &st
		}
		if req.Teacher.Contact != nil {
			c, err := teacherdo.NewContactInfo(req.Teacher.Contact.Phone, req.Teacher.Contact.Email)
			if err != nil {
				return nil, err
			}
			in.Contact = &c
		}
	}
	if err := s.cmd.SaveTeacher(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteTeacher 删除讲师。
func (s *TeacherService) DeleteTeacher(ctx context.Context, req *teacherv1.DeleteTeacherRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteTeacher(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListTeachers 分页查询讲师。
func (s *TeacherService) ListTeachers(ctx context.Context, req *teacherv1.ListTeachersRequest) (*teacherv1.TeacherSet, error) {
	items, err := s.qry.PageTeachers(ctx, teacherqry.PageTeachers{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	set := &teacherv1.TeacherSet{Teachers: make([]*teacherv1.Teacher, 0, len(items))}
	for _, do := range items {
		set.Teachers = append(set.Teachers, toTeacher(do))
	}
	return set, nil
}

// toTeacher 领域对象 → proto 消息。
func toTeacher(do *teacherdo.Teacher) *teacherv1.Teacher {
	contact := do.Contact()
	return &teacherv1.Teacher{
		Id:        do.ID(),
		StudentId: do.StudentID(),
		Name:      do.Name(),
		Title:     do.Title(),
		Contact: &courseschedulingv1.ContactInfo{
			Phone: contact.Phone(),
			Email: contact.Email(),
		},
		Status:    teacherv1.TeacherStatus(do.Status()),
		CreatedAt: toTimestamp(do.CreatedAt()),
		UpdatedAt: toTimestamp(do.UpdatedAt()),
	}
}
