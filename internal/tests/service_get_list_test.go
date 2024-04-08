package tests

import (
	"simactive/internal/tests/suite"
	"testing"

	"github.com/stretchr/testify/require"
)

func ServiceGetListTest_HappyPath(t *testing.T) {
	ctx, s := suite.NewSuite(t)

	list, err := s.ServiceClient.GetServiceList(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, list)
}
