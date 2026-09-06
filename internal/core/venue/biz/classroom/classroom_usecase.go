package classroom

// ClassroomUsecase is the classroom usecase.
type ClassroomUsecase struct {
	Repo ClassroomRepo
}

// NewClassroomUsecase creates a new ClassroomUsecase.
func NewClassroomUsecase(repo ClassroomRepo) *ClassroomUsecase {
	return &ClassroomUsecase{Repo: repo}
}
