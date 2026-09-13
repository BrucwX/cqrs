package service

// TeacherService is the teaching teacher service.
//
// TODO: depend on the teacher usecases in app/command once the teaching
// write model lands.
type TeacherService struct{}

// NewTeacherService creates a new TeacherService.
func NewTeacherService() *TeacherService { return &TeacherService{} }

// CourseService is the teaching course service.
//
// TODO: depend on the course usecases in app/command once the teaching
// write model lands.
type CourseService struct{}

// NewCourseService creates a new CourseService.
func NewCourseService() *CourseService { return &CourseService{} }
