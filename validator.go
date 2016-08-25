package validator

import (
	//"fmt"
	"reflect"
	"regexp"
	"strings"
)

const TagKey  = "valid"

const (
	optionID       = iota
	optionOptional
	optionEmail
	optionSize
)

var optionRegexps = map[int]*regexp.Regexp{
	optionID:       regexp.MustCompile(`^id=(-\w)$`),
	optionOptional: regexp.MustCompile(`^optional$`),
	optionEmail:    regexp.MustCompile(`^email$`),
	optionSize:     regexp.MustCompile(`^size(==|!=|>=|<=|>|<)([-\w]+)$`),
}

var (
	rtFloat64 = reflect.TypeOf(float64(0))
	rtInt     = reflect.TypeOf(0)
	rtString  = reflect.TypeOf("")
)

type Validator struct {
	filters []filter
	rt      reflect.Type
}

func New(v interface{}) *Validator {
	if rt := reflect.TypeOf(v); rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		panic("validator: expected of type pointer to struct")
	}
	validator := &Validator{rt: reflect.TypeOf(v)}
	rt := reflect.TypeOf(v).Elem()
	fields := make([]*field, 0)
	for _, index := range fieldIndexes(rt) {
		f := rt.FieldByIndex(index)
		tag, ok := f.Tag.Lookup(TagKey)
		if !ok || tag != "-" {
			continue
		}
		field := &field{id: strings.ToLower(f.Name), index: index, options: strings.Split(tag, ",")}
		if matches, options := extractOption(`optional`, field.options); matches != nil {
			field.options = options
			field.optional = true
		}
		if matches, options := extractOption(`id=([-\w]+)`, field.options); matches != nil {
			field.options = options
			field.id = matches[1]
		}
		fields = append(fields, field)
	}
	filters, err := findFilters(fields)
	if err != nil {
		panic("validator: " + err.Error())
	}
	return &Validator{filters: filters, rt: rt}
}

func extractOption(pattern string, options []string) ([]string, []string) {
	re := regexp.MustCompile(pattern)
	for i, option := range options {
		if matches := re.FindStringSubmatch(option); matches != nil {
			return matches, append(options[:i], options[i+1:]...)
		}
	}
	return nil, options
}

//func (v *Validator) ValidateNonEmpty() []error {
//}
//
func (v *Validator) Validate() []error {
}

func findIndexes(rt reflect.Type) (indexes [][]int) {
	indexes = make([][]int, 0)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		switch f.Type.Kind() {
		case reflect.Struct:
			for _, m := range fieldIndexes(f.Type) {
				indexes = append(indexes, append([]int{i}, m...))
			}
		default:
			indexes = append(indexes, f.Index)
		}
	}
	return
}

type field struct {
	id       string
	index    []int
	optional bool
	rt       reflect.Type
}

func findFilters(rt reflect.Type) []filter {
	for _, index := range fieldIndexes(rt) {
		tag, ok := rt.FieldByIndex(index).Tag.Lookup(TagKey)
		if !ok {
			continue
		}
		for _, option := range strings.Split(tag, ",") {
			for name, re := range optionRegexps {
				matches := re.FindStringSubmatch(option)
				if matches == nil {
					continue
				}
				switch name {
				case optionID:
					field.id = matches[1]
				case optionOptional:
					field.optional = true
				default:
					findFilter
				}

				filters = append(filters, filter)
				goto matched
			}
			return nil, fmt.Errorf("unknown tag option %s", option)
		matched:
		}
		field := &field{id: strings.ToLower(f.Name), index: index, options: }

	}
	filters := make([]filter, 0)
	for _, target := range fields {
		for _, option := range target.options {
		}
	}
	return filters, nil
}

type options []string

func (o *options) match(pattern string) []string {
	// TODO
}


func fieldID(field reflect.StructField) string {
}

func fieldOptional(field reflect.StructField) bool {
}

func fieldSize(field reflect.StructField) (filter, error) {
}
