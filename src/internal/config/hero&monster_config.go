package config

import "github.com/gdamore/tcell/v2"

const (
	StartLevelMonstersPerRoom = 0.7
	LastLevelMonstersPerRoom  = 5
)

const (
	BaseHitChance = 0.8
	MinHitChance  = 0.2
	MaxHitChance  = 0.9
	AgilityImpact = 0.1

	SnakeSleepChance            = 0.5
	VamprireHitMaxHealthPercent = 10

	StatsIncreasePerLvl      = 5
	ItemsStatsIncreasePerLvl = StatsIncreasePerLvl / 3.0

	VeryHighBaseHealth = 200.0
	HighBaseHealth     = 150.0
	AverageBaseHealth  = 100.0
	LowBaseHealth      = 50.0

	VeryHighBaseStrength = 20.0
	HighBaseStrength     = 10.0
	AverageBaseStrength  = 5.0
	LowBaseStrength      = 3.0

	VeryHighBaseAgility = 20.0
	HighBaseAgility     = 10.0
	AverageBaseAgility  = 5.0
	LowBaseAgility      = 3.0
)

const (
	PlayerBaseHealth    = HighBaseHealth
	PlayerBaseMaxHealth = HighBaseHealth
	PlayerBaseAgility   = VeryHighBaseAgility
	PlayerBaseStrength  = HighBaseStrength

	ZombieBaseHealth    = HighBaseHealth
	ZombieBaseMaxHealth = HighBaseHealth
	ZombieBaseAgility   = LowBaseAgility
	ZombieBaseStrength  = AverageBaseStrength

	VampireBaseHealth    = HighBaseHealth
	VampireBaseMaxHealth = HighBaseHealth
	VampireBaseAgility   = HighBaseAgility
	VampireBaseStrength  = AverageBaseStrength

	GhostBaseHealth    = LowBaseHealth
	GhostBaseMaxHealth = LowBaseHealth
	GhostBaseAgility   = HighBaseAgility
	GhostBaseStrength  = LowBaseStrength

	OgreBaseHealth    = VeryHighBaseHealth
	OgreBaseMaxHealth = VeryHighBaseHealth
	OgreBaseAgility   = LowBaseAgility
	OgreBaseStrength  = VeryHighBaseStrength

	SnakeBaseHealth    = AverageBaseHealth
	SnakeBaseMaxHealth = AverageBaseHealth
	SnakeBaseAgility   = VeryHighBaseAgility
	SnakeBaseStrength  = AverageBaseStrength

	MimicBaseHealth    = HighBaseHealth
	MimicBaseMaxHealth = HighBaseHealth
	MimicBaseAgility   = HighBaseAgility
	MimicBaseStrength  = LowBaseStrength
)

type Character struct {
	Health    float64
	MaxHealth float64
	Agility   float64
	Strength  float64
	Treasure  int
	Buffs     []*ItemEffect
}

type MonsterCfg struct {
	ID         rune
	Name       string
	Character  Character
	Hostility  Hostility
	ColorStyle tcell.Style
}

// monsterTemplates
var MonsterConfig = map[MonsterType]MonsterCfg{
	Zombie: {ID: 'Z', Name: "Zombie",
		Character: Character{Health: ZombieBaseHealth, MaxHealth: ZombieBaseMaxHealth,
			Agility: ZombieBaseAgility, Strength: ZombieBaseStrength, Treasure: 3}, Hostility: HostilityAverage, ColorStyle: StyleZombie},

	Vampire: {ID: 'V', Name: "Vampire",
		Character: Character{Health: VampireBaseHealth, MaxHealth: VampireBaseMaxHealth,
			Agility: VampireBaseAgility, Strength: VampireBaseStrength, Treasure: 5}, Hostility: HostilityHigh, ColorStyle: StyleVampire},

	Ghost: {ID: 'G', Name: "Ghost", Character: Character{Health: GhostBaseHealth, MaxHealth: GhostBaseMaxHealth,
		Agility: GhostBaseAgility, Strength: GhostBaseStrength, Treasure: 5}, Hostility: HostilityLow, ColorStyle: StyleGhost},

	Ogre: {ID: 'O', Name: "Ogre", Character: Character{Health: OgreBaseHealth, MaxHealth: OgreBaseMaxHealth,
		Agility: OgreBaseAgility, Strength: OgreBaseStrength, Treasure: 5}, Hostility: HostilityAverage, ColorStyle: StyleOgre},

	Snake: {ID: 'S', Name: "Snake", Character: Character{Health: SnakeBaseHealth, MaxHealth: SnakeBaseMaxHealth,
		Agility: SnakeBaseAgility, Strength: SnakeBaseStrength, Treasure: 10}, Hostility: HostilityHigh, ColorStyle: StyleSnake},

	Mimic: {ID: 'M', Name: "Mimic", Character: Character{Health: MimicBaseHealth, MaxHealth: MimicBaseMaxHealth,
		Agility: MimicBaseAgility, Strength: MimicBaseStrength, Treasure: 10}, Hostility: HostilityLow, ColorStyle: StyleMimic},
}

// MonsterType represents kind of monster
type MonsterType int

const (
	MonsterTypeInvalid MonsterType = iota
	Zombie
	Vampire
	Ghost
	Ogre
	Snake
	Mimic
)

// Hostility represents how agressive the monster is and radius of chase
type Hostility int

const (
	HostilityLow Hostility = (iota + 1) * 2
	HostilityAverage
	HostilityHigh
)
