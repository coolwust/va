package validator

import (
	"fmt"
)

type MalformedError struct {
	Target *Field
	Type   string
}

func (e *MalformedError) Error() string {
	return fmt.Sprintf("malformed %s", e.Type)
}

type RelationalError struct {
	Target    *Field
	Reference interface{}
	Operator  string
}

func (e *RelationalError) Error() string {
	m := map[string]string{
		">":  "not greater than",
		">=": "neither greater than nor equal to",
		"==": "not equal to",
		"<":  "not less than",
		"<=": "neither less than nor equal to",
		"!=": "equal to",
	}
	ref := e.Reference
	if f, ok := e.Reference.(*Field); ok {
		ref = f.ID
	}
	return fmt.Sprintf("is %s %v", e.Target, m[e.Operator], ref)
}
