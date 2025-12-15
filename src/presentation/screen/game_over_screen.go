package screen

import (
	"fmt"
	data "rogue_game/datalayer"
	"github.com/gdamore/tcell/v2"
)

var (
	StatDivider = "+-----+--------------+-----------------+------+---------+" +
		"---------+---------+---------+---------+---------+---------+-------" +
		"--+---------+"
	StatHeader1 = "| N   | Game         | Name            |Status|Treasure | " +
		"Deepest |Monsters | Food    | Potions | Scrolls | Hits    | Hits    " +
		"|  Tiles  |"
	StatHeader2 = "|     | started      |                 |      |collected| " +
		"level   |defeated |consumed | drunk   | read    | dealt   |received " +
		"|travelled|"
	StatFormat = "| %-3s | %-12s | %-15s | %-4s | %-7s | %-7s | %-7s | %-7s " +
		"| %-7s | %-7s | %-7s | %-7s | %-7s |"
)

func DrawStatHeader(screen tcell.Screen, x, y int,
	style tcell.Style) (int, int) {
	header := []string{
		StatDivider,
		StatHeader1,
		StatHeader2,
		StatDivider,
	}
	x, y = DrawLines(screen, header, x, y, style)
	return x, y
}

func DrawStatData(screen tcell.Screen, stat *data.GameStats,
	x, y, n int, style tcell.Style) (int, int) {
	line := []string{
		fmt.Sprintf(StatFormat, stat.StatToString(n)...),
		StatDivider,
	}
	x, y = DrawLines(screen, line, x, y, style)

	return x, y
}

func DrawStaticStatisticView(screen tcell.Screen, summaryStat *data.GameStats, place,
	totalGames, x, y int) (int, int) {
	yGap := 2
	// place in leaderboard
	placeString := fmt.Sprintf("You are number %d in this race of %d games",
		place+1, totalGames)
	x, y = DrawLine(screen, placeString, x, y, YellowBold)
	// sum af all stat fields from recorded statistics
	x, y = DrawLine(screen, "Summary result since time immemorial",
		x, y+yGap, YellowBold)
	x, y = DrawStatHeader(screen, x, y, YellowBold)
	x, y = DrawStatData(screen, summaryStat, x, y, totalGames, Yellow)
	return x, y
}

func DrawAllStatsPages(screen tcell.Screen, stats []*data.GameStats,
	page, perPage, totalPages, x, y int) {
	yGap := 2
	pages := fmt.Sprintf("Page: %d out of %d", page, totalPages)
	prompt := []string{
		"* Press:",
		"      'A' for <-",
		"      'D' for ->",
		"      'Q' or 'Esc' to quit",
	}
	// Header for all results
	x, y = DrawLine(screen, "All results: "+pages,
		x, y+yGap, YellowBold)
	x, y = DrawLines(screen, prompt, x, y, Gray)
	x, y = DrawStatHeader(screen, x, y, YellowBold)

	for i := range len(stats) {
		n := (page-1)*perPage + i + 1
		x, y = DrawStatData(screen, stats[i], x, y, n, Yellow)
	}
}
