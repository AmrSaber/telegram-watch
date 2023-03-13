package utils

import "testing"

func TestFormatTime(t *testing.T) {
	cases := map[int]string{
		100:            "100.000 ns",
		10_000:         "10.000 us",
		10_000_000:     "10.000 ms",
		10_000_000_000: "10.000 sec",
	}

	for input, expected := range cases {
		found := FormatTime(input)
		if found != expected {
			t.Errorf("expected %s, found %s", expected, found)
		}
	}
}
