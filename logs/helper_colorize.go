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

func getTitle(level string) string {
	if level == "WARNING" {
		level = "WARN"
	}
	return level
}

func getLevelColor(level string) int {
	switch level {
	case "DEBUG":
		return colorBlue
	case "INFO":
		return colorGreen
	case "ERROR", "PANIC":
		return colorRed + 10
	case "WARNING":
		return colorYellow
	case "FATAL":
		return colorBrightRed
	default:
		return colorBlue
	}
}
