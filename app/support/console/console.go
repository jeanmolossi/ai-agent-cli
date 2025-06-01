package console

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/pterm/pterm"
)

func TwoColumnDetail(first, second string, filler ...rune) string {
	margin := func(s string, left, right int) string {
		var builder strings.Builder
		if left > 0 {
			builder.WriteString(strings.Repeat(" ", left))
		}
		builder.WriteString(s)
		if right > 0 {
			builder.WriteString(strings.Repeat(" ", right))
		}
		return builder.String()
	}
	width := func(s string) int {
		return runewidth.StringWidth(pterm.RemoveColorFromString(s))
	}
	first = margin(first, 2, 1)
	if w := width(second); w > 0 {
		second = margin(second, 1, 2)
	} else {
		second = margin(second, 0, 2)
	}
	fillingText := ""
	if w := pterm.GetTerminalWidth() - width(first) - width(second); w > 0 {
		fillingText = strings.Repeat(string(append(filler, '.')[0]), w)
	}

	return first + fillingText + second
}
