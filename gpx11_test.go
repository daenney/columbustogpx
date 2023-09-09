package main

import "testing"

func expectPanic(t *testing.T, f func(string) string, val string) {
	t.Helper()
	defer func() { _ = recover() }()
	f(val)
	t.Fatalf("expected a panic")
}

func TestLatLon(t *testing.T) {
	tests := []struct {
		in    string
		out   string
		panic bool
	}{
		{in: "59.0N", out: "59.0"},
		{in: "59.0E", out: "59.0"},
		{in: "59.0S", out: "-59.0"},
		{in: "59.0W", out: "-59.0"},
		{in: "59.0", out: "", panic: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			if tt.panic {
				expectPanic(t, latLon, tt.in)
			} else {
				if res := latLon(tt.in); res != tt.out {
					t.Fatalf("on input: %s, expected: %s, got: %s", tt.in, tt.out, res)
				}
			}
		})
	}
}
