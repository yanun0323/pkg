package logs

import "fmt"

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

const (
	colorDarkRed = colorRed + 10
)

const (
	colorBrightBlack = iota + 90
	colorBrightRed
	colorBrightGreen
	colorBrightYellow
	colorBrightBlue
	colorBrightMagenta
	colorBrightCyan
	colorBrightWhite
)

func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
