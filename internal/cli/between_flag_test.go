package cli_test

import (
	"testing"

	"github.com/user/logslice/internal/cli"
)

func TestParseBetweenFlag_Empty(t *testing.T) {
	filters, err := cli.ParseBetweenFlag(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(filters) != 0 {
		t.Errorf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseBetweenFlag_Valid(t *testing.T) {
	filters, err := cli.ParseBetweenFlag([]string{"latency:10:200", "size:0:1024"})
	if err != nil {
		t.Fatal(err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
	if filters[0].Field() != "latency" {
		t.Errorf("expected latency, got %s", filters[0].Field())
	}
	if filters[0].Min() != 10 || filters[0].Max() != 200 {
		t.Errorf("unexpected range: %v %v", filters[0].Min(), filters[0].Max())
	}
	if filters[1].Field() != "size" {
		t.Errorf("expected size, got %s", filters[1].Field())
	}
}

func TestParseBetweenFlag_Errors(t *testing.T) {
	cases := []struct {
		input string
		desc  string
	}{
		{"latency:10", "missing max"},
		{"latency:abc:100", "non-numeric min"},
		{"latency:10:xyz", "non-numeric max"},
		{"latency:100:10", "min exceeds max"},
		{":10:100", "empty field"},
	}
	for _, tc := range cases {
		_, err := cli.ParseBetweenFlag([]string{tc.input})
		if err == nil {
			t.Errorf("case %q: expected error", tc.desc)
		}
	}
}
