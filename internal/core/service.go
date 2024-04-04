package core

type Service struct {
	id   int
	name string
}

func NewService(id int, name string) Service {
	return Service{
		id:   id,
		name: name,
	}
}

func (s *Service) Id() int {
	return s.id
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) SetID(id int) {
	s.id = id
}

func (s *Service) SetName(name string) {
	s.name = name
}
