package core

import "database/sql"

type Sim struct {
	id            int
	number        string
	providerId    int
	isActivated   bool
	isBlocked     bool
	activateUntil int64
}

func NewSim(id int, number string, providerId int, isActivated bool, activateUntil int64, isBlocked bool) Sim {
	return Sim{
		id:            id,
		number:        number,
		providerId:    providerId,
		isActivated:   isActivated,
		activateUntil: activateUntil,
		isBlocked:     isBlocked,
	}
}
func (s Sim) Id() int {
	return s.id
}
func (s Sim) Number() string {
	return s.number
}
func (s Sim) ProviderID() int {
	return s.providerId
}
func (s Sim) IsBlocked() bool {
	return s.isBlocked
}
func (s Sim) IsActivated() bool {
	return s.isActivated
}
func (s Sim) ActivateUntil() int64 {
	return s.activateUntil
}

func (s *Sim) SetID(id int) {
	s.id = id
}
func (s *Sim) SetActivated(status bool) {
	s.isActivated = status
}
func (s *Sim) SetNumber(number string) {
	s.number = number
}
func (s *Sim) SetBlocked(status bool) {
	s.isBlocked = status
}
func (s *Sim) SetActivateUntil(aunt int64) {
	s.activateUntil = aunt
}
func (s *Sim) SetProviderID(pid int) {
	s.providerId = pid
}

// [Scan] return object of [Sim] whitch is [Scannable], and map index [int]
// If any errors ocured while scanning it will be in [error]
func (s *Sim) ScanRows(row *sql.Rows) (int, error) {
	err := row.Scan(&s.id, &s.number, &s.providerId, &s.isActivated, &s.activateUntil, &s.isBlocked)
	return s.id, err
}

func (s *Sim) ScanRow(row *sql.Row) error {
	err := row.Scan(&s.id, &s.number, &s.providerId, &s.isActivated, &s.activateUntil, &s.isBlocked)
	return err
}

func (s *Sim) GetKey() int {
	return s.Id()
}

func (s *Sim) SetKey(id int) {
	s.id = id
}
