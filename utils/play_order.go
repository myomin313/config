package utils

import "strings"

// Display order for play types — does NOT affect IDs
var PlayOrder = map[string]int{
	"tiger":         1,
	"fish":          2,
	"shrimp":        3, // Gourd
	"crab":          4,
	"chicken":       5,
	"gold_coin":     6,
	"klouk":         7,
	"double_animal": 8,
	"triple_symbol": 9,  // Double Dices
	"triple_color":  10, // keep last (not in your list)
}

func PlayOrderIndex(key string) int {
	k := strings.TrimSpace(strings.ToLower(key))
	if v, ok := PlayOrder[k]; ok {
		return v
	}
	return 999
}
