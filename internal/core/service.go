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

// Id returns service ID.
func (s *Service) Id() int {
	return s.id
}

// Name returns service name.
func (s *Service) Name() string {
	return s.name
}

// SetID sets service ID.
func (s *Service) SetID(id int) {
	s.id = id
}

// SetName sets service name.
func (s *Service) SetName(name string) {
	s.name = name
}

// WithName sets the name of the service.
//
// name: the name to set for the service.
// Service: the updated service with the new name set.
func (s Service) WithName(name string) Service {
	s.name = name
	return s
}

// ScanRow scans the row into the service struct.
// It returns an error if the scan fails.
//
// row: the sql rows to scan.
// error: the error, if any, that occurred during the scan.
func (s *Service) ScanRow(row *sql.Row) error {
	err := row.Scan(&s.id, &s.name)
	return err
}

// ScanRows scans the rows into the service struct.
// It returns the id of the service and an error if the scan fails.
//
// rows: the sql rows to scan.
// id: the id of the service.
// error: the error, if any, that occurred during the scan.
func (s *Service) ScanRows(rows *sql.Rows) (int, error) {
	err := rows.Scan(&s.id, &s.name)
	return s.id, err
}

// GetKey returns the map key of the Service that used in the List.
//
// id: the id of the service.
func (s *Service) GetKey() int {
	return s.Id()
}

// SetKey sets the id of the service.
//
// id: the id of the service.
func (s *Service) SetKey(id int) {
	s.id = id
}
