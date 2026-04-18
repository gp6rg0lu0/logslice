package filter

import (
	"fmt"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// TypeFilter matches log lines where a field's JSON value type matches the expected type.
// Valid types: string, number, bool, null, object, array
type TypeFilter struct {
	field    string
	wantType string
}

var validTypes = map[string]bool{
	"string": true,
	"number": true,
	"bool":   true,
	"null":   true,
	"object": true,
	"array":  true,
}

// NewTypeFilter returns a TypeFilter that matches lines where field has the given JSON type.
func NewTypeFilter(field, typeName string) (*TypeFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("typefilter: field must not be empty")
	}
	if !validTypes[typeName] {
		return nil, fmt.Errorf("typefilter: unknown type %q; valid types: string, number, bool, null, object, array", typeName)
	}
	return &TypeFilter{field: field, wantType: typeName}, nil
}

// Field returns the field name being checked.
func (f *TypeFilter) Field() string { return f.field }

// WantType returns the expected type name.
func (f *TypeFilter) WantType() string { return f.wantType }

// Match returns true if the field value's type matches wantType.
func (f *TypeFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	raw, ok := line.RawField(f.field)
	if !ok {
		return f.wantType == "null"
	}
	switch raw.(type) {
	case string:
		return f.wantType == "string"
	case float64:
		return f.wantType == "number"
	case bool:
		return f.wantType == "bool"
	case nil:
		return f.wantType == "null"
	case map[string]interface{}:
		return f.wantType == "object"
	case []interface{}:
		return f.wantType == "array"
	default:
		return false
	}
}
