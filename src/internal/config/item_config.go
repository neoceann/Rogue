package config

import "github.com/gdamore/tcell/v2"

const (
	StartLevelFoodPerRoom    = 2
	LastLevelFoodPerRoom     = 1
	StartLevelPotionsPerRoom = 1
	LastLevelPotionsPerRoom  = 0.5
	StartLevelScrollsPerRoom = 1
	LastLevelScrollsPerRoom  = 0.5
	WeaponInFirstLevel       = false // if false - first weapon on 2nd level
)

const (
	UseWeaponInput = 'h'
	UseFoodInput   = 'j'
	UsePotionInput = 'k'
	UseScrollInput = 'e'
)

const (
	BackpackItemStackSize = 5
	BackpackCapacity      = 9
)

// базовые виды предметов
type ItemType int

const (
	ItemTypeInvalid ItemType = iota
	ItemTypeFood
	ItemTypePotion
	ItemTypeScroll
	ItemTypeWeapon
	ItemTypeKey
)

// подтипы предметов
type ItemSubType int

const (
	ItemSubTypeInvalid ItemSubType = iota

	ItemSubTypeFoodApple
	ItemSubTypeFoodPear
	ItemSubTypeFoodBeefSteak

	ItemSubTypePotionStr
	ItemSubTypePotionAgi
	ItemSubTypePotionMaxHealth

	ItemSubTypeScrollStr
	ItemSubTypeScrollAgi
	ItemSubTypeScrollMaxHealth

	ItemSubTypeWeaponKnife
	ItemSubTypeWeaponDagger
	ItemSubTypeWeaponSword
	ItemSubTypeWeaponAxe
	ItemSubTypeWeaponKatana
	ItemSubTypeWeaponScythe

	ItemSubTypeKeyRed
	ItemSubTypeKeyGreen
	ItemSubTypeKeyBlue
)

var ItemsToCreate = map[ItemType][]ItemSubType{
	ItemTypeFood:   {ItemSubTypeFoodApple, ItemSubTypeFoodPear, ItemSubTypeFoodBeefSteak},
	ItemTypePotion: {ItemSubTypePotionStr, ItemSubTypePotionAgi, ItemSubTypePotionMaxHealth},
	ItemTypeScroll: {ItemSubTypeScrollStr, ItemSubTypeScrollAgi, ItemSubTypeScrollMaxHealth},
	ItemTypeWeapon: {ItemSubTypeWeaponKnife, ItemSubTypeWeaponDagger, ItemSubTypeWeaponSword,
		ItemSubTypeWeaponAxe, ItemSubTypeWeaponKatana, ItemSubTypeWeaponScythe},
	ItemTypeKey: {ItemSubTypeKeyRed, ItemSubTypeKeyGreen, ItemSubTypeKeyBlue},
}

// на какую характеристику влияет
type EffectType int

const (
	Strength EffectType = iota
	Agility
	Health
	MaxHealth
)

// свойства предмета
type ItemEffect struct {
	EffectTo    EffectType
	EffectValue float64
	Duration    int     //Длительность в количестве ходов. 0 для постоянных эффектов или мгновенного действия
}

type ItemCfg struct {
	ID           rune
	ItemType     ItemType
	SubType      ItemSubType
	Name         string
	Stackable    bool
	MaxStack     int
	StackCounter int
	Effect       ItemEffect
	ColorStyle   tcell.Style
}

