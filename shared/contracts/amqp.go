package contracts

import (
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
	"github.com/dinno7/ride-sharing/shared/types"
)

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    []byte `json:"data"`
}

const TripExchange = "trip_exchange"

// Routing keys - using consistent event/command patterns
const (
	// Trip events (trip.event.*)
	TripEventCreated             = "trip.event.created"
	TripEventDriverAssigned      = "trip.event.driver_assigned"
	TripEventNoDriversFound      = "trip.event.no_drivers_found"
	TripEventDriverNotInterested = "trip.event.driver_not_interested"

	// Driver commands (driver.cmd.*)
	DriverCmdTripRequest = "driver.cmd.trip_request"
	DriverCmdTripAccept  = "driver.cmd.trip_accept"
	DriverCmdTripDecline = "driver.cmd.trip_decline"
	DriverCmdLocation    = "driver.cmd.location"
	DriverCmdRegister    = "driver.cmd.register"

	// Payment events (payment.event.*)
	PaymentEventSessionCreated = "payment.event.session_created"
	PaymentEventSuccess        = "payment.event.success"
	PaymentEventFailed         = "payment.event.failed"
	PaymentEventCancelled      = "payment.event.cancelled"

	// Payment commands (payment.cmd.*)
	PaymentCmdCreateSession = "payment.cmd.create_session"
)

type TripCreatedEventData struct {
	Trip *pb.Trip `json:"trip"`
}

type DriverResponseToTripData struct {
	TripID string `json:"tripID"`
	// RiderID string `json:"riderID"`
	Driver struct {
		Id             string           `json:"id"`
		Location       types.Coordinate `json:"location"`
		GeoHash        string           `json:"geohash"`
		Name           string           `json:"name"`
		ProfilePicture string           `json:"profilePicture"`
		PackageSlug    string           `json:"packageSlug"`
		CarPlate       string           `json:"carPlate"`
	} `json:"driver"`
}
