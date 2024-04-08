package tests

import (
	"fmt"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddSim_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	phone := generateFakePhoneNumber()
	activateUntil := gofakeit.Date().Unix()

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

func generateFakePhoneNumber() string {
	countryCode := gofakeit.Number(1, 999)
	operatorCode := gofakeit.Number(100, 999)
	phoneNumber := gofakeit.Number(1000, 9999)

	return fmt.Sprintf("+%d%d%d", countryCode, operatorCode, phoneNumber)
}
