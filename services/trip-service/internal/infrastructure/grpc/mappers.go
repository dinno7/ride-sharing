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
		out[i] = rideFareToGrpc(rf)
	}
	return out
}

func rideFareToGrpc(rideFares *domain.RideFare) *pb.RideFare {
	return &pb.RideFare{
		ID:                rideFares.ID,
		UserID:            rideFares.UserID,
		PackageSlug:       rideFares.PackageSlug,
		TotalPriceInCents: rideFares.TotalPriceInCents,
	}
}

func tripToGrpc(domainTrip *domain.Trip) *pb.Trip {
	return &pb.Trip{
		ID: domainTrip.ID,
		// SelectedRideFare: *pb.RideFare,
		Route: *pb.Route,
	}
}
