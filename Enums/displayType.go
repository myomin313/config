package enums

import "fmt"

type DisplayDuration int

const (
	Schedule  DisplayDuration = 1
	Permanent DisplayDuration = 2
)

// String returns the string representation of DisplayDuration
func (d DisplayDuration) String() string {
	switch d {
	case 1:
		return "schedule"
	case 2:
		return "permanent"
	default:
		return fmt.Sprintf("unknown(%d)", d)
	}
}

func ParseDisplayDuration(DisplayDuration string) DisplayDuration {
	switch DisplayDuration {
	case "schedule":
		return 1
	case "permanent":
		return 2
	default:
		return 0
	}
}
