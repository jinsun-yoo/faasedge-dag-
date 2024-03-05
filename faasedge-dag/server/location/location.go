package location

type Coord struct {
	x float64 `json:"x"`
	y float64 `json:"y"`
}

type AreaOfInterest struct {
	LowerLeftCoord  Coord
	UpperRightCoord Coord
	Clients         map[int]bool
	Clientlist      []int
}

func (aoi *AreaOfInterest) Belongs(coord Coord) bool {
	if coord.x >= aoi.LowerLeftCoord.x && coord.x <= aoi.UpperRightCoord.x &&
		coord.y >= aoi.LowerLeftCoord.y && coord.y <= aoi.UpperRightCoord.y {
		return true
	} else {
		return false
	}
}

func (aoi *AreaOfInterest) ActiveClients(coord Coord) []int {
	var activeClients []int
	for _, clientID := range aoi.Clientlist {
		if isActive, exists := aoi.Clients[clientID]; exists && isActive {
			activeClients = append(activeClients, clientID)
		}
	}
	return activeClients
}

type LocationTracker struct {
	AreaOfInterestList []AreaOfInterest
}

func (lt *LocationTracker) RegisterLocation(clientId int, coord Coord) *AreaOfInterest {
	var clientAoI *AreaOfInterest

	for i, aoi := range lt.AreaOfInterestList {
		clientInAoi, exists := aoi.Clients[clientId]

		if exists && clientInAoi {
			lt.AreaOfInterestList[i].Clients[clientId] = false
		}

		if aoi.Belongs(coord) && clientAoI == nil {
			if !exists {
				lt.AreaOfInterestList[i].Clientlist = append(lt.AreaOfInterestList[i].Clientlist, clientId)
			}
			lt.AreaOfInterestList[i].Clients[clientId] = true
			clientAoI = &lt.AreaOfInterestList[i]
		}

	}

	return clientAoI
}

func (lt *LocationTracker) LookupAreaOfInterest(coord Coord) AreaOfInterest {
	var aoiFound AreaOfInterest

	for _, aoi := range lt.AreaOfInterestList {

		if aoi.Belongs(coord) {
			return aoi
		}
	}

	return aoiFound
}

func (lt *LocationTracker) RelatedVehicles(clientId int) []int {
	var res []int
	for _, aoi := range lt.AreaOfInterestList {
		clientInAoi, exists := aoi.Clients[clientId]

		if exists && clientInAoi {
			for _, client := range aoi.Clientlist {
				related, _ := aoi.Clients[client]
				if related {
					res = append(res, client)
				}
			}

			return res
		}
	}
	return res
}
