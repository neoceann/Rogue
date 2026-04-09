package config

import "github.com/gdamore/tcell/v2"

// Game balance & world rules — may be adjusted by game designer
const (
	LevelNum         = 21
	FogOfWar         = false
	GameFieldWidth   = 100                          // X
	GameFieldLength  = 40                           // Y
	RoomsInWidth     = 3
	RoomsInLength    = 3
	RoomsNum         = RoomsInWidth * RoomsInLength 
	CenteredGeometry = false                        // true for centered rooms and corridors, false for dancing (diffused)

	// precise dimensions
	MinCorridorLength = 5
	MinRoomWidth      = 10
	MaxRoomWidth      = GameFieldWidth/RoomsInWidth - MinCorridorLength/2
	MinRoomLength     = 5
	MaxRoomLength     = GameFieldLength/RoomsInLength - MinCorridorLength/2

	BalanceAvgHealthItemsUsedTolerance = 2
	BalanceAvgHealthItemsUsedPerLevel = 5
	BalanceLevelDifficultyCoefStep = 0.2
	BalanceCheckPerLevel = 3
)

var (
	DefaultStyle    = tcell.StyleDefault
	StylePlayer     = DefaultStyle.Foreground(tcell.ColorGold)
	StyleZombie     = DefaultStyle.Foreground(tcell.ColorDarkGreen)
	StyleVampire    = DefaultStyle.Foreground(tcell.ColorRed)
	StyleGhost      = DefaultStyle.Foreground(tcell.ColorGhostWhite)
	StyleOgre       = DefaultStyle.Foreground(tcell.ColorYellow)
	StyleSnake      = DefaultStyle.Foreground(tcell.ColorGhostWhite)
	StyleMimic      = DefaultStyle.Foreground(tcell.ColorGhostWhite)
	StyleItemStr    = DefaultStyle.Foreground(tcell.ColorMoccasin)
	StyleItemAgi    = DefaultStyle.Foreground(tcell.ColorLightGreen)
	StyleItemHealth = DefaultStyle.Foreground(tcell.ColorDarkRed)
	StyleKeyRed     = DefaultStyle.Foreground(tcell.ColorDarkRed)
	StyleKeyGreen   = DefaultStyle.Foreground(tcell.ColorForestGreen)
	StyleKeyBlue    = DefaultStyle.Foreground(tcell.ColorBlue)
	StyleWeapon     = DefaultStyle.Foreground(tcell.ColorFireBrick)
)
