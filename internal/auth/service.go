package auth

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Authenticate(username, password string) (string, error) {
	// TODO: validate credentials, return JWT
	return "token", nil
}
