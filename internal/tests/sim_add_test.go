package tests

import (
	"fmt"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/tests/suite"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAddSim_HappyPath is a test function for adding a sim card in a happy path scenario.
//
// It uses the testing context and suite, generates a fake phone number and date, then adds a sim card with the provided details.
// It also checks the response and expected sim data.
func TestAddSim_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	phone := suite.GenerateFakePhoneNumber()
	activateUntil := suite.GenerateFakeDateUnix()

	resp, err := st.SimClient.AddSim(
		ctx,
		&pb.AddSimRequest{
			SimData: &pb.AddSimData{
				Number:        phone,
				ProviderID:    1,
				IsActivated:   false,
				ActivateUntil: activateUntil,
				IsBlocked:     false,
			},
		},
	)

	require.Nil(t, err)
	assert.Greater(t, int(resp.Id), 0)
	assert.NotEmpty(t, resp.Message)

	expectedSimData := pb.SimData{
		ID:            resp.Id,
		Number:        phone,
		ProviderID:    1,
		IsActivated:   false,
		ActivateUntil: activateUntil,
		IsBlocked:     false,
	}

	r, err := st.SimClient.GetSimList(ctx, &pb.Empty{})
	require.Nil(t, err)
	assert.Contains(t, r.SimList, &expectedSimData)
}

// TestAddSim_DuplicateSim is a test function for adding a duplicate sim.
//
// It tests adding a sim with the same phone number and expects an error to be returned.
// It also checks the error code and message.
func TestAddSim_DuplicateSim(t *testing.T) {

	ctx, ss := suite.NewSuite(t)
	phone := suite.GenerateFakePhoneNumber()
	activateUntil := suite.GenerateFakeDateUnix()

	resp, err := ss.SimClient.AddSim(
		ctx,
		&pb.AddSimRequest{
			SimData: &pb.AddSimData{
				Number:        phone,
				ProviderID:    1,
				IsActivated:   false,
				ActivateUntil: activateUntil,
				IsBlocked:     false,
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetId())

	resp, err = ss.SimClient.AddSim(
		ctx,
		&pb.AddSimRequest{
			SimData: &pb.AddSimData{
				Number:        phone,
				ProviderID:    1,
				IsActivated:   false,
				ActivateUntil: activateUntil,
				IsBlocked:     false,
			},
		},
	)
	require.Error(t, err)

	// check status
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())

	// check response
	assert.Empty(t, resp.GetId())
	assert.ErrorContains(t, err, fmt.Sprintf("sim card with number %s already exists", phone))
}

// TestAddSim_FailCases tests the failure cases of the AddSim function.
//
// It tests various scenarios where the AddSim function should return an error.
// Also it checks the error code and message.
func TestAddSim_FailCases(t *testing.T) {

	ctx, s := suite.NewSuite(t)

	provider := gofakeit.Number(1, 999)
	randomNumber := suite.GenerateFakePhoneNumber()

	tests := []struct {
		name               string
		number             string
		providerID         int
		expectedErr        string
		expectedStatusCode codes.Code
	}{
		{
			name:               "Add sim with 2empty phone number",
			number:             "",
			providerID:         provider,
			expectedErr:        "Bad phone number.",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Add sim with invalid phone number",
			number:             "invalid",
			providerID:         provider,
			expectedErr:        "Bad phone number.",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Add sim with too large phone number",
			number:             "12312312312312312312312312312312312",
			providerID:         provider,
			expectedErr:        "Bad phone number.",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Add sim with empty provider id",
			number:             randomNumber,
			providerID:         0,
			expectedErr:        "Bad provider id.",
			expectedStatusCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.SimClient.AddSim(ctx, &pb.AddSimRequest{
				SimData: &pb.AddSimData{
					Number:        tt.number,
					ProviderID:    int32(tt.providerID),
					IsActivated:   gofakeit.Bool(),
					ActivateUntil: suite.GenerateFakeDateUnix(),
					IsBlocked:     gofakeit.Bool(),
				},
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedStatusCode, st.Code())
		})
	}
}
