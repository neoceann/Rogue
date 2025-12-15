package controller

import (
	"fmt"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"

	"github.com/gdamore/tcell/v2"
)

func CreateScreen() (tcell.Screen, error) {

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("new screen: %w", err)
		// panic(err)
	}
	if err := screen.Init(); err != nil {
		screen.Fini()
		return nil, fmt.Errorf("screen init: %w", err)
	}

	return screen, nil
}

func GetPlayerInput(screen tcell.Screen) (signal entities.GameSignal) {
	signal = entities.GsError
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				r := ev.Rune()
				switch {
				case r == 'w' || r == 'W':
					return entities.GsUp
				case r == 's' || r == 'S':
					return entities.GsDown
				case r == 'a' || r == 'A':
					return entities.GsLeft
				case r == 'd' || r == 'D':
					return entities.GsRight
				case r == config.UseWeaponInput:
					return entities.GsWeapon
				case r == config.UseFoodInput:
					return entities.GsFood
				case r == config.UsePotionInput:
					return entities.GsPotion
				case r == config.UseScrollInput:
					return entities.GsScroll
				case r == 'q' || r == 'Q':
					return entities.GsQuit
				case r == 'l' || r == 'L':
					return entities.GsNextLevel
				case r >= '1' && r <= '9':
					offset := r - '1' // 0..8
					return entities.GsChoice1 + entities.GameSignal(offset)
				case r == 'y' || r == 'Y':
					return entities.GsYes
				case r == 'n' || r == 'N':
					return entities.GsNo
				}
			case tcell.KeyEsc, tcell.KeyCtrlC:
				return entities.GsQuit
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}

func UpdateMapShot(m *entities.MapShot, mc *entities.MovementChanges) {
	m[mc.OldPos.X][mc.OldPos.Y] = mc.OldElement
	m[mc.NewPos.X][mc.NewPos.Y] = mc.NewElement
}

func RemoveItemFromGS(gs *entities.GameSession, item *entities.Item) {
	for i, it := range gs.Items {
		if item == it {
			gs.Items = append(gs.Items[:i], gs.Items[i+1:]...)
			break
		}
	}
}
