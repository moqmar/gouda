package gouda

import (
	"strings"
)

// ASCII Colors
const (
	Yellow = "\033[33m"
	Bold   = "\033[1m"
	Reset  = "\033[0m"
)

var Cheese = strings.Replace(`\033[1;33m
    ▁▁   ▁▁▁▁
   ╱╲ ╰─╯    ╲     gouda 0.1
  ╮  ╲ ⭘   ╭─╮╲    the no-worries documentation tool
 ─╯   ╲    ╰─╯ ╲   
╱   ○  ╲▁╭─╮▁▁▁▁╲  \033[1;90mWrite and unite documentation for your projects from\033[1;33m
╲ ⭘    ╱ ╰─╯╭─╮ ╱  \033[1;90mvarious sources to create a uniquely consistent online\033[1;33m
 ╲   ╭─╮    ╰─╯╱   \033[1;90mdocumentation for every single part of it, and use lots\033[1;33m
  ╲  ╰─╯ ◯    ╱    \033[1;90mof plugins to make it as easy as possible for your users\033[1;33m
   ╲╱▁▁▁▁▁▁▁▁╱     \033[1;90mto get started with your project.\033[1;33m

\033[0m
`, `\033`, "\033", -1)

var progress = 0

func Progress() {
	switch progress {
	case 0:
		print("\033[2D▛ ")
	case 1:
		print("\033[2D▜ ")
	case 2:
		print("\033[2D▟ ")
	case 3:
		print("\033[2D▙ ")
	}
	progress = (progress + 1) % 4
}
