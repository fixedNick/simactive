package repository

import "simactive/internal/core/sim"

type SimRepository interface {
	Save(s sim.Sim) error
	Remove(s sim.Sim) error
	SimList() (*sim.SimList, error)
	ById(id int) (sim.Sim, error)
	ByNumber(number string) (sim.Sim, error)
}

type SimSqlRepository interface {
	Save(s sim.Sim) (id int, err error)
	Remove(u sim.Sim) error
	SimList() (*sim.SimList, error)
	ByID(id int) (sim.Sim, error)
	ByNumber(number string) (sim.Sim, error)
}
