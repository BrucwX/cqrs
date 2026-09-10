package service

// TeacherService is the teaching teacher service.
//
// TODO: depend on the teacher usecases in app/command once the teaching
// write model lands.
type TeacherService struct{}

// NewTeacherService creates a new TeacherService.
func NewTeacherService() *TeacherService { return &TeacherService{} }

// StudentService is the teaching student service.
//
// TODO: depend on the student usecases in app/command once the teaching
// write model lands.
type StudentService struct{}

// NewStudentService creates a new StudentService.
func NewStudentService() *StudentService { return &StudentService{} }

// CourseService is the teaching course service.
//
// TODO: depend on the course usecases in app/command once the teaching
// write model lands.
type CourseService struct{}

// NewCourseService creates a new CourseService.
func NewCourseService() *CourseService { return &CourseService{} }
