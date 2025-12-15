package entities

import "rogue_game/internal/config"

type GameSession struct {
	Level int
	Bounds
	Rooms     []Room
	Corridors []Corridor
	Items     []*Item
	*Player
	Monsters []*Monster
	Doors    []*Door
	Keys     []*Key
	Finish   Coordinates
	LevelDifficultyCoef float64
}

// Player defines player characteristics
type Player struct {
	RoomInd     int
	CorridorInd int
	Pos         Coordinates
	config.Character
	Backpack
	ActiveWeapon *Item
	FightStatus  PlayerFightStatus
	Balance PlayerGameBalance
}

type Item struct {
	Pos          Coordinates
	ItemType     config.ItemType
	SubType      config.ItemSubType
	ID           rune
	Name         string
	Description  string
	Stackable    bool
	MaxStack     int
	StackCounter int
	Effect       config.ItemEffect
}

type GameSignal int

const (
	GsError GameSignal = iota

	GsChoice1
	GsChoice2
	GsChoice3
	GsChoice4
	GsChoice5
	GsChoice6
	GsChoice7
	GsChoice8
	GsChoice9

	GsNoAction
	GsFinish
	GsNextLevel
	GsQuit
	GsUp
	GsDown
	GsLeft
	GsRight
	GsWeapon
	GsFood
	GsPotion
	GsScroll

	GsYes
	GsNo
	GsLose
)

// Defines coordinates on the field
type Coordinates struct {
	X, Y int
}

type MovementChanges struct {
	OldPos     Coordinates
	OldElement any
	NewPos     Coordinates
	NewElement any
}

type MonsterTrick bool

const (
	firstHitToVampire  MonsterTrick = true
	invisibleGhost     MonsterTrick = true
	ogreAttackCooldown MonsterTrick = false
	snakeSleepProc     MonsterTrick = false
	mimicAsItem        MonsterTrick = true
)

// Monster defines type of monster with all characteristics
type Monster struct {
	ID   rune
	Name string
	Type config.MonsterType
	Pos  Coordinates
	config.Character
	Hostility     config.Hostility
	IsChasing     bool
	Trick         MonsterTrick
	FightStatus   BaseFightStatus
	LastDirection Direction
	RoomBounds    Bounds
}

type Backpack map[config.ItemSubType]*Item

// Defines object bounds (rectangle: upper left included and bottom rigth point excluded)
type Bounds struct {
	Pos0 Coordinates // top left point, included
	Pos1 Coordinates // bottom right point, excluded
}

// Defines room Bounds
type Room struct {
	RoomInd int
	Bounds
}

// Defines slice of consequent coordinates
// Delete width
type Corridor struct {
	Points      []Coordinates
	FromRoomInd int
	ToRoomInd   int
}

// Defines map {X, Y} of game session
type MapShot [config.GameFieldWidth][config.GameFieldLength]any

type FogShot struct {
	*MapShot
	RoomNow     *Room
	CorridorNow *Corridor
	Walls       []Coordinates
	Visible     []Coordinates
}

type Field uint8

const (
	RoomFloor Field = iota
	CorridorFloor
	Wall
	Outside
	Finish
)

type BaseFightStatus struct {
	InFight    bool
	MissingHit bool
	DamageDone float64
}

type PlayerFightStatus struct {
	BaseFightStatus
	Asleep            bool
	TargetMonsterName string
}

type PlayerGameBalance struct{
	HealthItemsUsed uint8
}

type DoorColor int

const (
	RedDoor DoorColor = iota
	BlueDoor
	GreenDoor
)

type Door struct {
	Pos   Coordinates
	Color DoorColor
}

type Key struct {
	Pos    Coordinates
	Color  DoorColor
	RoomID int
}