package sim

type Sim struct {
	id     int
	number string
}
type SimList map[string]Sim

func (s Sim) Id() int {
	return s.id
}
func (s Sim) Number() string {
	return s.number
}
