package main

import (
  "fmt"
	"faasedge-dag/m/v2/test"
)

func main() {
	topLeft := test.Location{Latitude: 33.777450, Longitude: -84.397073} // 33.777450, -84.397073
	topRight := test.Location{Latitude: 33.777379, Longitude: -84.391408} // 33.777379, -84.391408
	bottomLeft := test.Location{Latitude: 33.771600, Longitude: -84.395356} // 33.771600, -84.395356
	bottomRight := test.Location{Latitude: 33.771600, Longitude: -84.390936} // 33.771600, -84.390936

	square := test.Square {
		TopLeft:     topLeft,
		TopRight:    topRight,
		BottomLeft: 	bottomLeft,
		BottomRight: bottomRight,
	}

	userLocation := test.Location{Latitude: 33.772385, Longitude: -84.393640} // 33.772385, -84.393640

	err, result := test.IsUserInInnerSquare(square, userLocation)
	fmt.Println(result)
	fmt.Println(err)
}