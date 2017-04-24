package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var strs = regexp.MustCompile("([a-z])([A-Z])")

// StringToHyphens returns the value of the string hypened.
func StringToHyphens(s string) string {
	matches := strs.FindStringSubmatch(s)
	return strings.Join(matches, "-")
}

var digits = regexp.MustCompile("\\d+")
var chars = regexp.MustCompile("\\D+")

// GetDuration returns the giving duration measure from the provided string value.
func GetDuration(ts string) (time.Duration, error) {
	value := chars.ReplaceAllString(ts, "")
	measure := digits.ReplaceAllString(ts, "")

	valInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	switch measure {
	case "ns":
		return time.Duration(valInt) * time.Nanosecond, nil
	case "ms":
		return time.Duration(valInt) * time.Millisecond, nil
	case "s":
		return time.Duration(valInt) * time.Second, nil
	case "m":
		return time.Duration(valInt) * time.Minute, nil
	case "h":
		return time.Duration(valInt) * time.Hour, nil
	default:
		return 0, errors.New("Invalid Measure for duration")
	}
}

// StringMatcher defines a function type to match a giving string against.
type StringMatcher func(string) bool

// MatchAny checks if the giving string matches any of the possiblities, it
// converts the giving type as necessary and if the possiblity is a function
// runs the function with the expected signature of `MatchString`.
// If runs until it finds a match then stops else returns false if non matches.
// Only checks for int,uint,float,rune,string and StringMatcher types.
func MatchAny(target string, possibilities ...interface{}) bool {
	var state bool

	for _, item := range possibilities {
		if state {
			break
		}

		switch item.(type) {
		case StringMatcher:
			if (item.(StringMatcher))(target) {
				state = true
				continue
			}
		case rune:
			if string(item.(rune)) == target {
				state = true
				continue
			}
		case string:
			if item.(string) == target {
				state = true
				continue
			}
		case int, int64:
			if target == strconv.FormatInt(item.(int64), 10) {
				state = true
				continue
			}
		case uint, uint32, uint64:
			if target == strconv.FormatUint(item.(uint64), 10) {
				state = true
				continue
			}
		case float32, float64:
			if target == strconv.FormatFloat(item.(float64), 'f', 1, 64) {
				state = true
				continue
			}
		case *regexp.Regexp:
			if (item.(*regexp.Regexp)).MatchString(target) {
				state = true
				continue
			}
		}
	}

	return state
}

// MatchAll checks if the giving string matches all the provided possibilities
// else returns false. It converts the possiblity into a string when possible
// and if its a StringMatcher type then runs the giving function.
// Only checks for int,uint,float,rune,string and StringMatcher types.
func MatchAll(target string, possibilities ...interface{}) bool {
	state := true

	for _, item := range possibilities {
		if !state {
			break
		}

		switch item.(type) {
		case StringMatcher:
			if !(item.(StringMatcher))(target) {
				state = false
				continue
			}
		case rune:
			if string(item.(rune)) != target {
				state = false
				continue
			}
		case string:
			if item.(string) != target {
				state = false
				continue
			}
		case int, int64:
			if target != strconv.FormatInt(item.(int64), 10) {
				state = false
				continue
			}
		case uint, uint32, uint64:
			if target != strconv.FormatUint(item.(uint64), 10) {
				state = false
				continue
			}
		case float32, float64:
			if target == strconv.FormatFloat(item.(float64), 'f', 1, 64) {
				state = false
				continue
			}
		case *regexp.Regexp:
			if (item.(*regexp.Regexp)).MatchString(target) {
				state = true
				continue
			}
		}
	}

	return state
}

// UUID returns a new uuid compatible string, it provides a makeshift UUID
// generator.
func UUID() (uuid string) {
	for i := 0; i < 32; i++ {
		rand.Seed(time.Now().UnixNano() + int64(i))
		random := rand.Intn(16)
		switch i {
		case 8, 12, 16, 20:
			uuid += "-"
		}
		switch i {
		case 12:
			uuid += fmt.Sprintf("%x", 4)
		case 16:
			uuid += fmt.Sprintf("%x", random&3|8)
		default:
			uuid += fmt.Sprintf("%x", random)
		}
	}
	return
}
