package types

import pb "github.com/dinno7/ride-sharing/shared/proto/trip"

type Route struct {
	Distance float64     `json:"distance"`
	Duration float64     `json:"duration"`
	Geometry []*Geometry `json:"geometry"`
}

type Geometry struct {
	Coordinates []*Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (r *Route) ToGrpc() *pb.Route {
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
