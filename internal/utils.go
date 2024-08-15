package internal

import "fmt"

// ANSI color codes for terminal output
const (
  Reset  = "\033[0m"
  Red    = "\033[31m"
  Green  = "\033[32m"
  Yellow = "\033[33m"
  Blue   = "\033[34m"
)

// Colorize wraps text in ANSI color codes
func Colorize(text string, color string) string {
  return fmt.Sprintf("%s%s%s", color, text, Reset)
}

