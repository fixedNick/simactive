package core

import "fmt"

type ServiceList map[int]Service

func NewServiceList(services ...Service) ServiceList {
	sl := make(ServiceList)

	for _, s := range services {
		sl[s.Id()] = NewService(s.Id(), s.Name())
	}
	return sl
}

func (sl *ServiceList) ByID(id int) (*Service, error) {
	service, ok := (*sl)[id]
	if !ok {
		return nil, fmt.Errorf("cannot find service with id `%d` in list", id)
	}
	return &service, nil
}
