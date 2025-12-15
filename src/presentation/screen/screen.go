package screen

import (
	"fmt"
	"rogue_game/internal/config"
	data "rogue_game/datalayer"
	"rogue_game/internal/domain/character"
	"rogue_game/internal/domain/entities"
	"sort"

	"github.com/gdamore/tcell/v2"
)

const startPrintingX = config.GameFieldWidth + 1

func PrintMapUI(screen tcell.Screen, m entities.MapShot) {
	var ch rune
	var style tcell.Style
	for x := range m {
		for y := range m[0] {
			switch v := m[x][y].(type) {
			case entities.Field:
				switch v {
				case entities.RoomFloor:
					ch = ' '
					style = config.DefaultStyle
				case entities.CorridorFloor:
					ch = ' '
					style = config.DefaultStyle
				case entities.Wall:
					ch = '█'
					style = config.DefaultStyle
				case entities.Outside:
					ch = '░'
					style = config.DefaultStyle
				case entities.Finish:
					ch = 'F'
					style = config.DefaultStyle
				}
			case *entities.Player:
				ch = '⚇'
				style = config.StylePlayer
			case *entities.Item:
				ch = v.ID
				style = config.ItemConfig[v.SubType].ColorStyle
			case *entities.Monster:
				if v.Type == config.Mimic && bool(v.Trick) {
					ch = '§'
					style = config.StyleItemHealth
				} else if v.Type == config.Ghost && bool(v.Trick) {
					ch = ' '
					style = config.DefaultStyle
				} else {
					ch = v.ID
					style = config.MonsterConfig[v.Type].ColorStyle
				}

			case *entities.Door:
				switch v.Color{
				case entities.RedDoor:
					style = config.StyleKeyRed
				case entities.BlueDoor:
					style = config.StyleKeyBlue
				case entities.GreenDoor:
					style = config.StyleKeyGreen
				}
				ch = '#'

			default:
				ch = '?'
				style = config.DefaultStyle
			}

			screen.SetContent(x, y, ch, nil, style)
		}
	}
}

func PrintInfoUI(screen tcell.Screen, gs entities.GameSession, gstat *data.GameStats) {
	DrawText(screen, startPrintingX, 0, fmt.Sprintf("%s IS MY HERO! LEVEL: %d", gstat.PlayerName, gs.Level), tcell.StyleDefault.Bold(true))

	if gs.ActiveWeapon == nil {
		DrawText(screen, startPrintingX, 2, "ACTIVE WEAPON: NONE", tcell.StyleDefault.Bold(true))
	} else {
		DrawText(screen, startPrintingX, 2, "ACTIVE WEAPON:", tcell.StyleDefault.Bold(true))
		DrawText(screen, startPrintingX, 3, fmt.Sprintf("%s (%.1f %s)",
			gs.ActiveWeapon.Name, gs.ActiveWeapon.Effect.EffectValue,
			entities.GetEffectName(gs.ActiveWeapon.Effect.EffectTo)),
			tcell.StyleDefault.Bold(true))
	}

	DrawText(screen, startPrintingX, 5, fmt.Sprintf("HEALTH: %.1f / %.1f", gs.Health, gs.MaxHealth), tcell.StyleDefault.Bold(true))
	DrawText(screen, startPrintingX, 6, fmt.Sprintf("STRENGTH: %.1f", gs.Strength), tcell.StyleDefault.Bold(true))
	DrawText(screen, startPrintingX, 7, fmt.Sprintf("AGILITY: %.1f", gs.Agility), tcell.StyleDefault.Bold(true))
	DrawText(screen, startPrintingX, 8, fmt.Sprintf("TREASURES: %d", gs.Treasure), tcell.StyleDefault.Bold(true))

	DrawText(screen, startPrintingX, 10, "Buffs:", tcell.StyleDefault)
	for i, b := range gs.Player.Buffs {
		DrawText(screen, startPrintingX, i+11, fmt.Sprintf("%.1f %s for %d moves", b.EffectValue, entities.GetEffectName(b.EffectTo), b.Duration), tcell.StyleDefault)
	}

	DrawText(screen, startPrintingX+40, 0, fmt.Sprintf("Backpack: (capacity: %d)", config.BackpackCapacity), tcell.StyleDefault)
	sortedBackp := SortedBackpack(&gs.Player.Backpack)
	i := 0
	for _, elem := range sortedBackp {
		if gs.Backpack[elem].ItemType == config.ItemTypeWeapon {
			DrawText(screen, startPrintingX+40, i+1, fmt.Sprintf("%d. '%c': ", i + 1, config.UseWeaponInput)+gs.Backpack[elem].Name+gs.Backpack[elem].Description, tcell.StyleDefault.Bold(true))
		} else {
			var c rune
			switch gs.Backpack[elem].ItemType {
			case config.ItemTypeFood:
				c = config.UseFoodInput
			case config.ItemTypePotion:
				c = config.UsePotionInput
			case config.ItemTypeScroll:
				c = config.UseScrollInput
			}

			if gs.Backpack[elem].ItemType != config.ItemTypeKey {
				DrawText(screen, startPrintingX+40, i+1, fmt.Sprintf("%d. '%c': ", i + 1, c)+gs.Backpack[elem].Name+fmt.Sprintf(" (%d)",
					gs.Backpack[elem].StackCounter)+gs.Backpack[elem].Description, tcell.StyleDefault.Bold(true))
			} else {
				DrawText(screen, startPrintingX+40, i+1, gs.Backpack[elem].Name+gs.Backpack[elem].Description, tcell.StyleDefault.Bold(true))
			}
		}

		i++
	}
}

func DrawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

func PrintAvailableItems(screen tcell.Screen, items []*entities.Item, signal entities.GameSignal) {
	DrawText(screen, startPrintingX+40, 10, "Choose item:", tcell.StyleDefault)

	if signal == entities.GsWeapon {
		DrawText(screen, startPrintingX+40, 11, "0. Unequip active weapon (to backpack)", tcell.StyleDefault)
	} else {
		DrawText(screen, startPrintingX+40, 11, "0. Cancel", tcell.StyleDefault)
	}

	for i, item := range items {
		if item.ItemType == config.ItemTypeWeapon {
			DrawText(screen, startPrintingX+40, 12+i,
				fmt.Sprintf("%d. %s%s", i+1, item.Name, item.Description),
				tcell.StyleDefault)
		} else {
			DrawText(screen, startPrintingX+40, 12+i,
				fmt.Sprintf("%d. %s(%d)%s", i+1, item.Name, item.StackCounter, item.Description),
				tcell.StyleDefault)
		}

	}
	screen.Show()
}

func PrintMonstersAroundHero(screen tcell.Screen, gs *entities.GameSession, m *entities.MapShot) {
	for i, monster := range character.MonstersAroundHero(gs, m) {
		DrawText(screen, startPrintingX+i*25, 25, monster.Name, tcell.StyleDefault)
		DrawText(screen, startPrintingX+i*25, 26, fmt.Sprintf("HEALTH: %.1f / %.1f", monster.Health, monster.MaxHealth), tcell.StyleDefault)
		DrawText(screen, startPrintingX+i*25, 27, fmt.Sprintf("STRENGTH: %.1f", monster.Strength), tcell.StyleDefault)
		DrawText(screen, startPrintingX+i*25, 28, fmt.Sprintf("AGILITY: %.1f", monster.Agility), tcell.StyleDefault)
		DrawText(screen, startPrintingX+i*25, 29, fmt.Sprintf("TREASURES: %d", monster.Treasure), tcell.StyleDefault)
	}
}

func PrintFightInfo(screen tcell.Screen, gs *entities.GameSession, m *entities.MapShot) {
	pstatus := gs.Player.FightStatus

	if !pstatus.InFight {
		return
	}

	var message string

	if pstatus.MissingHit {
		message = "You missed"
	} else if pstatus.Asleep {
		message = "You fell asleep for 1 action"
	} else {
		message = fmt.Sprintf("You hits %s for %.1f hp", pstatus.TargetMonsterName, pstatus.DamageDone)
	}

	DrawText(screen, startPrintingX, 17, message, tcell.StyleDefault)

	for i, m := range character.MonstersAroundHero(gs, m) {
		if m.FightStatus.InFight {
			if m.FightStatus.MissingHit {
				message = fmt.Sprintf("%s missed", m.Name)
			} else {
				message = fmt.Sprintf("%s hits you for %.1f hp", m.Name, m.FightStatus.DamageDone)
			}
			DrawText(screen, startPrintingX, 18+i, message, tcell.StyleDefault)
		}
	}
}

func SortedBackpack(b *entities.Backpack) []config.ItemSubType {
	keys := make([]config.ItemSubType, 0, len(*b))
	for k := range *b {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}