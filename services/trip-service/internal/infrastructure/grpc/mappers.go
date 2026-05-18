package grpc

import (
	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
	"github.com/dinno7/ride-sharing/shared/types"
)

func tripPreviewToGrpc(i *ports.PreviewTripOutput) *pb.PreviewTripResponse {
	return &pb.PreviewTripResponse{
		Route:     routeToGrpc(i.Route),
		RideFares: rideFaresToGrpc(i.RideFares),
	}
}

func routeToGrpc(r *types.Route) *pb.Route {
	geometry := []*pb.Geometry{}
	for _, g := range r.Geometry {
		cordinates := []*pb.Cordinate{}
		for _, c := range g.Coordinates {
			cordinates = append(cordinates, &pb.Cordinate{
				Latitude:  c.Latitude,
				Longitude: c.Longitude,
			})
		}
		newGeo := &pb.Geometry{Cordinates: cordinates}
		geometry = append(geometry, newGeo)
	}
	return &pb.Route{
		Geometry: geometry,
		Distance: r.Distance,
		Duration: r.Duration,
	}
}

func rideFaresToGrpc(rideFares []*domain.RideFare) []*pb.RideFare {
	out := make([]*pb.RideFare, len(rideFares))
	for i, rf := range rideFares {
		out[i] = &pb.RideFare{
			ID:                rf.ID,
			UserID:            rf.UserID,
			PackageSlug:       rf.PackageSlug,
			TotalPriceInCents: rf.TotalPriceInCents,
		}
	}
	return out
}
