package location

type Coord struct {
	x float64
	y float64
}

type AreaOfInterest struct {
	lowerLeftCoord  Coord
	upperRightCoord Coord
	clients         map[int]bool
	clientlist      []int
}

func (aoi *AreaOfInterest) Belongs(coord Coord) bool {
	if coord.x >= aoi.lowerLeftCoord.x && coord.x <= aoi.upperRightCoord.x &&
		coord.y >= aoi.lowerLeftCoord.y && coord.y <= aoi.upperRightCoord.y {
		return true
	} else {
		return false
	}
}

func (aoi *AreaOfInterest) Clients(coord Coord) []int {
	var activeClients []int
	for _, clientID := range aoi.clientlist {
		if isActive, exists := aoi.clients[clientID]; exists && isActive {
			activeClients = append(activeClients, clientID)
		}
	}
	return activeClients
}

type LocationTracker struct {
	areaOfInterestList []AreaOfInterest
}

func (lt *LocationTracker) RegisterLocation(clientId int, coord Coord) *AreaOfInterest {
	var clientAoI *AreaOfInterest

	for i, aoi := range lt.areaOfInterestList {
		clientInAoi, exists := aoi.clients[clientId]

		if exists && clientInAoi {
			lt.areaOfInterestList[i].clients[clientId] = false
		}

		if aoi.Belongs(coord) && clientAoI == nil {
			if !exists {
				lt.areaOfInterestList[i].clientlist = append(lt.areaOfInterestList[i].clientlist, clientId)
			}
			lt.areaOfInterestList[i].clients[clientId] = true
			clientAoI = &lt.areaOfInterestList[i]
		}

	}

	return clientAoI
}

func (lt *LocationTracker) LookupAreaOfInterest(coord Coord) AreaOfInterest {
	var aoiFound AreaOfInterest

	for _, aoi := range lt.areaOfInterestList {

		if aoi.Belongs(coord) {
			return aoi
		}
	}

	return aoiFound
}

func (lt *LocationTracker) RelatedVehicles(clientId int) []int {
	var res []int
	for _, aoi := range lt.areaOfInterestList {
		clientInAoi, exists := aoi.clients[clientId]

		if exists && clientInAoi {
			for _, client := range aoi.clientlist {
				related, _ := aoi.clients[client]
				if related {
					res = append(res, client)
				}
			}

			return res
		}
	}
	return res
}
