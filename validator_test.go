package validator

import (
	"reflect"
	"testing"
)

type V0 struct {
	A string
	B int
	C struct {
		D string
		E int
		F struct {
			G string
			H int
		}
	}
}

var fieldIndexesTest = [][]int{
	[]int{0},
	[]int{1},
	[]int{2, 0},
	[]int{2, 1},
	[]int{2, 2, 0},
	[]int{2, 2, 1},
}

func TestFieldIndexes(t *testing.T) {
	if indexes := fieldIndexes(reflect.TypeOf(V0{})); !reflect.DeepEqual(indexes, fieldIndexesTest) {
		t.Errorf("indexes = %#v, want %#v", indexes, fieldIndexesTest)
	}
}

type V1 struct {
	Email string `va:"email"`
	Age   int    `va:"optional,size>=18"`
}

func TestValidatorValidate(t *testing.T) {
}
