package types

type Route struct {
	RouteId           string            `json:"route_id" csv:"route_id"`
	AgencyId          string            `json:"agency_id" csv:"agency_id"`
	RouteShortName    string            `json:"route_short_name" csv:"route_short_name"`
	RouteLongName     string            `json:"route_long_name" csv:"route_long_name"`
	RouteDesc         string            `json:"route_desc" csv:"route_desc"`
	RouteType         RouteType         `json:"route_type" csv:"route_type"`
	RouteUrl          string            `json:"route_url" csv:"route_url"`
	RouteColor        string            `json:"route_color" csv:"route_color"`
	RouteTextColor    string            `json:"route_text_color" csv:"route_text_color"`
	RouteSortOrder    int               `json:"route_sort_order" csv:"route_sort_order"`
	ContinuousPickup  ContinuousPickup  `json:"continuous_pickup" csv:"continuous_pickup"`
	ContinuousDropOff ContinuousDropOff `json:"continuous_drop_off" csv:"continuous_drop_off"`
	NetworkId         int               `json:"network_id" csv:"network_id"`
	CEMVSupport       CEMVSupport       `json:"cemv_support" csv:"cemv_support"`
}

type RouteType int

const (
	RouteTypeTram RouteType = iota
	RouteTypeMetro
	RouteTypeRail
	RouteTypeBus
	RouteTypeFerry
	RouteTypeCableTram
	RouteTypeCableCar
	RouteTypeFunicular
	RouteTypeTrolleybus
	RouteTypeMonorail
)

type ContinuousPickup int

const (
	ContinuousPickupAvailable ContinuousPickup = iota
	ContinuousPickupNotAvailable
	ContinuousPickupRequiresContactingAgency
	ContinuousPickupRequiresContactingDriver
)

type ContinuousDropOff int

const (
	ContinuousDropOffAvailable ContinuousDropOff = iota
	ContinuousDropOffNotAvailable
	ContinuousDropOffRequiresContactingAgency
	ContinuousDropOffRequiresContactingDriver
)
