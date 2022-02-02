package tok

import (
	"math"
	"testing"
)

func TestReadInt(t *testing.T) {
	cases := []struct {
		inp     string
		base    int
		bitSize int
		exp     int64
		tail    string
	}{
		// ---------------------------------------------------------  8
		{"22", 10, 8, 22, ""},
		{"1c", 16, 8, 28, ""},
		{"0", 16, 8, 0, ""},
		{"-70", 8, 8, -56, ""},
		// boundaries
		{"-128", 10, 8, -128, ""},
		{"127", 10, 8, 127, ""},
		{"-2345", 10, 8, -23, "45"},
		{"128", 10, 8, 12, "8"},
		// with leading zeros
		{"-0046", 10, 8, -46, ""},
		// ignore other data at the end
		{"32-blocks", 10, 8, 32, "-blocks"},
		// --------------------------------------------------------- 16
		{"18", 10, 16, 18, ""},
		{"30df", 16, 16, 12511, ""},
		{"-4E3", 16, 16, -1251, ""},
		{"7561", 8, 16, 3953, ""},
		// boundaries
		{"-32768", 10, 16, math.MinInt16, ""},
		{"32767", 10, 16, math.MaxInt16, ""},
		{"32768", 10, 16, 3276, "8"},
		// ignore other data at the end
		{"345wxyz", 10, 16, 345, "wxyz"},
		// --------------------------------------------------------- 32
		// boundaries
		{"-2147483648", 10, 32, math.MinInt32, ""},
		{"2147483647", 10, 32, math.MaxInt32, ""},
		{"-2147483649", 10, 32, -214748364, "9"},
		// --------------------------------------------------------- 64
		{"42", 10, 64, 42, ""},
		{"aBcD", 16, 64, 43981, ""},
		{"-4a3F", 16, 64, -19007, ""},
		// boundaries
		{"-9223372036854775808", 10, 64, math.MinInt64, ""},
		{"9223372036854775807", 10, 64, math.MaxInt64, ""},
		{"9223372036854775808", 10, 64, 922337203685477580, "8"},
		// ignore other data at the end
		{"777 oth", 10, 64, 777, " oth"},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		val, err := sca.ReadInt(c.base, c.bitSize)
		if err != nil {
			t.Errorf("%d %q unexpected error: %v", i, c.inp, err)
		} else if val != c.exp {
			t.Errorf("%d unexpected result: %d != %d", i, val, c.exp)
		} else if sca.Tail() != c.tail {
			t.Errorf("%d %q scanner at wrong positiong, tail >%s<", i, c.inp, sca.Tail())
		}
	}

	failed := []struct {
		inp     string
		base    int
		bitSize int
	}{
		{"-abcd", 10, 64},
	}
	for i, f := range failed {
		sca := NewScanner(f.inp)
		_, err := sca.ReadInt(f.base, f.bitSize)
		if err == nil {
			t.Errorf("%d expected error for %s", i, f.inp)
		}
	}
}

func TestReadUint(t *testing.T) {
	cases := []struct {
		inp     string
		base    int
		bitSize int
		exp     uint64
		tail    string
	}{
		// ---------------------------------------------------------  8
		{"18", 10, 8, 18, ""},
		{"1c", 16, 8, 28, ""},
		{"F0", 16, 8, 240, ""},
		{"70", 8, 8, 56, ""},
		// boundaries
		{"0", 10, 8, 0, ""},
		{"255", 10, 8, 255, ""},
		{"2345", 10, 8, 234, "5"},
		{"256", 10, 8, 25, "6"},
		// with leading zeros
		{"00460", 10, 8, 46, "0"},
		// --------------------------------------------------------- 16
		// boundaries
		{"0", 10, 16, 0, ""},
		{"65535", 10, 16, 65535, ""},
		{"65536", 10, 16, 6553, "6"},
		// --------------------------------------------------------- 32
		// boundaries
		{"4294967295", 10, 32, math.MaxUint32, ""},
		{"4294967296", 10, 32, 429496729, "6"},
		// --------------------------------------------------------- 64
		// general
		{"30df", 16, 64, 12511, ""},
		{"4E3", 16, 64, 1251, ""},
		{"7561", 8, 64, 3953, ""},
		// boundaries
		{"18446744073709551615,0", 10, 64, math.MaxUint64, ",0"},
		// ignore other data at the end
		{"1170343number", 10, 64, 1170343, "number"},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		val, err := sca.ReadUint(c.base, c.bitSize)
		if err != nil {
			t.Errorf("%d %q unexpected error: %v", i, c.inp, err)
		} else if val != c.exp {
			t.Errorf("%d unexpected result: %d != %d", i, val, c.exp)
		} else if sca.Tail() != c.tail {
			t.Errorf("%d scanner at wrong positiong, tail >%s<", i, sca.Tail())
		}
	}
}
