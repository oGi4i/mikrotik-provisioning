package models

import (
	"errors"
	"fmt"
	"strconv"
)

type RouterOSDuration int64

const (
	Second RouterOSDuration = 1
	Minute                  = 60 * Second
	Hour                    = 60 * Minute
	Day                     = 24 * Hour
)

var (
	errLeadingInt = errors.New("time: bad [0-9]*") // never printed
	unitMap       = map[string]int64{
		"s": int64(Second),
		"m": int64(Minute),
		"h": int64(Hour),
		"d": int64(Day),
	}
	units = []string{"d", "h", "m", "s"}
)

// String returns a string representing the duration in the form "1d18h3m5s".
func (d *RouterOSDuration) String() string {
	u := int64(*d)
	var result string

	result += fmtUnit(0, u)

	return result
}

func fmtUnit(index int, v int64) string {
	var result string
	if v/unitMap[units[index]] > 0 {
		result += fmt.Sprintf("%d%s", v/unitMap[units[index]], units[index])
		v = v % unitMap[units[index]]

		if index+1 < len(units) && v%unitMap[units[index]] != 0 {
			result += fmtUnit(index+1, v)
			return result
		} else {
			return result
		}
	} else {
		if index+1 < len(units) && v%unitMap[units[index]] != 0 {
			result += fmtUnit(index+1, v)
			return result
		} else {
			return ""
		}
	}
}

// ParseDuration parses a duration string.
// Valid time units are "s", "m", "h", "d".
func ParseDuration(s string) (RouterOSDuration, error) {
	// ([0-9]*[a-z]+)+
	orig := s
	var d int64

	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, errors.New("time: invalid duration " + orig)
	}
	for s != "" {
		var v int64
		var err error

		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, errors.New("time: invalid duration " + orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		if !pre {
			// no digits
			return 0, errors.New("time: invalid duration " + orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, errors.New("time: missing unit in duration " + orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, errors.New("time: unknown unit " + u + " in duration " + orig)
		}
		if v > (1<<63-1)/unit {
			// overflow
			return 0, errors.New("time: invalid duration " + orig)
		}
		v *= unit
		d += v
		if d < 0 {
			// overflow
			return 0, errors.New("time: invalid duration " + orig)
		}
	}

	return RouterOSDuration(d), nil
}

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x int64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > (1<<63-1)/10 {
			// overflow
			return 0, "", errLeadingInt
		}
		x = x*10 + int64(c) - '0'
		if x < 0 {
			// overflow
			return 0, "", errLeadingInt
		}
	}
	return x, s[i:], nil
}

func (d *RouterOSDuration) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		return nil
	}

	duration, err := ParseDuration(unquoted)
	if err != nil {
		return nil
	}

	*d = duration
	return nil
}

func (d *RouterOSDuration) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}
