package utils

import (
	"strings"
	"time"
)

var dayNamesID = map[string]string{
	"Sunday":    "Minggu",
	"Monday":    "Senin",
	"Tuesday":   "Selasa",
	"Wednesday": "Rabu",
	"Thursday":  "Kamis",
	"Friday":    "Jumat",
	"Saturday":  "Sabtu",
}

var monthNamesID = map[string]string{
	"January":   "Januari",
	"February":  "Februari",
	"March":     "Maret",
	"April":     "April",
	"May":       "Mei",
	"June":      "Juni",
	"July":      "Juli",
	"August":    "Agustus",
	"September": "September",
	"October":   "Oktober",
	"November":  "November",
	"December":  "Desember",
}

var shortMonthNamesID = map[string]string{
	"Jan": "Jan",
	"Feb": "Feb",
	"Mar": "Mar",
	"Apr": "Apr",
	"May": "Mei",
	"Jun": "Jun",
	"Jul": "Jul",
	"Aug": "Agt",
	"Sep": "Sep",
	"Oct": "Okt",
	"Nov": "Nov",
	"Dec": "Des",
}

func FormatDateID(t time.Time, format string) string {
	result := t.Format(format)

	for en, id := range dayNamesID {
		result = strings.ReplaceAll(result, en, id)
	}

	for en, id := range monthNamesID {
		result = strings.ReplaceAll(result, en, id)
	}

	for en, id := range shortMonthNamesID {
		result = strings.ReplaceAll(result, en, id)
	}

	return result
}
