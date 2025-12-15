package generation

import (
	"math/rand"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
)

func (g *Generator) DoorsKeysGeneration(gs *entities.GameSession){//, m *entities.MapShot) {

	doors, keys := generateDoorsAndKeys(gs, g)

	gs.Doors = append(gs.Doors, doors...)
	gs.Keys = append(gs.Keys, keys...)

	for _, key := range keys {
		var sbtype config.ItemSubType
		switch key.Color {
		case entities.RedDoor:
			sbtype = config.ItemSubTypeKeyRed
		case entities.BlueDoor:
			sbtype = config.ItemSubTypeKeyBlue
		case entities.GreenDoor:
			sbtype = config.ItemSubTypeKeyGreen
		}
		keyAsItem, err := entities.NewItem(config.ItemTypeKey, sbtype, gs)
		if err != nil {
			panic("bad key")
		}
		keyAsItem.Pos = entities.Coordinates{X: key.Pos.X, Y: key.Pos.Y}
		gs.Items = append(gs.Items, keyAsItem)
	}
}

func getConnectedCorridors(corridors []entities.Corridor, roomID int) []entities.Corridor {
	var result []entities.Corridor
	for _, c := range corridors {
		if c.FromRoomInd == roomID || c.ToRoomInd == roomID {
			result = append(result, c)
		}
	}
	return result
}

func getOtherRoom(corridor entities.Corridor, roomID int) int {
	if corridor.FromRoomInd == roomID {
		return corridor.ToRoomInd
	}
	return corridor.FromRoomInd
}

func findDoorInCorridor(doors []*entities.Door, corridor entities.Corridor) *entities.Door {
	for _, door := range doors {
		for _, point := range corridor.Points {
			if point.X == door.Pos.X && point.Y == door.Pos.Y {
				return door
			}
		}
	}
	return nil
}

func getDoorTargetRoom(door *entities.Door, corridors []entities.Corridor) int {
	for _, corr := range corridors {
		for _, point := range corr.Points {
			if point.X == door.Pos.X && point.Y == door.Pos.Y {
				return corr.ToRoomInd
			}
		}
	}
	return -1
}

func generateDoorsAndKeys(gs *entities.GameSession, g *Generator) ([]*entities.Door, []*entities.Key) {
	for attempt := 0; attempt < 10; attempt++ {
		doors, keys := tryGenerate(gs, g)
		if validateAccessibility(gs.Rooms, gs.Corridors, doors, keys, gs.Player.RoomInd) {
			return doors, keys
		}
	}
	return nil, nil
}

func tryGenerate(gs *entities.GameSession, g *Generator) ([]*entities.Door, []*entities.Key) {
	var doors []*entities.Door
	var keys []*entities.Key
	var keyPos entities.Coordinates

	availableColors := []entities.DoorColor{
		entities.RedDoor,
		entities.BlueDoor,
		entities.GreenDoor,
	}

	doorCount := rand.Intn(3) + 1
	for i := 0; i < doorCount && i < len(gs.Corridors); i++ {
		colorIndex := rand.Intn(len(availableColors))
		color := availableColors[colorIndex]

		availableColors = append(availableColors[:colorIndex], availableColors[colorIndex+1:]...)

		corr := gs.Corridors[rand.Intn(len(gs.Corridors))]
		if len(corr.Points) > 2 {
			pos := corr.Points[rand.Intn(len(corr.Points)-2)+1]
			doors = append(doors, &entities.Door{Pos: entities.Coordinates{X: pos.X, Y: pos.Y}, Color: color})
		}
	}

	for _, door := range doors {
		targetRoom := getDoorTargetRoom(door, gs.Corridors)
		var availableRooms []entities.Room
		for _, room := range gs.Rooms {
			if room.RoomInd != targetRoom {
				availableRooms = append(availableRooms, room)
			}
		}

		if len(availableRooms) > 0 {
			room := availableRooms[rand.Intn(len(availableRooms))]

			attempt := 0
			maxAttempts := 100
			for {
				newPos, _ := g.randomCoordinatesInsideRoom(gs, room.RoomInd)

				if isPosFree(gs, newPos) {
					keyPos = newPos
					break
				}
				attempt++
				if attempt == maxAttempts {
					panic("keys error")
				}
			}
			keys = append(keys, &entities.Key{Pos: entities.Coordinates{X: keyPos.X, Y: keyPos.Y}, Color: door.Color, RoomID: room.RoomInd})
		}
	}

	return doors, keys
}

func validateAccessibility(rooms []entities.Room, corridors []entities.Corridor, doors []*entities.Door, keys []*entities.Key, startRoomID int) bool {
	if len(doors) == 0 {
		return true
	}

	neededKeys := make(map[entities.DoorColor]bool)
	for _, door := range doors {
		neededKeys[door.Color] = true
	}

	visited := make(map[int]bool)
	collectedKeys := make(map[entities.DoorColor]bool)
	queue := []int{startRoomID}
	visited[startRoomID] = true

	for len(queue) > 0 {
		roomID := queue[0]
		queue = queue[1:]

		for _, key := range keys {
			if key.RoomID == roomID {
				collectedKeys[key.Color] = true
			}
		}

		for _, corridor := range getConnectedCorridors(corridors, roomID) {
			nextRoom := getOtherRoom(corridor, roomID)
			if visited[nextRoom] {
				continue
			}
			if nextRoom < 0 || nextRoom >= len(rooms) {
				continue
			}

			if door := findDoorInCorridor(doors, corridor); door != nil {
				if !collectedKeys[door.Color] {
					continue
				}
			}

			visited[nextRoom] = true
			queue = append(queue, nextRoom)
		}
	}

	for color := range neededKeys {
		if !collectedKeys[color] {
			return false
		}
	}

	return true
}
