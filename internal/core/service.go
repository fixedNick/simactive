package core

type Service struct {
	id   int
	name string
}

func (s *Service) Id() int {
	return s.id
}

func (s *Service) Name() string {
	return s.name
}
