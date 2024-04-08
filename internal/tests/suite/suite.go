package suite

import (
	"context"
	"fmt"
	"net"
	"simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/config"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg       config.Config
	SimClient SimHelper.SimClient
}

const (
	grpcHost   = "127.0.0.1"
	configPath = "../../config/app/local.yaml"
)

func NewSuite(t *testing.T) (context.Context, *Suite) {

	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath(configPath)
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		grpcArrdress(*cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failer: %v", err)
	}

	return ctx, &Suite{
		T:         t,
		Cfg:       *cfg,
		SimClient: SimHelper.NewSimClient(cc),
	}
}

func grpcArrdress(cfg config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}

func GenerateFakePhoneNumber() string {
	countryCode := gofakeit.Number(1, 999)
	operatorCode := gofakeit.Number(100, 999)
	phoneNumber := gofakeit.Number(1000, 9999)

	return fmt.Sprintf("%d%d%d", countryCode, operatorCode, phoneNumber)
}

func GenerateFakeDateUnix() int64 {
	return int64(gofakeit.Number(int(time.Now().Unix()), int(time.Now().Unix()+int64((time.Hour*24*365).Seconds()))))
}
