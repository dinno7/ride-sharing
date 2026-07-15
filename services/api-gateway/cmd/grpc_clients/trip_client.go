package grpcclients

import (
	"github.com/dinno7/ride-sharing/shared/env"
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
	"github.com/dinno7/ride-sharing/shared/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tripGrpcServiceClient struct {
	conn   *grpc.ClientConn
	Client pb.TripServiceClient
}

func NewTripServiceClient() (*tripGrpcServiceClient, error) {
	tripServiceURL := env.GetString("TRIP_SERVICE_URL", "trip-service:9000")

	dialOptions := append(
		tracing.GRPCClientTracingOpts(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	conn, err := grpc.NewClient(
		tripServiceURL,
		dialOptions...,
	)
	if err != nil {
		return nil, err
	}

	return &tripGrpcServiceClient{
		conn:   conn,
		Client: pb.NewTripServiceClient(conn),
	}, nil
}

func (tc *tripGrpcServiceClient) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}
