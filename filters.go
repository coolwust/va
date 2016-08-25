package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type fieldError struct {
	target *field
	want   reflect.Type
}

func (e *fieldError) Error() string {
	return ""
}

func fieldFilter() ([]filter, error) {
	regexps := make(map[string]*regexp.Regexp)
	for name, re := range filterRegexps {
		regexps[name] = regexp.MustCompile(re)
	}
	filters := make([]filter, 0)
	for _, target := range fields {
		for _, option := range target.options {
			for name, re := range regexps {
				matches := re.FindStringSubmatch(option)
				if len(matches) == 0 {
					continue
				}
				var filter filter
				var err error
				switch name {
				case FilterEmail:
					filter, err = caseEmail(target)
				case FilterSize:
					filter, err = caseSize(target, matches, fields)
				}
				if err != nil {
					return nil, err
				}
				filters = append(filters, filter)
				goto matched
			}
			return nil, fmt.Errorf("unknown tag option %s", option)
		matched:
		}
	}
	return filters, nil
}

func caseSize(target *field, matches []string, fields []*field) (*size, error) {
	if !target.rt.ConvertibleTo(rtFloat64) {
		return nil, &fieldError{target, rtFloat64}
	}
	n := &size{target: target, op: matches[1]}
	if ref, err := strconv.ParseFloat(matches[2], 64); err == nil {
		n.ref = ref
		return n, nil
	} else if err != nil && err.(*strconv.NumError).Err == strconv.ErrRange {
		return nil, fmt.Errorf("%v", err)
	}
	for _, f := range fields {
		if f.id == matches[2] {
			if !f.rt.ConvertibleTo(rtFloat64) {
				return nil, &fieldError{f, rtString}
			}
			n.ref = f
			return n, nil
		}
	}
	str := "field %s is convertible to neither an float64 sizeber nor a field ID"
	return nil, fmt.Errorf(str, matches[2])
}

func caseEmail(target *field) (*email, error) {
	if target.rt.Kind() != reflect.String {
		return nil, &fieldError{target, rtString}
	}
	return &email{target: target}, nil
}

type filter interface {
	filter(structrv reflect.Value) error
}

type email struct {
	target   []int
	id       string
	optional bool
}

func (e *email) filter(structrv reflect.Value) error {
	target := structrv.FieldByIndex(e.target.index).Interface().(string)
	if !IsEmail(target) {
		return &MalformedError{
			Target: &Field{ID: e.target.id, Index: e.target.index, Value: target},
			Type:   "email",
		}
	}
	return nil
}

//type cap struct {
//	target *field
//	ref    interface{}
//	op     string
//}
//
//func (c *cap) filter(structrv reflect.Value) error {
//	return nil
//}
//
//type len struct {
//	target *field
//	ref    interface{}
//	op     string
//}
//
//func (l *len) filter(structrv reflect.Value) error {
//	return nil
//}
