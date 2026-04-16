package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseExistsFlag parses a flag value of the form "field" or "!field".
// A leading '!' means the field must be absent; otherwise it must be present.
// Returns a configured ExistsFilter or an error.
func ParseExistsFlag(value string) (*filter.ExistsFilter, error) {
	if value == "" {
		return nil, fmt.Errorf("exists flag: value must not be empty")
	}
	mustExist := true
	field := value
	if strings.HasPrefix(value, "!") {
		mustExist = false
		field = strings.TrimPrefix(value, "!")
	}
	if field == "" {
		return nil, fmt.Errorf("exists flag: field name must not be empty after '!'")
	}
	return filter.NewExistsFilter(field, mustExist)
}
