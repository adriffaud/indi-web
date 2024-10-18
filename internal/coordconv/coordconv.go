package coordconv

import (
	"fmt"
	"math"
	"strconv"
)

func DecimalToSexagesimal(decimal string) (string, error) {
	decimalf, err := strconv.ParseFloat(decimal, 64)
	if err != nil {
		return "", err
	}

	hours, remainder := math.Modf(decimalf)
	minutesf := remainder * 60
	minutes, remainder := math.Modf(minutesf)
	seconds := int(remainder * 60)

	return fmt.Sprintf("%02.0f:%02.0f:%02d", hours, minutes, seconds), nil
}
