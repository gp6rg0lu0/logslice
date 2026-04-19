package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// VersionFilter matches log lines where a semver-like field satisfies a minimum
// and/or maximum version constraint (major.minor.patch, compared numerically).
type VersionFilter struct {
	field  string
	minVer [3]int
	maxVer [3]int
	hasMin bool
	hasMax bool
}

func parseSemver(s string) ([3]int, error) {
	parts := strings.SplitN(s, ".", 3)
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	var v [3]int
	for i := 0; i < 3; i++ {
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return v, fmt.Errorf("invalid version segment %q: %w", parts[i], err)
		}
		v[i] = n
	}
	return v, nil
}

func compareVersion(a, b [3]int) int {
	for i := 0; i < 3; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

// NewVersionFilter creates a filter that matches lines where field's semver
// value is within [minVer, maxVer]. Pass empty string to skip a bound.
func NewVersionFilter(field, minVer, maxVer string) (*VersionFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("version filter: field must not be empty")
	}
	if minVer == "" && maxVer == "" {
		return nil, fmt.Errorf("version filter: at least one of min or max must be set")
	}
	f := &VersionFilter{field: field}
	if minVer != "" {
		v, err := parseSemver(minVer)
		if err != nil {
			return nil, fmt.Errorf("version filter min: %w", err)
		}
		f.minVer = v
		f.hasMin = true
	}
	if maxVer != "" {
		v, err := parseSemver(maxVer)
		if err != nil {
			return nil, fmt.Errorf("version filter max: %w", err)
		}
		f.maxVer = v
		f.hasMax = true
	}
	if f.hasMin && f.hasMax && compareVersion(f.minVer, f.maxVer) > 0 {
		return nil, fmt.Errorf("version filter: min %q exceeds max %q", minVer, maxVer)
	}
	return f, nil
}

func (f *VersionFilter) Field() string { return f.field }

func (f *VersionFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	v, err := parseSemver(val)
	if err != nil {
		return false
	}
	if f.hasMin && compareVersion(v, f.minVer) < 0 {
		return false
	}
	if f.hasMax && compareVersion(v, f.maxVer) > 0 {
		return false
	}
	return true
}
