package resample

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"fmt"
	"runtime"
	"strconv"
	"strings"
	urllib "net/url"
)

type Params struct {
	Fraction string
}

func (params Params) GetFraction() [2]int {
	if !strings.Contains(params.Fraction, "/") {
		x, _ := strconv.Atoi(params.Fraction)
		return [2]int{x, 1}
	}
	parts := strings.Split(params.Fraction, "/")
	numerator, _ := strconv.Atoi(parts[0])
	denominator, _ := strco