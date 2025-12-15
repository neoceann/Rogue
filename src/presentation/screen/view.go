package screen

import (
	"github.com/gdamore/tcell/v2"
)

var (
	YellowBold = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	Yellow     = tcell.StyleDefault.Foreground(tcell.ColorYellow)
	RedBold    = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	GreenBold  = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	BlueBold   = tcell.StyleDefault.Foreground(tcell.ColorBlue).Bold(true)
	Gray       = tcell.StyleDefault.Foreground(tcell.ColorGray)
)

// draws one line of text (string)
func DrawLine(screen tcell.Screen, line string, x, y int,
	style tcell.Style) (int, int) {
	for i, ch := range line {
		screen.SetContent(x+i, y, ch, nil, style)
		// screen.Show() later check
	}
	screen.Show() // also take away?
	return x, y + 1
}

// draw all lines of text ([]string)
func DrawLines(screen tcell.Screen, lines []string, x, y int,
	style tcell.Style) (int, int) {
	for _, line := range lines {
		x, y = DrawLine(screen, line, x, y, style)
	}
	return x, y
}

func DrawWinBanner(screen tcell.Screen, x, y int) (int, int) {
	banner := []string{
		"",
		" _  _  _____  __  __    _    _  ____  _  _",
		"( \\/ )(  _  )(  )(  )  ( \\/\\/ )(_  _)( \\( )",
		" \\  /  )(_)(  )(__)(    )    (  _)(_  )  (",
		" (__) (_____)(______)  (__/\\__)(____)(_)\\_)",
		""}
	x, y = DrawLines(screen, banner, x, y, GreenBold)
	return x, y
}

func DrawLoseBanner(screen tcell.Screen, x, y int) (int, int) {
	banner := []string{
		"",
		" _  _  _____  __  __    __    _____  ___  ____",
		"( \\/ )(  _  )(  )(  )  (  )  (  _  )/ __)( ___)",
		" \\  /  )(_)(  )(__)(    )(__  )(_)( \\__ \\ )__)",
		" (__) (_____)(______)  (____)(_____)(___/(____)",
		"",
	}
	x, y = DrawLines(screen, banner, x, y, RedBold)
	return x, y
}

func DrawSavedBanner(screen tcell.Screen, x, y int) (int, int) {
	banner := []string{
		"",
		"  ___    __    __  __  ____    ___    __  _  _  ____  ____",
		" / __)  /__\\  (  \\/  )( ___)  / __)  /__\\( \\/ )( ___)(  _ \\",
		"( (_-. /(__)\\  )    (  )__)   \\__ \\ /(__)\\\\  /  )__)  )(_) )",
		" \\___/(__)(__)(_/\\/\\_)(____)  (___/(__)(__)\\/  (____)(____/",
		"",
	}
	x, y = DrawLines(screen, banner, x, y, BlueBold)
	return x, y
}

func DrawStartPicUI(screen tcell.Screen, x, y int) (int, int) {
	welcomeDraw := []string{
		"                                         .             -:.     - -",
		"                                ::=      ..    -.:.   . ::  :      .    : :.",
		"                               ... =    :  @@@@@@@@@@@@@@--- : - -+.      +",
		"                .    . .         . .-=* :@@@@@@@@@@@@@@@@@@.=   . =      =   :    ..",
		"               :    = : .     =: +-: .:=@@@@@@@@@@@@@@@@@@@@   . -:.*   :.       -",
		"               :     .      . :  =+  :.*@@@@@-  :@@@@@@@@@@@.::.:: .-.-=.# ..   :-.    :",
		"                          - - - + :  -:+@@@@@*  #@@@@@@@@@@@#",
		"                                       *@@@@%   :@@@@@@%@@@@.:.- --=  .   ::=-  .",
		"                      =-+.:. ..  .:=.:  .=@#      #@@@ @@@@@-  .  . . .-    ..-= =",
		"                   - -     - *:. -#-.-==@%      .  @@ @@@@@@#    =.::        .    :",
		"                   :    .   -.         @@@@@#:  .: @ %@@@@@@- :-:. . =                      :",
		"                            -      . *.:@@@*.    %  %@@@@@@@:.=.  +             -",
		"                   -  :        :..   .-=@@@@     **=@@@@@@@@::.- :      :-+     :      . ",
		"              :    =::  .     ==:* .   .@@@@+ *= *@@@@@@@@@@+==- . ::             .  .     .",
		":           :        -   .+       *...-=@@@@% *%  @@@@@@@@@@+    ----=   :      : : -. -",
		"                         -  --. #:=- - #@@@@  %@: @@@@@@@@@@:- -::  ::: -     .              .",
		"                    .  .  : =:: .:.#: . @@= +#%@  #@%%%@##*%.--. - -- :     : .     ::     :.",
		"                :.-  ..- .  .::.=. :  ...+  .-%=- =  - *. :  =   -  .:  . :               :.",
		"",
		"",
		"                 +-----------------------------------------------------------------+",
		"                 |            Welcome to Roguelike, stranger...                    |",
		"",
		"                 |   Death is permanent. Glory is fleeting. The dungeon awaits.    |",
		"                 +-----------------------------------------------------------------+",
		"",
		"",
		"                      +-------------------------------------------------------+",
		"                      |           Do you want to check saved games?           |",
		"                      | 'N' start new game | 'Y' yes, let's continue old game |",
		"                      +-------------------------------------------------------+",
		"",
	}
	screen.Clear()
	x, y = DrawLines(screen, welcomeDraw, x, y, YellowBold)
	screen.Show()
	return x, y
}
