package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/internal/token"
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

func CreateSession(deviceID int) (*model.Tokens, error) {
	tokenJWT := token.New(g.Config["token_secret"])

	accessTokenTTL, err := time.ParseDuration(g.Config["access_token_ttl"])
	if err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	refreshTokenTTL, err := time.ParseDuration(g.Config["refresh_token_ttl"])
	if err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	accessToken, err := tokenJWT.NewJWT(deviceID, accessTokenTTL)
	if err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	refreshToken, err := tokenJWT.NewRefreshToken()
	if err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	tokens := &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := store.I.Users().AddRefreshToken(deviceID, tokens.RefreshToken, refreshTokenTTL); err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	return tokens, nil
}

func GetNumber(v interface{}) (int, error) {
	switch v := v.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case string:
		if v, err := strconv.Atoi(v); err == nil {
			return v, nil
		}
	}

	return 0, errors.Wrap(errors.Errorf("value is not number (%T)", v), "GetNumber")
}
