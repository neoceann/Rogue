package screen

import (
	"fmt"
	data "rogue_game/datalayer"

	"github.com/gdamore/tcell/v2"
)

var (
	SavedDivider = "+--------+---------------------+-------------------------" +
		"---+----------+----------+"
	SavedHeader = "| Choose | Saved               | Player name              " +
		"  | Treasure |  Level   |"
	SavedFormat = "| %-6s | %-19s | %-26s | %-8s | %-8s |"
)

func NameInputUI(screen tcell.Screen, x, y int) string {
	prompt := "» Enter your name, stranger: "
	var name []rune
	DrawLines(screen, []string{prompt}, x, y, YellowBold)
	x = len(prompt) + 2
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter, tcell.KeyCtrlJ:
				screen.Sync()
				return string(name)

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(name) > 0 {
					name = name[:len(name)-1]
					x--
					// delete last symbol
					screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
					screen.Show()
				}
			case tcell.KeyRune:
				r := ev.Rune()
				// only letters, figures, '-', '_', ' '
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
					(r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' ' {
					name = append(name, r)
					screen.SetContent(x, y, r, nil,
						tcell.StyleDefault.Foreground(tcell.ColorYellow))
					x++
					screen.Show()
				}
			case tcell.KeyEscape:
				continue
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}

func DrawDownloads(screen tcell.Screen, games []*data.GameSave,
	page int, totalPages int) {
	pages := fmt.Sprintf("Page: %d out of %d", page, totalPages)

	header := []string{
		pages,
		"",
		SavedDivider,
		SavedHeader,
		SavedDivider,
	}

	prompt := []string{
		"* Press:",
		"       1-9 to download game",
		"      'A' for <-",
		"      'D' for ->",
		"      'Esc' to cancel and start new game",
	}
	screen.Clear()
	x := 2
	y := 2
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	x, y = DrawLines(screen, header, x, y, style)
	style = tcell.StyleDefault.Foreground(tcell.ColorYellow)
	for i := range len(games) {
		choose := fmt.Sprintf("%3d", i+1)
		savedTime := games[i].SaveTime.Format("15:04 Jan 2")
		name := games[i].GameStats.PlayerName // make function to cut
		treasure := fmt.Sprintf("%3d", games[i].GameStats.TreasureCollected)
		level := fmt.Sprintf("%3d", games[i].GameStats.DeepestLevel)
		g := []string{
			fmt.Sprintf(SavedFormat, choose, savedTime, name, treasure, level),
			SavedDivider,
		}
		x, y = DrawLines(screen, g, x, y, style)
	}
	style = tcell.StyleDefault.Foreground(tcell.ColorGray)
	DrawLines(screen, prompt, x, y, style)
}
