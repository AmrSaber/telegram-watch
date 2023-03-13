package utils

import (
	"fmt"
)

type Unit struct {
	Name       string
	Multiplier float64 // How much of this unit equals the next one
}

func FormatNumber[T Number](value T, units []Unit) string {
	floatValue := float64(value)

	level := 0
	for floatValue >= units[level].Multiplier && level < len(units)-1 {
		floatValue /= units[level].Multiplier
		level++
	}

	return fmt.Sprintf("%.3f %s", floatValue, units[level].Name)
}

func FormatTime[T Number](nanoSeconds T) string {
	timeUnits := []Unit{
		{"ns", 1000},
		{"us", 1000},
		{"ms", 1000},
		{"sec", 60},
		{"min", 60},
		{"hr", -1},
	}

	return FormatNumber(nanoSeconds, timeUnits)
}