var ItemConfig = map[ItemSubType]ItemCfg{

	// Food
	ItemSubTypeFoodApple: {ID: '☼', ItemType: ItemTypeFood, SubType: ItemSubTypeFoodApple,
		Name: "Magic apple", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{Health, 10, 0}, ColorStyle: StyleItemHealth},
	ItemSubTypeFoodPear: {ID: '◒', ItemType: ItemTypeFood, SubType: ItemSubTypeFoodPear,
		Name: "Juicy pear", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{Health, 20, 0}, ColorStyle: StyleItemHealth},
	ItemSubTypeFoodBeefSteak: {ID: '♨', ItemType: ItemTypeFood, SubType: ItemSubTypeFoodBeefSteak,
		Name: "Steak from unknown meat", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{Health, 30, 0}, ColorStyle: StyleItemHealth},

	// Potions
	ItemSubTypePotionStr: {ID: '⁂', ItemType: ItemTypePotion, SubType: ItemSubTypePotionStr,
		Name: "Potion of strength", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{Strength, 10, 30}, ColorStyle: StyleItemStr},
	ItemSubTypePotionAgi: {ID: '⁂', ItemType: ItemTypePotion, SubType: ItemSubTypePotionAgi,
		Name: "Potion of agility", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{Agility, 10, 30}, ColorStyle: StyleItemAgi},
	ItemSubTypePotionMaxHealth: {ID: '⁂', ItemType: ItemTypePotion, SubType: ItemSubTypePotionMaxHealth,
		Name: "Potion of health", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{MaxHealth, 10, 30}, ColorStyle: StyleItemHealth},

	// Scrolls
	ItemSubTypeScrollStr: {ID: '§', ItemType: ItemTypeScroll, SubType: ItemSubTypeScrollStr,
		Name: "Scroll of strength", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{Strength, 2, 0}, ColorStyle: StyleItemStr},
	ItemSubTypeScrollAgi: {ID: '§', ItemType: ItemTypeScroll, SubType: ItemSubTypeScrollAgi,
		Name: "Scroll of agility", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{Agility, 2, 0}, ColorStyle: StyleItemAgi},
	ItemSubTypeScrollMaxHealth: {ID: '§', ItemType: ItemTypeScroll, SubType: ItemSubTypeScrollMaxHealth,
		Name: "Scroll of health", Stackable: true, MaxStack: BackpackItemStackSize,
		Effect: ItemEffect{MaxHealth, 2, 0}, ColorStyle: StyleItemHealth},

	// Weapons
	ItemSubTypeWeaponKnife: {ID: '†', ItemType: ItemTypeWeapon, SubType: ItemSubTypeWeaponKnife,
		Name: "Kitchen knife", Stackable: false, MaxStack: 1,
		Effect: ItemEffect{Strength, 5, 0}, ColorStyle: StyleWeapon},
	ItemSubTypeWeaponDagger: {ID: '‡', ItemType: ItemTypeWeapon, SubType: ItemSubTypeWeaponDagger,
		Name: "Blunt dagger", Stackable: false, MaxStack: 1,
		Effect: ItemEffect{Strength, 10, 0}, ColorStyle: StyleWeapon},
	ItemSubTypeWeaponSword: {ID: '⚔', ItemType: ItemTypeWeapon, SubType: ItemSubTypeWeaponSword,
		Name: "Knight's swords", Stackable: false, MaxStack: 1,
		Effect: ItemEffect{Strength, 15, 0}, ColorStyle: StyleWeapon},
	ItemSubTypeWeaponAxe: {ID: '⚒', ItemType: ItemTypeWeapon, SubType: ItemSubTypeWeaponAxe,
		Name: "Axes of the furious Ogre", Stackable: false, MaxStack: 1,
		Effect: ItemEffect{Strength, 20, 0}, ColorStyle: StyleWeapon},
	ItemSubTypeWeaponKatana: {ID: '⚚', ItemType: ItemTypeWeapon, SubType: ItemSubTypeWeaponKatana,
		Name: "Demonic katana", Stackable: false, MaxStack: 1,
		Effect: ItemEffect{Strength, 25, 0}, ColorStyle: StyleWeapon},
	ItemSubTypeWeaponScythe: {ID: '☭', ItemType: ItemTypeWeapon, SubType: ItemSubTypeWeaponScythe,
		Name: "Scythe of DEVIL itself", Stackable: false, MaxStack: 1,
		Effect: ItemEffect{Strength, 30, 0}, ColorStyle: StyleWeapon},

	// Keys
	ItemSubTypeKeyRed: {ID: '⚷', ItemType: ItemTypeKey, SubType: ItemSubTypeKeyRed,
		Name: "Red key", Stackable: false, MaxStack: 1, ColorStyle: StyleKeyRed},
	ItemSubTypeKeyGreen: {ID: '⚷', ItemType: ItemTypeKey, SubType: ItemSubTypeKeyGreen,
		Name: "Green key", Stackable: false, MaxStack: 1, ColorStyle: StyleKeyGreen},
	ItemSubTypeKeyBlue: {ID: '⚷', ItemType: ItemTypeKey, SubType: ItemSubTypeKeyBlue,
		Name: "Blue key", Stackable: false, MaxStack: 1, ColorStyle: StyleKeyBlue},
}
