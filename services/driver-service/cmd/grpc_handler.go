package main

import (
	"context"

	pb "github.com/dinno7/ride-sharing/shared/proto/driver"
	"google.golang.org/grpc"
)

type driverGrpcHandler struct {
	driverService *DriverService
	pb.UnimplementedDriverServiceServer
}

func NewDriverGrpcHandler(server *grpc.Server, driverService *DriverService) *driverGrpcHandler {
	handler := &driverGrpcHandler{driverService: driverService}

	pb.RegisterDriverServiceServer(server, handler)
	return handler
}

func (h *driverGrpcHandler) RegisterDriver(
	ctx context.Context,
	req *pb.DriverRequest,
) (*pb.DriverResponse, error) {
	driver, err := h.driverService.RegisterDriver(req.GetDriverId(), req.PackageSlug)
	if err != nil {
		return nil, err
	}
	return &pb.DriverResponse{
		Driver: driver,
	}, nil
}

func (h *driverGrpcHandler) UnregisterDriver(
	ctx context.Context,
	req *pb.DriverRequest,
) (*pb.DriverResponse, error) {
	h.driverService.UnregisterDriver(req.GetDriverId())
	return &pb.DriverResponse{
		Driver: nil,
	}, nil
}
