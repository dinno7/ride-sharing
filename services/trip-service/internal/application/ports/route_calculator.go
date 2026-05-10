package ports

import "github.com/dinno7/ride-sharing/shared/types"

type RouteCalculator interface {
	CalcRoutes(pickup, destination *types.Coordinate) (*types.Route, error)
}
