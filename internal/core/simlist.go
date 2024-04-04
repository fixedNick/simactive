package core

import "fmt"

type SimList map[int]Sim

func NewSimList(s ...Sim) SimList {
	list := make(SimList)
	for _, sim := range s {
		list[sim.Id()] = sim
	}
	return list
}

func (sl *SimList) PtrByID(id int) (*Sim, error) {
	s, ok := (*sl)[id]
	if !ok {
		return nil, fmt.Errorf("sim with id `%d` not found in list", id)
	}
	return &s, nil
}
