package enums

import (
	"fmt"
)

type Status uint

const (
	Active   Status = 1
	Disabled Status = 0
)

// String returns the string representation of DisplayDuration
func (d Status) String() string {
	switch d {
	case 0:
		return "disabled"
	case 1:
		return "active"
	default:
		return fmt.Sprintf("unknown(%d)", d)
	}
}

func ParseStatus(status string) Status {
	switch status {
	case "active":
		return Active
	default:
		return 0
	}
}
