package tests

import (
	"fmt"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/core"
	"simactive/internal/tests/suite"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeleteSim_HappyPath tests the successful deletion of a SIM.
// Uses a fake SIM data to add a SIM, then deletes it and verifies the deletion.
func TestDeleteSim_HappyPath(t *testing.T) {
	ctx, ss := suite.NewSuite(t)

	sim := core.NewSim(
		0,
		suite.GenerateFakePhoneNumber(),
		gofakeit.Number(1, 999),
		gofakeit.Bool(),
		suite.GenerateFakeDateUnix(),
		gofakeit.Bool(),
	)

	resp, err := ss.SimClient.AddSim(ctx, &pb.AddSimRequest{
		SimData: &pb.AddSimData{
			Number:        sim.Number(),
			ProviderID:    int32(sim.ProviderID()),
			IsActivated:   sim.IsActivated(),
			ActivateUntil: sim.ActivateUntil(),
			IsBlocked:     sim.IsBlocked(),
		},
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetId())

	sim.SetID(int(resp.GetId()))

	deleteResp, err := ss.SimClient.DeleteSim(
		ctx,
		&pb.DeleteSimRequest{
			Id: int32(sim.Id()),
		},
	)
	require.NoError(t, err)
	assert.NotEmpty(t, deleteResp.GetId())
	assert.Equal(t, sim.Id(), int(deleteResp.GetId()))

	// Get list and check that it do not contain deleted sim
	list, err := ss.SimClient.GetSimList(ctx, nil)
	require.NoError(t, err)
	assert.NotContains(t, list.GetSimList(), sim)
	assert.NotEmpty(t, deleteResp)
}

// TestDeleteSim_FailCases tests the failure cases of the DeleteSim function.
// It tests deleting a SIM card with a not existing ID and with an empty ID.
func TestDeleteSim_FailCases(t *testing.T) {
	ctx, s := suite.NewSuite(t)
	// invalid input
	// sim not found
	tests := []struct {
		name               string
		id                 int
		expectedErr        string
		expectedStatusCode codes.Code
	}{
		{
			name:               "Delete sim with not existing id",
			id:                 999999999, // hope that sim with this id does not exist ;)
			expectedErr:        fmt.Sprintf("sim card with id %d not found", 999999999),
			expectedStatusCode: codes.NotFound,
		},
		{
			name:               "Delete sim with empty id",
			id:                 0,
			expectedErr:        "Invalid id, id must be greater than 0",
			expectedStatusCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteResp, err := s.SimClient.DeleteSim(
				ctx,
				&pb.DeleteSimRequest{
					Id: int32(tt.id),
				},
			)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
			assert.Empty(t, deleteResp.GetId())

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tt.expectedStatusCode, st.Code())
		})
	}
}
