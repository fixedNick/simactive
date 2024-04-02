package used

type Used struct {
	id        int
	simId     int
	serviceId int
	blocked   bool
}

func (u Used) Id() int {
	return u.id
}
func (u Used) SimID() int {
	return u.simId
}
func (u Used) ServiceID() int {
	return u.serviceId
}
func (u Used) Blocked() bool {
	return u.blocked
}
