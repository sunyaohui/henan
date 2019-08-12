package vo

type GeoRadius struct {
	Member        []byte        `json:member`
	Distance      float64       `json:distance`
	GeoCoordinate GeoCoordinate `json:geoCoordinate`
}

type GeoCoordinate struct {
	Longitude float64
	Latitude  float64
}
