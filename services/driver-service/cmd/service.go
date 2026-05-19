package main

import (
	math "math/rand/v2"
	"sync"

	pb "github.com/dinno7/ride-sharing/shared/proto/driver"
	"github.com/dinno7/ride-sharing/shared/util"

	"github.com/mmcloughlin/geohash"
)

type driverInMap struct {
	Driver *pb.Driver
	// Index int
	// TODO: route
}

type DriverService struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

func NewDriverService() *DriverService {
	return &DriverService{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *DriverService) RegisterDriver(driverId string, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIndex]

	// we can ignore this property for now, but it must be sent to the frontend.
	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	randomAvatar := util.GetRandomAvatar(1)
	randomPlate := GenerateRandomPlate()

	driver := &pb.Driver{
		Geohash:        geohash,
		Location:       &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
		Name:           "Lando Norris",
		Id:             driverId,
		PackageSlug:    packageSlug,
		ProfilePicture: randomAvatar,
		CarPlate:       randomPlate,
	}

	s.drivers = append(s.drivers, &driverInMap{
		Driver: driver,
	})

	return driver, nil
}

func (s *DriverService) UnregisterDriver(driverId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, driver := range s.drivers {
		if driver.Driver.Id == driverId {
			s.drivers = append(s.drivers[:i], s.drivers[i+1:]...)
		}
	}
}
