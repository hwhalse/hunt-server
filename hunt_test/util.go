package main

import (
	"math"
	"math/rand"
)

func randomUSCoordinate() (lat, lon float64) {
	const (
		minLat = 24.396308
		maxLat = 49.384358
		minLon = -125.0
		maxLon = -66.93457
	)
	return minLat + rand.Float64()*(maxLat-minLat),
		minLon + rand.Float64()*(maxLon-minLon)
}

func randomNearbyCoordinate(lat, lon float64) (newLat, newLon float64) {
	const earthRadius = 6371000
	const maxDistance = 10000

	distance := rand.Float64() * maxDistance
	bearing := rand.Float64() * 2 * math.Pi

	latRad := lat * math.Pi / 180
	lonRad := lon * math.Pi / 180

	newLatRad := math.Asin(math.Sin(latRad)*math.Cos(distance/earthRadius) +
		math.Cos(latRad)*math.Sin(distance/earthRadius)*math.Cos(bearing))
	newLonRad := lonRad + math.Atan2(
		math.Sin(bearing)*math.Sin(distance/earthRadius)*math.Cos(latRad),
		math.Cos(distance/earthRadius)-math.Sin(latRad)*math.Sin(newLatRad))

	return newLatRad * 180 / math.Pi, newLonRad * 180 / math.Pi
}