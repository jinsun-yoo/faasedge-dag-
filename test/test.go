package test

import (
	"fmt"
	// "math"
)

type Location struct {
	Latitude  float64
	Longitude float64
}
type Square struct {
	TopLeft     Location
	TopRight    Location
	BottomLeft  Location
	BottomRight Location
}

func IsUserInInnerSquare(square Square, user Location) (bool, string) {
	innerSquareTop    := square.TopLeft.Latitude - 0.1*(square.TopLeft.Latitude - square.BottomLeft.Latitude)
	innerSquareBottom := square.BottomLeft.Latitude + 0.1*(square.TopLeft.Latitude - square.BottomLeft.Latitude)
	innerSquareLeft   := square.TopLeft.Longitude + 0.1*(square.TopRight.Longitude - square.TopLeft.Longitude)
	innerSquareRight  := square.TopRight.Longitude - 0.1*(square.TopRight.Longitude - square.TopLeft.Longitude)

	if user.Latitude >= innerSquareBottom && user.Latitude <= innerSquareTop &&
		user.Longitude >= innerSquareLeft && user.Longitude <= innerSquareRight {
		return true, "User is in the inner square."
	}

	var direction string

	if user.Latitude > innerSquareTop {
		direction = "North"
	} else if user.Latitude < innerSquareBottom {
		direction = "South"
	} else if user.Longitude < innerSquareLeft {
		direction = "West"
	} else if user.Longitude > innerSquareRight {
		direction = "East"
	}

	return false, fmt.Sprintf("User is headed %s.", direction)
}
