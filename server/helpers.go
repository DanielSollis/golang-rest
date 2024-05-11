package server

import "math"

// Calculates the great circle distance between two points.
// The haversine formula assumes points on a perfect sphere
// (the earth isn't a perfect sphere) so the haversine error
// can be up to 0.5%
func haversine(user, sensor *Coordinates) float64 {
	userLat, sensorLat := user.Latitude, sensor.Latitude
	userLon, sensorLon := user.Longitude, sensor.Longitude

	// Distance between latitudes in radians
	latDistanceRad := (sensorLat - userLat) * math.Pi / 180
	lonDistanceRad := (sensorLon - userLon) * math.Pi / 180

	// Latitudes in radians
	userLatRad := userLat * math.Pi / 180
	sensorLatRad := sensorLat * math.Pi / 180

	// Calculate the square of half the chord length between two points 'a'
	latPower := math.Pow(math.Sin(latDistanceRad/2), 2)
	lonPower := math.Pow(math.Sin(lonDistanceRad/2), 2)
	latCosine := math.Cos(userLatRad) * math.Cos(sensorLatRad)
	a := latPower + lonPower*latCosine

	// Calculate the angular between the two points 'c'
	c := 2 * math.Asin(math.Sqrt(a))

	earthsRadius := 6371 // In km (approximately)
	return float64(earthsRadius) * c
}

func CreateSensor(name, unit string, lat, lon float64) *Sensor {
	return &Sensor{
		Name: name,
		Location: Coordinates{
			Latitude:  lat,
			Longitude: lon,
		},
		Tags: SensorTags{
			Name: name,
			Unit: unit,
		},
	}
}
