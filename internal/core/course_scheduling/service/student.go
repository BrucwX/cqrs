package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	courseschedulingv1 "cqrs/api/v1/course_scheduling"
	studentv1 "cqrs/api/v1/course_scheduling/student"
	studentcmd "cqrs/internal/core/course_scheduling/app/command/student"
	studentqry "cqrs/internal/core/course_scheduling/app/query/student"
	studentdo "cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// StudentService 学员服务：把 proto DTO 转成领域输入，调 app 层，再把结果转回 DTO。
//
// 嵌入 UnimplementedStudentServiceServer 后，还没实现的 GetStudent 会
// 自动返回 codes.Unimplemented，先占着位。
type StudentService struct {
	studentv1.UnimplementedStudentServiceServer

	cmd *studentcmd.Handler
	qry *studentqry.Handler
}

// NewStudentService 创建学员服务。
func NewStudentService(cmd *studentcmd.Handler, qry *studentqry.Handler) *StudentService {
	return &StudentService{cmd: cmd, qry: qry}
}

// SaveStudent 新增或更新学员。
func (s *StudentService) SaveStudent(ctx context.Context, req *studentv1.SaveStudentRequest) (*emptypb.Empty, error) {
	in := studentcmd.StudentInput{}
	if req.Student != nil {
		// proto3 里 id=0 表示"没传"，当作新增；非 0 当作更新。
		if req.Student.Id != 0 {
			id := req.Student.Id
			in.ID = &id
		}
		in.Name = &req.Student.Name
		if req.Student.StudentType != studentv1.StudentType_STUDENT_TYPE_UNSPECIFIED {
			t := studentdo.StudentType(req.Student.StudentType)
			in.StudentType = &t
		}
		if req.Student.Status != studentv1.StudentStatus_STUDENT_STATUS_UNSPECIFIED {
			st := studentdo.Status(req.Student.Status)
			in.Status = &st
		}
		if req.Student.Contact != nil {
			c, err := studentdo.NewContactInfo(req.Student.Contact.Phone, req.Student.Contact.Email)
			if err != nil {
				return nil, err
			}
			in.Contact = &c
		}
	}
	if err := s.cmd.SaveStudent(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteStudent 删除学员。
func (s *StudentService) DeleteStudent(ctx context.Context, req *studentv1.DeleteStudentRequest) (*emptypb.Empty, error) {
	if err := s.cmd.DeleteStudent(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListStudents 分页查询学员。
func (s *StudentService) ListStudents(ctx context.Context, req *studentv1.ListStudentsRequest) (*studentv1.StudentSet, error) {
	items, err := s.qry.PageStudents(ctx, studentqry.PageStudents{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	set := &studentv1.StudentSet{Students: make([]*studentv1.Student, 0, len(items))}
	for _, do := range items {
		set.Students = append(set.Students, toStudent(do))
	}
	return set, nil
}

// toStudent 领域对象 → proto 消息。
func toStudent(do *studentdo.Student) *studentv1.Student {
	contact := do.Contact()
	return &studentv1.Student{
		Id:          do.ID(),
		Name:        do.Name(),
		StudentType: studentv1.StudentType(do.StudentType()),
		Contact: &courseschedulingv1.ContactInfo{
			Phone: contact.Phone(),
			Email: contact.Email(),
		},
		Status:    studentv1.StudentStatus(do.Status()),
		CreatedAt: timestamppb.New(do.CreatedAt()),
		UpdatedAt: timestamppb.New(do.UpdatedAt()),
	}
}
