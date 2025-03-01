package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

var Rnd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

var dayOfWeekEnToRuMap = map[string]string{
	"Monday":    "пн",
	"Tuesday":   "вт",
	"Wednesday": "ср",
	"Thursday":  "чт",
	"Friday":    "пт",
	"Saturday":  "сб",
	"Sunday":    "вс",
}

func DayOfWeekEnToRu(v string) string {
	return dayOfWeekEnToRuMap[v]
}

var monthEnToRuMap = map[string]string{
	"January":   "Январь",
	"February":  "Февраль",
	"March":     "Март",
	"April":     "Апрель",
	"May":       "Май",
	"June":      "Июнь",
	"July":      "Июль",
	"August":    "Август",
	"September": "Сентябрь",
	"October":   "Октябрь",
	"November":  "Ноябрь",
	"December":  "Декабрь",
}

func MonthEnToRu(month string, short bool) string {
	if short {
		if m, ok := monthEnToRuMap[month]; ok {
			return m[:3]
		}
		return ""
	} else {
		return monthEnToRuMap[month]
	}
}

func MD5(v string) string {
	h := md5.New()
	h.Write([]byte(v))
	return hex.EncodeToString(h.Sum(nil))
}
