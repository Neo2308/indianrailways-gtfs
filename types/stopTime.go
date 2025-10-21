package types

type StopTime struct {
	TripId                   string            `json:"trip_id" csv:"trip_id"`
	ArrivalTime              string            `json:"arrival_time" csv:"arrival_time"`
	DepartureTime            string            `json:"departure_time" csv:"departure_time"`
	StopId                   string            `json:"stop_id" csv:"stop_id"`
	LocationGroupId          string            `json:"location_group_id" csv:"location_group_id"`
	LocationId               string            `json:"location_id" csv:"location_id"`
	StopSequence             int               `json:"stop_sequence" csv:"stop_sequence"`
	StopHeadsign             string            `json:"stop_headsign" csv:"stop_headsign"`
	StartPickupDropOffWindow string            `json:"start_pickup_drop_off_window" csv:"start_pickup_drop_off_window"`
	EndPickupDropOffWindow   string            `json:"end_pickup_drop_off_window" csv:"end_pickup_drop_off_window"`
	PickupType               PickupType        `json:"pickup_type" csv:"pickup_type"`
	DropOffType              DropOffType       `json:"drop_off_type" csv:"drop_off_type"`
	ContinuousPickup         ContinuousPickup  `json:"continuous_pickup" csv:"continuous_pickup"`
	ContinuousDropOff        ContinuousDropOff `json:"continuous_drop_off" csv:"continuous_drop_off"`
	ShapeDistTraveled        float64           `json:"shape_dist_traveled" csv:"shape_dist_traveled"`
	Timepoint                Timepoint         `json:"timepoint" csv:"timepoint"`
	PickupBookingRuleId      string            `json:"pickup_booking_rule_id" csv:"pickup_booking_rule_id"`
	DropOffBookingRuleId     string            `json:"drop_off_booking_rule_id" csv:"drop_off_booking_rule_id"`
}

type Timepoint int

const (
	// TimepointApproximateTimes Times are considered approximate.
	TimepointApproximateTimes Timepoint = iota

	// TimepointExactTimes Times are considered exact.
	TimepointExactTimes
)

type PickupType int

const (
	// PickupTypeRegularlyScheduled Regularly scheduled pickup.
	PickupTypeRegularlyScheduled PickupType = iota

	// PickupTypeNoPickupAvailable No pickup available.
	PickupTypeNoPickupAvailable

	// PickupTypeRequiresContactingAgency Must phone agency to arrange pickup.
	PickupTypeRequiresContactingAgency

	// PickupTypeRequiresContactingDriver Must coordinate with driver to arrange pickup.
	PickupTypeRequiresContactingDriver
)

type DropOffType int

const (
	// DropOffTypeRegularlyScheduled Regularly scheduled drop off.
	DropOffTypeRegularlyScheduled DropOffType = iota

	// DropOffTypeNoDropOffAvailable No drop off available.
	DropOffTypeNoDropOffAvailable

	// DropOffTypeRequiresContactingAgency Must phone agency to arrange drop off.
	DropOffTypeRequiresContactingAgency

	// DropOffTypeRequiresContactingDriver Must coordinate with driver to arrange drop off.
	DropOffTypeRequiresContactingDriver
)
