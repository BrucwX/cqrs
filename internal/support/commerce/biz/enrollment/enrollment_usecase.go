package enrollment

// EnrollmentUsecase is the enrollment usecase.
type EnrollmentUsecase struct {
	Repo EnrollmentRepo
}

// NewEnrollmentUsecase creates a new EnrollmentUsecase.
func NewEnrollmentUsecase(repo EnrollmentRepo) *EnrollmentUsecase {
	return &EnrollmentUsecase{Repo: repo}
}
