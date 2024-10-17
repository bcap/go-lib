package numstat

import (
	"math"
	"strings"
)

// https://www.compart.com/en/unicode/block/U+2580
var blocks = []rune{
	'\u258F', // Left One Eighth Block:     ▏
	'\u258E', // Left One Quarter Block:    ▎
	'\u258D', // Left Three Eighths Block:  ▍
	'\u258C', // Left Half Block:           ▌
	'\u258B', // Left Five Eighths Block:   ▋
	'\u258A', // Left Three Quarters Block: ▊
	'\u2589', // Left Seven Eighths Block:  ▉
	'\u2588', // Full Block:                █
}

// Generates a unicode graphical bar
// The bar has always 102 characters: 2 for borders and a 100 for blocks
// Examples:
// - an empty bar:        ┃                                                                                                    ┃
// - a full bar:          ┃████████████████████████████████████████████████████████████████████████████████████████████████████┃
// - a 79.84% filled bar: ┃███████████████████████████████████████████████████████████████████████████████▊                    ┃
// - a 13.47% filled bar: ┃█████████████▍                                                                                      ┃
// Borders can be disabled with borders = false
func Bar(value float64, max float64, borders bool) string {
	if value > max {
		value = max
	}

	percentage := value / max * 100
	result := strings.Builder{}

	if borders {
		// start border
		// https://www.compart.com/en/unicode/block/U+2500
		result.WriteRune('\u2503') // Box Drawings Heavy Vertical: ┃
	}

	contentLength := 0

	// whole part
	for i := 0; i < int(percentage); i++ {
		result.WriteRune(blocks[7])
		contentLength++
	}

	// remainder
	remainder := (percentage - math.Trunc(percentage)) * 8.0
	if int(remainder) > 0 {
		result.WriteRune(blocks[int(remainder)-1])
		contentLength++
	}

	// fill with spaces
	for contentLength < 100 {
		result.WriteRune(' ')
		contentLength++
	}

	if borders {
		result.WriteRune('\u2503') // Box Drawings Heavy Vertical: ┃
	}

	return result.String()
}
