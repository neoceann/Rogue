package generation

import (
	"fmt"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
	"slices"
)

const (
	MaxAttempts = 100
)

// distributes objects evenly and randomly in rooms which are not
// in except list
// returns map of room index and q-ty of objects
func (g *Generator) allocateObjectsInRooms(num int,
	except []int) (map[int]int, error) {
	perRoom := num / (config.RoomsNum - len(except))
	numLeft := num % (config.RoomsNum - len(except))

	// 'rooms[room index]=q-ty' for rooms which are not in except list
	rooms := make(map[int]int)
	// slice of room indexes which are not in except list
	roomIndexes := []int{}
	// spreads equal number for all rooms which are not in except
	for i := range config.RoomsNum {
		if !slices.Contains(except, i) {
			rooms[i] = perRoom
			roomIndexes = append(roomIndexes, i)
		} else {
			rooms[i] = 0
		}
	}
	// spreads left objects randomly among rooms which are not in except
	for range numLeft {
		if len(roomIndexes) < 1 {
			return nil, fmt.Errorf("cannot randomly choose room, "+
				"len(rooms with objects)=%d", len(roomIndexes))
		}
		rnd := g.rng.Intn(len(roomIndexes))
		rooms[roomIndexes[rnd]]++
		roomIndexes = removeFromSlice(roomIndexes, rnd)
	}
	// can be deleted - for self check
	sumCheck := 0
	for _, n := range rooms {
		sumCheck += n
	}
	if sumCheck != num {
		return nil, fmt.Errorf("spread: %d objects, should be %d",
			sumCheck, num)
	}
	return rooms, nil
}

// appoints checked coordinates within restrictions of room: qty
// append objects to Gamesession
func (g *Generator) spreadObjectsInRooms(gs *entities.GameSession,
	objectTemplates []*entities.Item, rooms map[int]int) (int, error) {
	countObjects := 0
	for i, n := range rooms {
		for range n {
			if len(objectTemplates) < 1 {
				return countObjects, fmt.Errorf("len(objectTemplates) < 1, " +
					"cant use rnd")
			}
			objTmp := objectTemplates[g.rng.Intn(len(objectTemplates))]
			newItem, err := entities.NewItem(objTmp.ItemType,
				objTmp.SubType, gs)
			if err != nil {
				return countObjects, fmt.Errorf("random item positioning")
			}
			pos, err := g.attemptToPlaceObject(gs, &gs.Rooms[i], MaxAttempts)
			if err != nil {
				return countObjects, fmt.Errorf("spreadObjectsInRooms: %w", err)
			}
			newItem.Pos = pos
			gs.Items = append(gs.Items, newItem)
			countObjects++
		}
	}
	return countObjects, nil
}

// makes up to MaxAttempt trials to put object on random free cell wiwthin
// room bounds
func (g *Generator) attemptToPlaceObject(gs *entities.GameSession,
	r *entities.Room, maxAttempts int) (entities.Coordinates, error) {
	attempt := 0
	for {
		newPos, err := g.randomCoordinatesInsideRoom1(r)
		if err != nil {
			return entities.Coordinates{},
				fmt.Errorf("could not create random location in the room")
		}
		if isPosFree(gs, newPos) {
			return newPos, nil
		}
		attempt++
		if attempt == maxAttempts {
			return entities.Coordinates{},
				fmt.Errorf("too many attempts (%d) to create random location,"+
					" increase room size", attempt)
		}
	}
}

// common function
func removeFromSlice[T any](slice []T, ind int) []T {
	right := []T{}
	if ind < len(slice)-1 {
		right = slice[ind+1:]
	}
	return append(slice[:ind], right...)
}
