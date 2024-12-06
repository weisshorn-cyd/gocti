package api_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weisshorn-cyd/gocti/api"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Int    int       `gocti:"int"`
		Float  float64   `gocti:"float"`
		String string    `gocti:"string"`
		Bool   bool      `gocti:"bool"`
		Array  []string  `gocti:"array"`
		Time   time.Time `gocti:"time"`
	}

	type args struct {
		input  any
		output any
	}

	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "Basic types OK",
			args: args{
				input: map[string]any{
					"int":    1,
					"float":  0.54,
					"string": "test",
					"bool":   true,
					"array":  []any{"a", "b"},
					"time":   "2024-07-09T08:54:17.503Z",
				},
				output: &testStruct{},
			},
			want: &testStruct{
				Int:    1,
				Float:  0.54,
				String: "test",
				Bool:   true,
				Array:  []string{"a", "b"},
				Time:   time.Date(2024, 7, 9, 8, 54, 17, 503000000, time.UTC),
			},
			wantErr: false,
		}, {
			name: "Decoder creation failure",
			args: args{
				input:  map[string]any{"int": 1},
				output: testStruct{},
			},
			want:    testStruct{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := api.Decode(test.args.input, test.args.output)
			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, test.args.output)
		})
	}
}

type Interface struct {
	F1 string `gocti:"f1"`
	F2 string `gocti:"f2"`

	Remain map[string]any `gocti:",remain"`
}

func (i Interface) Implementations() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[Implementation](),
	}
}

func (i Interface) Remainder() map[string]any {
	return i.Remain
}

type Implementation struct {
	F1 string `gocti:"f1"`
	F2 string `gocti:"f2"`
	F3 string `gocti:"f3"`
	F4 string `gocti:"f4"`
}

type NotImplementation struct {
	F1 string `gocti:"f1"`
	F2 string `gocti:"f2"`
	F3 string `gocti:"f3"`
	F4 string `gocti:"f4"`
}

func TestDecodeInterface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   api.GraphQLInterface
		output  any
		want    any
		wantErr error
	}{
		{
			name: "Decode Interface - Not implementation",
			input: Interface{
				F1:     "v1",
				F2:     "v2",
				Remain: map[string]any{},
			},
			output: &NotImplementation{},
			want:   nil,
			wantErr: api.NotImplementingError{
				InterfaceType:      reflect.TypeFor[Interface](),
				ImplementationType: reflect.TypeFor[NotImplementation](),
			},
		}, {
			name: "Decode Interface - Not a pointer",
			input: Interface{
				F1:     "v1",
				F2:     "v2",
				Remain: map[string]any{},
			},
			output:  Implementation{},
			want:    nil,
			wantErr: api.ErrNotAPointer,
		}, {
			name: "Decode Interface - No remain",
			input: Interface{
				F1:     "v1",
				F2:     "v2",
				Remain: map[string]any{},
			},
			output: &Implementation{},
			want: &Implementation{
				F1: "v1",
				F2: "v2",
			},
			wantErr: nil,
		}, {
			name: "Decode Interface - Ok",
			input: Interface{
				F1: "v1",
				F2: "v2",
				Remain: map[string]any{
					"f3": "v3",
				},
			},
			output: &Implementation{},
			want: &Implementation{
				F1: "v1",
				F2: "v2",
				F3: "v3",
			},
			wantErr: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := api.DecodeInterface(test.input, test.output)
			if test.wantErr != nil {
				require.ErrorContains(t, err, test.wantErr.Error())

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, test.output)
		})
	}
}
