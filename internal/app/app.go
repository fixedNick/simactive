package app

import repository "simactive/internal/storage"

type Application struct {
	InMemStorage repository.SimRepository
}
