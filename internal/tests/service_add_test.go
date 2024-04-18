package tests

// TODO
// In adding new sim we have to check that this sim have valid PROVIDER id and name

import (
	"fmt"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"simactive/internal/tests/suite"
	"testing"
)

func TestAddService_HappyPath(t *testing.T) {

	ctx, s := suite.NewSuite(t)

	serviceName := suite.GenerateFakeString(30)

	respAdd, err := s.ServiceClient.AddService(ctx, &pb.AddServiceRequest{Name: serviceName})

	require.NoError(t, err)
	require.NotNil(t, respAdd)
	require.NotEmpty(t, respAdd.GetId())
	assert.Greater(t, int(respAdd.GetId()), 0)

	list, err := s.ServiceClient.GetServiceList(ctx, &pb.Empty{})
	require.NoError(t, err)

	// check that list contains the new service
	found := false
	for _, s := range list.GetServices() {
		if s.Id == respAdd.GetId() {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestAddService_FailCases(t *testing.T) {

	// service already exists
	// service name is empty
	// service name is too long
	ctx, s := suite.NewSuite(t)

	// for case when service already exists we produce firstly adding a service
	var largeServiceName string
	var alreadyExistServiceName string

	// generate random string for service name with length 10-60
	alreadyNameLen := gofakeit.Number(10, 60)
	for {
		if len(largeServiceName) > 64 {
			break
		}

		if len(alreadyExistServiceName) < alreadyNameLen {
			alreadyExistServiceName += gofakeit.Letter()
		}
		largeServiceName += gofakeit.Letter()
	}

	_, err := s.ServiceClient.AddService(ctx, &pb.AddServiceRequest{Name: alreadyExistServiceName})
	require.NoError(t, err)

	tests := []struct {
		name               string
		serviceName        string
		expectedErr        string
		expectedStatusCode codes.Code
	}{
		{
			name:               "Add service with existing name",
			serviceName:        alreadyExistServiceName,
			expectedErr:        fmt.Sprintf("service with name %s already exists", alreadyExistServiceName),
			expectedStatusCode: codes.AlreadyExists,
		},
		{
			name:               "Add service with empty name",
			serviceName:        "",
			expectedErr:        "service name cannot be empty",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Add service with too long name",
			serviceName:        largeServiceName,
			expectedErr:        "service name cannot be longer than 64 characters",
			expectedStatusCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.ServiceClient.AddService(ctx, &pb.AddServiceRequest{Name: tt.serviceName})

			require.Error(t, err)
			require.ErrorContains(t, err, tt.expectedErr)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tt.expectedStatusCode, st.Code())
		})
	}
}
