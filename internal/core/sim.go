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

// NewSim creates a new Sim object with the given parameters.
// Parameters:
//   - id: the id of the Sim.
//   - number: the number of the Sim.
//   - providerId: the provider id of the Sim.
//   - isActivated: a boolean indicating if the Sim is activated.
//   - activateUntil: the timestamp until when the Sim is activated.
//   - isBlocked: a boolean indicating if the Sim is blocked.
// Returns:
//   - Sim: the newly created Sim object.
func NewSim(id int, number string, providerId int, isActivated bool, activateUntil int64, isBlocked bool) Sim {
	// Create and return a new Sim object with the given parameters.
	return Sim{
		id:            id,            // Set the id of the Sim.
		number:        number,        // Set the number of the Sim.
		providerId:    providerId,    // Set the provider id of the Sim.
		isActivated:   isActivated,   // Set if the Sim is activated.
		activateUntil: activateUntil, // Set the timestamp until when the Sim is activated.
		isBlocked:     isBlocked,     // Set if the Sim is blocked.
	}
}

// Getters

func (s Sim) Id() int              { return s.id }
func (s Sim) Number() string       { return s.number }
func (s Sim) ProviderID() int      { return s.providerId }
func (s Sim) IsBlocked() bool      { return s.isBlocked }
func (s Sim) IsActivated() bool    { return s.isActivated }
func (s Sim) ActivateUntil() int64 { return s.activateUntil }

// Setters

func (s *Sim) SetID(id int)                { s.id = id }
func (s *Sim) SetNumber(number string)     { s.number = number }
func (s *Sim) SetProviderID(pid int)       { s.providerId = pid }
func (s *Sim) SetBlocked(status bool)      { s.isBlocked = status }
func (s *Sim) SetActivated(status bool)    { s.isActivated = status }
func (s *Sim) SetActivateUntil(aunt int64) { s.activateUntil = aunt }

// ScanRow scans the values from the given sql.Rows into the fields of the Sim struct.
//
// It takes a pointer to a sql.Row as parameter and returns an error.
// Scanned values are stored in the fields of the pointer to the Sim struct.
func (s *Sim) ScanRows(row *sql.Rows) (int, error) {
	err := row.Scan(&s.id, &s.number, &s.providerId, &s.isActivated, &s.activateUntil, &s.isBlocked)
	return s.id, err
}

// ScanRow scans the values from the given sql.Row into the fields of the Sim struct.
//
// It takes a pointer to a sql.Row as parameter and returns an error.
// Scanned values are stored in the fields of the pointer to the Sim struct.
func (s *Sim) ScanRow(row *sql.Row) error {
	err := row.Scan(&s.id, &s.number, &s.providerId, &s.isActivated, &s.activateUntil, &s.isBlocked)
	return err
}

// GetKey returns the map key of the Sim that used in the List.
//
// Returns an integer.
func (s *Sim) GetKey() int {
	return s.Id()
}

// SetKey sets the id for the Sim struct.
//
// Parameter: id int - the new id to set.
func (s *Sim) SetKey(id int) {
	s.id = id
}
