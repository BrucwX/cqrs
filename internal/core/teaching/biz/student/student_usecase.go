package student

// StudentUsecase is the student usecase.
type StudentUsecase struct {
	Repo StudentRepo
}

// NewStudentUsecase creates a new StudentUsecase.
func NewStudentUsecase(repo StudentRepo) *StudentUsecase {
	return &StudentUsecase{Repo: repo}
}
