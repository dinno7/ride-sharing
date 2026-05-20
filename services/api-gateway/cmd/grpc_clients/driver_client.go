package grpcclients

import (
	"github.com/dinno7/ride-sharing/shared/env"
	pb "github.com/dinno7/ride-sharing/shared/proto/driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type driverGrpcServiceClient struct {
	conn   *grpc.ClientConn
	Client pb.DriverServiceClient
}

func NewDriverServiceClient() (*driverGrpcServiceClient, error) {
	driverServiceURL := env.GetString("DRIVER_SERVICE_URL", "driver-service:9000")

	conn, err := grpc.NewClient(
		driverServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &driverGrpcServiceClient{
		conn:   conn,
		Client: pb.NewDriverServiceClient(conn),
	}, nil
}

func (tc *driverGrpcServiceClient) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}
