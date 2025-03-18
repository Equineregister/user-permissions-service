package permissions

type Service struct {
	repo ReaderWriter
}

// NewService creates a new permissions service
func NewService(repo ReaderWriter) *Service {
	return &Service{repo: repo}
}
