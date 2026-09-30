package enums

import (
	"fmt"
	"strconv"
	"strings"
)

type DisplayPlatform uint

const (
	H5     DisplayPlatform = 1
	Web    DisplayPlatform = 2
	Mobile DisplayPlatform = 3
)

var currencyMap = map[uint64]string{
	1: "USD",
	2: "KHR",
	3: "CNY",
	4: "THB",
	5: "VND",
	6: "SGD",
	7: "MYR",
}

// String returns the string representation of DisplayDuration
func (d DisplayPlatform) String() string {
	switch d {
	case 1:
		return "h5"
	case 2:
		return "web"
	case 3:
		return "mobile"
	default:
		return fmt.Sprintf("unknown(%d)", d)
	}
}

func ParseDisplayPlatform(platform string) DisplayPlatform {
	switch platform {
	case "h5":
		return H5
	case "web":
		return Web
	case "mobile":
		return Mobile
	default:
		return 0
	}
}

func ParseCurrency(currencyID uint64) string {
	if code, ok := currencyMap[currencyID]; ok {
		return code
	}
	return "USD"
}

func ParseCurrencyID(currency string) uint64 {
	cur := strings.ToUpper(strings.TrimSpace(currency))

	for id, code := range currencyMap {
		if code == cur {
			return id
		}
	}

	id, err := strconv.ParseUint(cur, 10, 64)
	if err != nil {
		return 1
	}

	return id
}
