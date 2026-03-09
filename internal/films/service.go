package films

type Film struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewRepository()}
}

func (s *Service) List() ([]Film, error) {
	return s.repo.FindAll()
}

func (s *Service) Get(id string) (*Film, error) {
	return s.repo.FindByID(id)
}
