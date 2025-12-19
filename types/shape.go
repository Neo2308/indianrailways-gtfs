package types

type Shape struct {
	ShapeId           string  `json:"shape_id" csv:"shape_id"`
	ShapePtLat        float64 `json:"shape_pt_lat" csv:"shape_pt_lat"`
	ShapePtLon        float64 `json:"shape_pt_lon" csv:"shape_pt_lon"`
	ShapePtSequence   int     `json:"shape_pt_sequence" csv:"shape_pt_sequence"`
	ShapeDistTraveled float64 `json:"shape_dist_traveled" csv:"shape_dist_traveled"`
}
