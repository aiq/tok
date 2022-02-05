package tok

import (
	"fmt"
	"math"
	"unicode/utf8"
)

// ReadRune reads one rune value from the scanner.
func (s *Scanner) ReadRune() (rune, error) {
	r, i := utf8.DecodeRuneInString(s.Tail())
	if r == utf8.RuneError {
		return utf8.RuneError, s.ErrorFor("rune")
	}
	return r, s.BoolErrorFor(s.Move(i), "rune")
}

// RevReadRune reads one rune value from the scanner the reverse way.
func (s *Scanner) RevReadRune() (rune, error) {
	r, i := utf8.DecodeLastRuneInString(s.Head())
	if r == utf8.RuneError {
		return r, s.ErrorFor("rune")
	}
	return r, s.BoolErrorFor(s.Move(-i), "rune")
}

// ReadBool reads bool value from the scanner.
// Valid format values are
// - "l" for true and false
// - "U" for TRUE and FALSE
// - "Cc" for True and False
// - "*" for all cases
// An empty format string will be interpreted as "*".
func (s *Scanner) ReadBool(format string) (bool, error) {
	if format == "" || format == "*" {
		if s.IfAny("true", "True", "TRUE") {
			return true, nil
		} else if s.IfAny("false", "False", "FALSE") {
			return false, nil
		}
		return false, s.ErrorFor("bool")
	}

	trueStr, falseStr := "", ""
	switch format {
	case "l":
		trueStr, falseStr = "true", "false"
	case "U":
		trueStr, falseStr = "TRUE", "FALSE"
	case "Cc":
		trueStr, falseStr = "True", "False"
	default:
		return false, fmt.Errorf("invalid format")
	}

	if s.If(trueStr) {
		return true, nil
	} else if s.If(falseStr) {
		return false, nil
	}
	return false, s.ErrorFor("bool")
}

// ReadInt reads a integer value from the scanner.
// Valid base values are 8. 10 and 16.
// Valid bitSize values are 8, 16, 32 and 64.
func (s *Scanner) ReadInt(base int, bitSize int) (int64, error) {
	var charFunc func(rune) int32
	switch base {
	case 8:
		charFunc = octValue
	case 10:
		charFunc = decValue
	case 16:
		charFunc = hexValue
	default:
		return 0, fmt.Errorf("invalid base value %d", base)
	}

	var min int64
	var max int64
	switch bitSize {
	case 8:
		min, max = math.MinInt8, math.MaxInt8
	case 16:
		min, max = math.MinInt16, math.MaxInt16
	case 32:
		min, max = math.MinInt32, math.MaxInt32
	case 64:
		min, max = math.MinInt64, math.MaxInt64
	default:
		return 0, fmt.Errorf("invalid bitSize value %d", bitSize)
	}

	var i64 int64
	var tmp int64
	n := -1
	marker := s.Mark()
	neg := false
	if s.IfRune('-') {
		neg = true
	} else if s.IfRune('+') {
		neg = false
	}
	for i, r := range s.Tail() {
		v := charFunc(r)
		if v == -1 {
			break
		}

		if neg {
			tmp = (i64)*int64(base) - int64(v)
			if i64 < tmp || tmp < min {
				break
			}
		} else {
			tmp = (i64 * int64(base)) + int64(v)
			if i64 > tmp || tmp > max {
				break
			}
		}

		i64 = tmp
		n = i
	}
	if n == -1 {
		s.ToMarker(marker)
		return 0, s.ErrorFor("integer")
	}

	s.Move(n + 1)
	return i64, nil
}

func octValue(r rune) int32 {
	if inRange('0', r, '7') {
		return r - '0'
	}
	return -1
}

func decValue(r rune) int32 {
	if inRange('0', r, '9') {
		return r - '0'
	}
	return -1
}

func hexValue(r rune) int32 {
	if inRange('0', r, '9') {
		return r - '0'
	} else if inRange('a', r, 'f') {
		return (r - 'a') + 10
	} else if inRange('A', r, 'F') {
		return (r - 'A') + 10
	}
	return -1
}

// ReadUint reads a unsigned integer value from the scanner.
// Valid base values are 8, 10 and 16.
// Valid bitSize values are 8, 16, 32 and 64.
func (s *Scanner) ReadUint(base int, bitSize int) (uint64, error) {
	var charFunc func(rune) int32
	switch base {
	case 8:
		charFunc = octValue
	case 10:
		charFunc = decValue
	case 16:
		charFunc = hexValue
	default:
		return 0, fmt.Errorf("invalid base value %d", base)
	}

	var max uint64
	switch bitSize {
	case 8:
		max = math.MaxUint8
	case 16:
		max = math.MaxUint16
	case 32:
		max = math.MaxUint32
	case 64:
		max = math.MaxUint64
	default:
		return 0, fmt.Errorf("invalid bitSize value %d", bitSize)
	}

	var u64 uint64
	var tmp uint64
	n := -1
	for i, r := range s.Tail() {
		v := charFunc(r)
		if v == -1 {
			break
		}

		tmp = (u64 * uint64(base)) + uint64(v)
		if u64 > tmp || tmp > max {
			break
		}

		u64 = tmp
		n = i
	}
	if n == -1 {
		return 0, s.ErrorFor("unsigned integer")
	}

	s.Move(n + 1)
	return u64, nil
}
