package types

type Trip struct {
	RouteId              string                  `json:"route_id" csv:"route_id"`
	ServiceId            string                  `json:"service_id" csv:"service_id"`
	TripId               string                  `json:"trip_id" csv:"trip_id"`
	TripShortName        string                  `json:"trip_short_name" csv:"trip_short_name"`
	TripHeadsign         string                  `json:"trip_headsign" csv:"trip_headsign"`
	DirectionId          DirectionId             `json:"direction_id" csv:"direction_id"`
	BlockId              string                  `json:"block_id" csv:"block_id"`
	ShapeId              string                  `json:"shape_id" csv:"shape_id"`
	WheelchairAccessible WheelChairAccessibility `json:"wheelchair_accessible" csv:"wheelchair_accessible"`
	BikesAllowed         BikesAllowed            `json:"bikes_allowed" csv:"bikes_allowed"`
	CarsAllowed          CarsAllowed             `json:"cars_allowed" csv:"cars_allowed"`
}

type DirectionId int

const (
	// DirectionIdOutbound Travel in one direction (i.e. from source to destination)
	DirectionIdOutbound DirectionId = iota

	// DirectionIdInbound Travel in the opposite direction (i.e. from destination to source)
	DirectionIdInbound
)

type WheelChairAccessibility int

const (
	// WheelChairAccessibilityNoInfo No accessibility information for the trip.
	WheelChairAccessibilityNoInfo WheelChairAccessibility = iota

	// WheelChairAccessibilityAccessible Vehicle being used on this particular trip can accommodate at least one rider in a wheelchair
	WheelChairAccessibilityAccessible

	// WheelChairAccessibilityNotAccessible No riders in wheelchairs can be accommodated on this trip.
	WheelChairAccessibilityNotAccessible
)

type BikesAllowed int

const (
	// BikesAllowedNoInfo No bike information for the trip.
	BikesAllowedNoInfo BikesAllowed = iota

	// BikesAllowedAllowed Vehicle being used on this particular trip can accommodate at least one bicycle.
	BikesAllowedAllowed

	// BikesAllowedNotAllowed No bicycles are allowed on this trip.
	BikesAllowedNotAllowed
)

type CarsAllowed int

const (
	// CarsAllowedNoInfo No car information for the trip.
	CarsAllowedNoInfo CarsAllowed = iota

	// CarsAllowedAllowed Vehicle being used on this particular trip can accommodate at least one car.
	CarsAllowedAllowed

	// CarsAllowedNotAllowed No cars are allowed on this trip.
	CarsAllowedNotAllowed
)
