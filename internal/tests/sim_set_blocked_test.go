package tests

import (
	"simactive/internal/core"
	"simactive/internal/tests/suite"
	"testing"

	pb "simactive/api/generated/github.com/fixedNick/SimHelper"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestActivateSim_HappyPath is a test function for activating a sim card in a happy path scenario.
func TestSetBlockedSim_HappyPath(t *testing.T) {

	ctx, s := suite.NewSuite(t)

	// add new not blocked sim
	// block sim
	// check that sim is blocked

	fakeProvider := core.Provider{}.WithName(suite.GenerateFakeString(16))
	sim := core.NewSim(
		0,
		suite.GenerateFakePhoneNumber(),
		&fakeProvider,
		gofakeit.Bool(),
		suite.GenerateFakeDateUnix(),
		false,
	)

	resp, err := s.SimClient.AddSim(ctx, &pb.AddSimRequest{
		SimData: &pb.AddSimData{
			Number:        sim.Number(),
			ProviderName:  sim.Provider().Name(),
			IsActivated:   sim.IsActivated(),
			ActivateUntil: sim.ActivateUntil(),
			IsBlocked:     sim.IsBlocked(),
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetId())

	sim.SetID(int(resp.GetId()))
	require.Equal(t, sim.Id(), int(resp.GetId()))

	_, err = s.SimClient.SetSimBlocked(ctx, &pb.SSBRequest{
		Id: int32(sim.Id()),
	})
	require.NoError(t, err)

	// get sim list
	// find sim by id
	// check that sim is blocked

	simList, err := s.SimClient.GetSimList(ctx, &pb.Empty{})
	require.NoError(t, err)
	require.NotNil(t, simList)
	require.NotEmpty(t, simList.GetSimList())

	isSimFound := false
	for _, v := range simList.SimList {
		if v.ID == int32(sim.Id()) {
			isSimFound = true
			assert.True(t, v.IsBlocked)
		}
	}

	require.True(t, isSimFound)
}

func TestSetBlockedSim_FailCases(t *testing.T) {
	ctx, s := suite.NewSuite(t)

	// id is invalid
	// sim not found

	tests := []struct {
		name               string
		id                 int
		expectedErr        string
		expectedStatusCode codes.Code
	}{
		{
			name:               "Block sim with invalid id",
			id:                 0,
			expectedErr:        "Invalid id, id must be greater than 0",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Block sim with not existing id",
			id:                 999999999,
			expectedErr:        "sim card with id 999999999 not found",
			expectedStatusCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.SimClient.SetSimBlocked(
				ctx,
				&pb.SSBRequest{
					Id: int32(tt.id),
				},
			)
			require.Error(t, err)
			require.Equal(t, tt.expectedStatusCode, status.Code(err))
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}

}
