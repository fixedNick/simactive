package tests

import (
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"simactive/internal/core"
	"simactive/internal/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceDelete_HappyPath(t *testing.T) {

	// add new service
	// check that list contains the new service
	// then remove it
	// check that list does not contain the new service

	ctx, s := suite.NewSuite(t)

	service := core.NewService(0, gofakeit.BS())

	addResp, err := s.ServiceClient.AddService(ctx, &pb.AddServiceRequest{Name: service.Name()})
	require.NoError(t, err)
	require.NotNil(t, addResp)
	require.Greater(t, int(addResp.GetId()), 0)

	service.SetID(int(addResp.GetId()))

	list, err := s.ServiceClient.GetServiceList(ctx, &pb.Empty{})
	require.NoError(t, err)

	// check that list contains the new service
	isListContainsNewService := slices.ContainsFunc(list.GetServices(), func(s *pb.ServiceData) bool {
		return int(s.Id) == service.Id()
	})
	require.True(t, isListContainsNewService)

	// remove service
	removedResp, err := s.ServiceClient.DeleteService(ctx, &pb.DeleteServiceRequest{ID: int32(service.Id())})
	require.NoError(t, err)
	assert.Equal(t, int(removedResp.GetId()), service.Id())

	// check that list does not contain the new service
	listAfterDelete, err := s.ServiceClient.GetServiceList(ctx, &pb.Empty{})
	require.NoError(t, err)

	// check that list does not contain the new service
	isListContainsNewServiceAfterDelete := slices.ContainsFunc(listAfterDelete.GetServices(), func(s *pb.ServiceData) bool {
		return int(s.Id) == service.Id()
	})
	require.False(t, isListContainsNewServiceAfterDelete)
}

func TestServiceDelete_FailCases(t *testing.T) {

	ctx, s := suite.NewSuite(t)

	// invalid id
	// service not found

	tests := []struct {
		name               string
		id                 int
		expectedErr        string
		expectedStatusCode codes.Code
	}{
		{
			name:               "Delete service with invalid id",
			id:                 0,
			expectedErr:        "Invalid id, id must be greater than 0",
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name:               "Delete service with not existing id",
			id:                 999999999,
			expectedErr:        "service with id 999999999 not found",
			expectedStatusCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.ServiceClient.DeleteService(ctx, &pb.DeleteServiceRequest{ID: int32(tt.id)})
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tt.expectedStatusCode, st.Code())
		})
	}
}
