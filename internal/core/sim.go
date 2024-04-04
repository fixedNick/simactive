package core

type Sim struct {
	id            int
	number        string
	providerId    int
	isActivated   bool
	isBlocked     bool
	activateUntil int64
}
type SimList map[int]Sim

func NewSimList(s ...Sim) SimList {
	list := make(SimList)
	for _, sim := range s {
		list[sim.Id()] = sim
	}
	return list
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
