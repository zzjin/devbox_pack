package utils

// ANSI color code constants
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorGray    = "\033[37m"
	ColorBold    = "\033[1m"
)

// Red returns red colored text
func Red(text string) string {
	return ColorRed + text + ColorReset
}

// Green returns green colored text
func Green(text string) string {
	return ColorGreen + text + ColorReset
}

// Yellow returns yellow colored text
func Yellow(text string) string {
	return ColorYellow + text + ColorReset
}

// Blue returns blue colored text
func Blue(text string) string {
	return ColorBlue + text + ColorReset
}

// Magenta returns magenta colored text
func Magenta(text string) string {
	return ColorMagenta + text + ColorReset
}

// Cyan returns cyan colored text
func Cyan(text string) string {
	return ColorCyan + text + ColorReset
}

// Gray returns gray colored text
func Gray(text string) string {
	return ColorGray + text + ColorReset
}

// Bold returns bold text
func Bold(text string) string {
	return ColorBold + text + ColorReset
}
