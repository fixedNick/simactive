package services

import (
	"context"
	"simactive/internal/core"
	repository "simactive/internal/infrastructure"
)

// UsedService is a service for handling operations related to used resources.
type UsedService struct {
	repository *repository.Repository
}

func NewUsedService(repo *repository.Repository) *UsedService {
	ss := &UsedService{
		repository: repo,
	}
	return ss
}

// UseSimForService is a method to mark a sim as used for a specific service.
// It creates a new entry in the 'used' table in the database.
//
// Parameters:
//   - ctx: The context.Context object for the request.
//   - simId: The ID of the sim that is being used.
//   - serviceId: The ID of the service for which the sim is being used.
//
// Returns:
//   - error: An error if the operation fails. nil if the operation is successful.
func (us *UsedService) UseSimForService(
	ctx context.Context,
	simId int,
	serviceId int,
) error {
	// Create a new used object with the provided IDs.
	used := core.Used{}.WithSimID(simId).WithServiceID(serviceId)

	// Save the used object to the 'used' table in the database.
	_, err := us.repository.UsedRepository.Add(ctx, used.SimID(), used.ServiceID(), used.IsBlocked(), used.BlockedInfo())

	// Return any error that occurred during the operation.
	return err
}
