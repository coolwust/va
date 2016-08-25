package validator

import (
	"reflect"
	"testing"
)

type F1 struct{ A string }

type F2 struct{ A, B int }

var findFiltersTests = []struct {
	fields []*field
	counts int
}{
	{
		[]*field{
			&field{rt: rtString, options: []string{"email"}},
			&field{rt: rtInt, options: []string{"size>=10"}},
		},
		2,
	},
	{
		[]*field{
			&field{id: "a", rt: rtInt, options: []string{"size>=10"}},
			&field{id: "b", rt: rtInt, options: []string{"size>=a"}},
		},
		2,
	},
}

func TestFindFilters(t *testing.T) {
	for i, tt := range findFiltersTests {
		filters, err := findFilters(tt.fields)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		if counts := len(filters); counts != tt.counts {
			t.Errorf("%d: counts = %d, want %d", counts, tt.counts)
		}
	}
}

var emailFilterTests = []struct {
	target   *field
	structrv reflect.Value
	ok       bool
}{
	{&field{index: []int{0}}, reflect.ValueOf(F1{"kevin@example.com"}), true},
	{&field{index: []int{0}}, reflect.ValueOf(F1{"example.com"}), false},
}

func TestEmailFilter(t *testing.T) {
	for i, tt := range emailFilterTests {
		e := &email{target: tt.target}
		if err := e.filter(tt.structrv); (err == nil) != tt.ok {
			t.Errorf("%d: ok == %t, want %t", i, err == nil, tt.ok)
		}
	}
}

var sizeFilterTests = []struct {
	target   *field
	ref      interface{}
	op       string
	structrv reflect.Value
	ok       bool
}{
	{&field{index: []int{0}}, 10.0, ">=", reflect.ValueOf(F2{A: 15}), true},
	{&field{index: []int{0}}, &field{index: []int{1}}, "<", reflect.ValueOf(F2{6, 12}), true},
	{&field{index: []int{0}}, &field{index: []int{1}}, "==", reflect.ValueOf(F2{10, 11}), false},
}

func TestSizeFilter(t *testing.T) {
	for i, tt := range sizeFilterTests {
		s := &size{target: tt.target, ref: tt.ref, op: tt.op}
		if err := s.filter(tt.structrv); (err == nil) != tt.ok {
			t.Errorf("%d: ok == %t, want %t", i, err == nil, tt.ok)
		}
	}
}
