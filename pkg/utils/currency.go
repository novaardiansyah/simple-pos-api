package utils

import (
	"fmt"
	"strings"
)

func FormatCurrency(number int64, currency string, showCurrency bool) string {
	isNegative := number < 0
	if isNegative {
		number = -number
	}

	formatted := formatWithDots(number)

	var result string
	if isNegative {
		result = "-" + currency + formatted
	} else {
		result = currency + formatted
	}

	if showCurrency {
		return result
	}

	return maskNumbers(result)
}

func FormatRupiah(number int64) string {
	return FormatCurrency(number, "Rp", true)
}

func formatWithDots(n int64) string {
	str := fmt.Sprintf("%d", n)
	length := len(str)

	if length <= 3 {
		return str
	}

	var result strings.Builder
	remainder := length % 3

	if remainder > 0 {
		result.WriteString(str[:remainder])
		if length > remainder {
			result.WriteString(".")
		}
	}

	for i := remainder; i < length; i += 3 {
		result.WriteString(str[i : i+3])
		if i+3 < length {
			result.WriteString(".")
		}
	}

	return result.String()
}

func maskNumbers(s string) string {
	result := []rune(s)
	for i, r := range result {
		if r >= '0' && r <= '9' {
			result[i] = '-'
		}
	}
	return string(result)
}

func FormatPercent(value int) string {
	return fmt.Sprintf("%d%%", value)
}
