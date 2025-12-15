package controller

import (
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
	"rogue_game/presentation/screen"

	"github.com/gdamore/tcell/v2"
)

func RequestItemForUse(gs *entities.GameSession, scr tcell.Screen, signal entities.GameSignal) (*entities.Item, bool) {

	var items []*entities.Item
	var item *entities.Item

	unequipFlag := false

	switch signal{
	case entities.GsWeapon:
		items = gs.Player.Backpack.GetItemsByType(config.ItemTypeWeapon)
		if items == nil && gs.Player.ActiveWeapon != nil {
			unequipFlag = true
		}
	case entities.GsFood:
		items = gs.Player.Backpack.GetItemsByType(config.ItemTypeFood)
	case entities.GsPotion:
		items = gs.Player.Backpack.GetItemsByType(config.ItemTypePotion)
	case entities.GsScroll:
		items = gs.Player.Backpack.GetItemsByType(config.ItemTypeScroll)
	default:
		items = nil
	}

	if items == nil && !unequipFlag {
		return item, unequipFlag
	}

	screen.PrintAvailableItems(scr, items, signal)
	
	correctChoice := false
	for !correctChoice {
		ev := scr.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case '0':
					if signal == entities.GsWeapon {
						unequipFlag = true
					}
					correctChoice = true
				case '1', '2', '3', '4', '5', '6', '7', '8', '9':
					 index := int(ev.Rune() - '1')

					 if index > -1 && index < len(items) {
					 	item = items[index]
						correctChoice = true
					}
				}
			}
		}
	}

	return item, unequipFlag
}