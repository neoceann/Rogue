package entities

import (
	"fmt"
	"rogue_game/internal/config"
)

// Coordinates: returns Manhattan distance between 2 points
func (c1 Coordinates) DistTo(c2 Coordinates) int {
	dX := c1.X - c2.X
	dY := c1.Y - c2.Y
	if dX < 0 {
		dX = -dX
	}
	if dY < 0 {
		dY = -dY
	}
	return dX + dY

}

func (c Coordinates) FindFirstEmptyAround(m *MapShot, b Backpack) (int, int, bool) {
	var ex, ey int = -1, -1
	correct := true
	for x := -1; x < 2; x++ {
        for y := -1; y < 2; y++ {
			cm,_ := CanMove(c.X + x,c.Y + y, m, false, b)
			if cm && m[c.X + x][c.Y + y] == RoomFloor {
				ex = c.X + x
				ey = c.Y + y
			}
		}
	}

	correct = ex != -1 && ey != -1

	return ex, ey, correct
}

// Direction: represents direction of object
// LU = left up ↖, ...
type Direction int

const (
	DirectionInvalid Direction = iota // 0
	U                                 // ↑ (or Forward)
	D                                 // ↓
	L                                 // ←
	R                                 // →
	UL                                // ↖
	UR                                // ↗
	DL                                // ↙
	DR                                // ↘
	Stop                              // no movement
)

func MoveUnit(d Direction, ghostx, ghosty int, m *MapShot, gs *GameSession, unit any, step int, mc *MovementChanges) {
	var uc *Coordinates
	isGhost := false
	switch v := unit.(type){
	case *Player:
		uc = &gs.Player.Pos
	case *Monster:
		monsterIndex := GetMonsterIndexByPointer(gs, v)
		uc = &gs.Monsters[monsterIndex].Pos
		isGhost = v.Type == config.Ghost || v.IsChasing
	default:
		panic(fmt.Errorf("unknown type for movement"))
	}

	mc.NewElement = unit
	mc.OldPos = *uc
	mc.OldElement = RoomFloor

	if isGhost {
		uc.X = ghostx
		uc.Y = ghosty
	} else {
		updateCoordinatesWithStep(d, uc, step)
	}

	mc.NewPos = *uc
}

func updateCoordinatesWithStep (d Direction, uc *Coordinates, step int) {
	switch d {
	case U:
		uc.Y -= step
	case D:
		uc.Y += step
	case L:
		uc.X -= step
	case R:
		uc.X += step
	case UL:
		uc.X -= step
		uc.Y -= step
	case UR:
		uc.X += step
		uc.Y -= step
	case DL:
		uc.X -= step
		uc.Y += step
	case DR:
		uc.X += step
		uc.Y += step
	}
}

func CanMove(x, y int, m *MapShot, triggeredByMonster bool, b Backpack) (bool, bool) {
	switch v := m[x][y].(type) {
	case *Monster:
		return false, !triggeredByMonster
	case Field:
		 return (v != Wall && v != Outside && (v != Finish || !triggeredByMonster)), false
	case *Item:
		return !triggeredByMonster, false
	case *Player:
		return false, false
	case *Door:
		return !triggeredByMonster && b.HaveKeyForDoor(v), false
	}

	return true, false
}