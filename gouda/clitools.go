package gouda

// ASCII Colors
const (
	Yellow = "\033[33m"
	Bold   = "\033[1m"
	Reset  = "\033[0m"
	Cheese = `
      ___ _____
     /\ (_)    \                           _
    /  \      (_,      __ _  ___  _   _  __| | __ _
   _)  _\   _    \    / _' |/ _ \| | | |/ _' |/ _' |
  /   (_)\_( )____\  | (_| | (_) | |_| | (_| | (_| |
  \_     /    _  _/   \__, |\___/ \__,_|\__,_|\__,_|
    ) /\/  _ ( )(     |___/
    \ \_) (_)   /
     \/________/


`
)

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
