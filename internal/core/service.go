package core

import "database/sql"

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

// [Scan] return object of [Service] whitch is [Scannable], and map index [int]
// If any errors ocured while scanning it will be in [error]
func (s Service) Scan(rows *sql.Rows) (Scannable, int, error) {
	scannedService := Service{}
	err := rows.Scan(&scannedService.id, &scannedService.name)
	if err != nil {
		return scannedService, 0, err
	}
	return scannedService, scannedService.id, err
}
