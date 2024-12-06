package api

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sampleInterface struct {
	P1 string `gocti:"p1"`
	P2 string `gocti:"p2"`
}

func (i sampleInterface) Implementations() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[sampleImplementation](),
	}
}

func (i sampleInterface) Remainder() map[string]any {
	return map[string]any{}
}

type sampleImplementation struct {
	C1 string `gocti:"c1"`
	C2 string `gocti:"c2"`
}

func TestParseAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		parseFunc func(options ...ParseOption) string
		want      string
	}{
		{
			name:      "Parse - Not a struct",
			parseFunc: ParseAttributes[float64],
			want:      "",
		}, {
			name: "Parse - No tags",
			parseFunc: ParseAttributes[struct {
				Field1 string
				Field2 int
			}],
			want: "",
		}, {
			name: "Parse ok",
			parseFunc: ParseAttributes[struct {
				Field1 string `gocti:"field1"`
				Field2 int    `gocti:"field2"`
			}],
			want: "field1\nfield2",
		}, {
			name: "Parse ok - SubStruct",
			parseFunc: ParseAttributes[struct {
				Field1 string `gocti:"field1"`
				Field2 struct {
					SubField1 string `gocti:"sub_field1"`
				} `gocti:"field2"`
			}],
			want: "field1\nfield2 {\nsub_field1\n}",
		}, {
			name: "Parse ok - SubStruct + untagged",
			parseFunc: ParseAttributes[struct {
				Field1 string `gocti:"field1"`
				Field2 struct {
					SubField1 string `gocti:"sub_field1"`
				} `gocti:"field2"`
				UntaggedField1 int64
			}],
			want: "field1\nfield2 {\nsub_field1\n}",
		}, {
			name:      "Parse ok - interface",
			parseFunc: ParseAttributes[sampleInterface],
			want:      "p1\np2\n... on sampleImplementation {\nc1\nc2\n}",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.parseFunc()
			got = strings.ReplaceAll(got, "\t", "")

			assert.Equal(t, test.want, got)
		})
	}
}
