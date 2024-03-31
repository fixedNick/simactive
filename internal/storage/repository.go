package repository

import "simactive/internal/core/sim"

type SimRepository interface {
	Save(s sim.Sim) error
	Remove(u sim.Sim) error
	SimList() (sim.SimList, error)
	ById(id int) sim.Sim
	ByNumber(number string) sim.Sim
}
