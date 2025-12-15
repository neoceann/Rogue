package generation

import (
	"fmt"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
)

// generateRandomSize returns random size of room which fits grid cell
// and respecting:
// a. Min/MaxRoomWidth/Length
// b. MinCorridorLength
func (g *Generator) generateRoomSizeRand() (int, int, error) {
	gridCellWidth := config.GameFieldWidth / config.RoomsInWidth
	gridCellLength := config.GameFieldLength / config.RoomsInLength

	// maximum width within restrictions
	maxWidth := min(gridCellWidth-config.MinCorridorLength,
		config.MaxRoomWidth)
	maxLength := min(gridCellLength-config.MinCorridorLength,
		config.MaxRoomLength)

	// range length of possible dimensions
	widthRange := maxWidth - config.MinRoomWidth + 1
	lengthRange := maxLength - config.MinRoomLength + 1

	// return error if rnage is negative - need to change config setttings
	if widthRange <= 0 || lengthRange <= 0 {
		return 0, 0, fmt.Errorf(
			"invalid room size range: width [%d..%d] (range=%d), "+
				"length [%d..%d] (range=%d)",
			config.MinRoomWidth, maxWidth, widthRange,
			config.MinRoomLength, maxLength, lengthRange,
		)
	}

	// random values within min dimension and max dimension
	randWidth := g.rng.Intn(widthRange) + config.MinRoomWidth
	randLength := g.rng.Intn(lengthRange) + config.MinRoomLength
	return randWidth, randLength, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// centeredRooms computes half-open bounds [pos0, pos1) of a room
// centered in grid cell `ind`
// Returns error if the room would exceed game field bounds
func (g *Generator) centeredRooms(ind int,
	roomWidth, roomLength int,
) (entities.Coordinates, entities.Coordinates, error) {
	gridCellWidth := config.GameFieldWidth / config.RoomsInWidth
	gridCellLength := config.GameFieldLength / config.RoomsInLength
	gridCol := ind % config.RoomsInWidth
	gridRow := ind / config.RoomsInWidth
	// Coordinate of center of room = center of grid cell
	x := gridCol*gridCellWidth + (gridCellWidth-roomWidth)/2
	y := gridRow*gridCellLength + (gridCellLength-roomLength)/2
	// check if room bounds are within game field
	if x < 0 || x+roomWidth > config.GameFieldWidth ||
		y < 0 || y+roomLength > config.GameFieldLength {
		return entities.Coordinates{}, entities.Coordinates{},
			fmt.Errorf("room %d: out of field bounds, func=centeredRooms", ind)
	}
	return entities.Coordinates{X: x, Y: y},
		entities.Coordinates{X: x + roomWidth, Y: y + roomLength},
		nil
}

// dancingRooms (set in config) computes half-open bounds [pos0, pos1)
// of a room randomly shifted in grid cell `ind`
// Returns error if the room would exceed game field bounds
func (g *Generator) dancingRooms(ind int,
	roomWidth, roomLength int,
) (entities.Coordinates, entities.Coordinates, error) {
	gridCellWidth := config.GameFieldWidth / config.RoomsInWidth
	gridCellLength := config.GameFieldLength / config.RoomsInLength
	gridCol := ind % config.RoomsInWidth
	gridRow := ind / config.RoomsInWidth

	// coordinates of upper left point of room ==
	// upper left corner of grid cell + 1/2 MinCorridorLength x, y
	xCorner := gridCol*gridCellWidth + config.MinCorridorLength/2
	yCorner := gridRow*gridCellLength + config.MinCorridorLength/2

	// range of random shift
	xRange := gridCellWidth - config.MinCorridorLength - roomWidth
	yRange := gridCellLength - config.MinCorridorLength - roomLength
	xShift := 0
	yShift := 0
	if xRange > 0 {
		xShift = g.rng.Intn(xRange)
	}
	if yShift > 0 {
		yShift = g.rng.Intn(yRange)
	}
	x := xCorner + xShift
	y := yCorner + yShift
	if x < 0 || x+roomWidth > config.GameFieldWidth ||
		y < 0 || y+roomLength > config.GameFieldLength {
		return entities.Coordinates{}, entities.Coordinates{},
			fmt.Errorf("room %d: out of field bounds, func=dancingRooms", ind)
	}
	return entities.Coordinates{X: x, Y: y},
		entities.Coordinates{X: x + roomWidth, Y: y + roomLength},
		nil
}
