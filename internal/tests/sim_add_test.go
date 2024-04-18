package tests

import (
	"fmt"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/core"
	"simactive/internal/tests/suite"
	"slices"
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
	provider := core.Provider{}.WithName(gofakeit.BS())

	resp, err := st.SimClient.AddSim(
		ctx,
		&pb.AddSimRequest{
			SimData: &pb.AddSimData{
				Number:        phone,
				ProviderName:  provider.Name(),
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
		ID:     resp.Id,
		Number: phone,
		Provider: &pb.ProviderData{
			Id:   int32(provider.Id()),
			Name: provider.Name(),
		},
		IsActivated:   false,
		ActivateUntil: activateUntil,
		IsBlocked:     false,
	}

	r, err := st.SimClient.GetSimList(ctx, &pb.Empty{})
	require.Nil(t, err)

	contains := slices.ContainsFunc(r.SimList, func(s *pb.SimData) bool {
		return s.Number == expectedSimData.Number &&
			s.Provider.Name == expectedSimData.Provider.Name &&
			s.ID == expectedSimData.ID &&
			s.ActivateUntil == expectedSimData.ActivateUntil &&
			s.IsActivated == expectedSimData.IsActivated &&
			s.IsBlocked == expectedSimData.IsBlocked
	})

	assert.True(t, contains)
}

// TestAddSim_DuplicateSim is a test function for adding a duplicate sim.
//
// It tests adding a sim with the same phone number and expects an error to be returned.
// It also checks the error code and message.
func TestAddSim_DuplicateSim(t *testing.T) {

	ctx, ss := suite.NewSuite(t)
	phone := suite.GenerateFakePhoneNumber()
	activateUntil := suite.GenerateFakeDateUnix()
	provider := core.Provider{}.WithName(gofakeit.BS())

	resp, err := ss.SimClient.AddSim(
		ctx,
		&pb.AddSimRequest{
			SimData: &pb.AddSimData{
				Number:        phone,
				ProviderName:  provider.Name(),
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
				ProviderName:  provider.Name(),
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

	provider := gofakeit.BS()
	randomNumber := suite.GenerateFakePhoneNumber()

	tests := []struct {
		name               string
		number             string
		provider           string
		expectedErr        string
		expectedStatusCode codes.Code
	}{
		{
			name:               "Add sim with empty phone number",
			number:             "",
			provider:           provider,
			expectedErr:        "Bad phone number.",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Add sim with invalid phone number",
			number:             "invalid",
			provider:           provider,
			expectedErr:        "Bad phone number.",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Add sim with too large phone number",
			number:             "12312312312312312312312312312312312",
			provider:           provider,
			expectedErr:        "Bad phone number.",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Add sim with empty provider name",
			number:             randomNumber,
			provider:           "",
			expectedErr:        "Provider name is required.",
			expectedStatusCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.SimClient.AddSim(ctx, &pb.AddSimRequest{
				SimData: &pb.AddSimData{
					Number:        tt.number,
					ProviderName:  tt.provider,
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
