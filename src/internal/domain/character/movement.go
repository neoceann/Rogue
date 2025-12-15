package character

import (
	"math/rand"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
)

const (
	PlayerStep			  = 1
	MonsterBaseStep		  = 1
	OgreStep              = 2
)

func MonsterDirectionHandler(monster *entities.Monster, monsterstep int) (d entities.Direction, x, y int) {
	switch monster.Type {
	case config.Snake:
		d, x, y = snakeMovePattern(monster, monsterstep)

	case config.Ghost:
		x, y = ghostMovePattern(monster)

	default:
		d, x, y = defaultMovePattern(monsterstep)
	}

	return
}

func defaultMovePattern (monsterstep int) (d entities.Direction, x, y int) {
	for d == entities.DirectionInvalid {
		d = entities.Direction(rand.Intn(int(entities.R + 1)))
	}
	switch d {
	case entities.U:
		y-= monsterstep
	case entities.D:
		y+= monsterstep
	case entities.R:
		x+= monsterstep
	case entities.L:
		x-= monsterstep
	}
	return
}

func ghostMovePattern(monster *entities.Monster) (x, y int){
	x = rand.Intn(monster.RoomBounds.Pos1.X + 1 - monster.RoomBounds.Pos0.X) + monster.RoomBounds.Pos0.X
	y = rand.Intn(monster.RoomBounds.Pos1.Y + 1 - monster.RoomBounds.Pos0.Y) + monster.RoomBounds.Pos0.Y
	return
}

func snakeMovePattern(monster *entities.Monster, monsterstep int) (d entities.Direction, x, y int){
	for {
		d = entities.Direction(rand.Intn(int(entities.DR - entities.UL + 1)) + int(entities.UL))
		if d != monster.LastDirection {
			monster.LastDirection = d
			break
		}
	}
	
	switch d {
	case entities.UL:
		x-= monsterstep
		y-= monsterstep
	case entities.UR:
		x+= monsterstep
		y-= monsterstep
	case entities.DL:
		x-= monsterstep
		y+= monsterstep
	case entities.DR:
		x+= monsterstep
		y+= monsterstep
	}

	return
}

func CalcMonsterPosTowardPlayer(monster *entities.Monster, player *entities.Player, m *entities.MapShot) *entities.Coordinates {
    current := &monster.Pos
    
    neighbors := getPassableNeighbors(current, m, player.Backpack)
    if len(neighbors) == 0 {
        return current
    }
    
    bestNeighbor := neighbors[0]
    bestDistance := bestNeighbor.DistTo(player.Pos)
    
    for _, neighbor := range neighbors[1:] {
        distance := neighbor.DistTo(player.Pos)
        if distance < bestDistance {
            bestDistance = distance
            bestNeighbor = neighbor
        }
    }
    
    return &bestNeighbor
}

func getPassableNeighbors(pos *entities.Coordinates, m *entities.MapShot, b entities.Backpack) []entities.Coordinates {
    directions := []entities.Coordinates{
        {X: 1, Y: 0},
        {X: -1, Y: 0},
        {X: 0, Y: 1},
        {X: 0, Y: -1},
    }
    
    var passable []entities.Coordinates
    for _, dir := range directions {
        newPos := entities.Coordinates{X: pos.X + dir.X, Y: pos.Y + dir.Y}
        if can, _ := entities.CanMove(newPos.X, newPos.Y, m, true, b); can {
            passable = append(passable, newPos)
        }
    }
    return passable
}