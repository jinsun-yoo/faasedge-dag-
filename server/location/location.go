package location

type AreaOfInterest struct {
	lowerLeftCoord Coord
	upperRightCoord Coord
	clients map[int]bool{}
	clientlist []int 

	Belongs(coord Coord) bool 
	Clients() []int 
}


type LocationTracker struct {
	areaOfInterestList []AreaOfInterest
	

	RegisterLocation (clientId int, coord Coord)
	LookupAreaOfInterest (coord Coord) AreaOfInterst
	RelatedVehicles (clientId int) 
}
