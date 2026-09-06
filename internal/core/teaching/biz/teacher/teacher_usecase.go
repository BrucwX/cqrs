package teacher

// TeacherUsecase is the teacher usecase.
type TeacherUsecase struct {
	Repo TeacherRepo
}

// NewTeacherUsecase creates a new TeacherUsecase.
func NewTeacherUsecase(repo TeacherRepo) *TeacherUsecase {
	return &TeacherUsecase{Repo: repo}
}
