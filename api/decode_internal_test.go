package api

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_goctiHookFunc(t *testing.T) {
	t.Parallel()

	type args struct {
		fromType reflect.Type
		toType   reflect.Type
		data     any
	}

	tests := []struct {
		name    string
		args    args
		want    any
		wantErr error
	}{
		{
			name: "DateTime (string) to time.Time",
			args: args{
				fromType: reflect.TypeFor[string](),
				toType:   reflect.TypeFor[time.Time](),
				data:     "2024-07-09T08:54:17.503Z",
			},
			want:    time.Date(2024, 7, 9, 8, 54, 17, 503000000, time.UTC),
			wantErr: nil,
		}, {
			name: "Error DateTime (int) to time.Time",
			args: args{
				fromType: reflect.TypeFor[int](),
				toType:   reflect.TypeFor[time.Time](),
				data:     int(1720515257),
			},
			want: time.Date(2024, 7, 9, 8, 54, 17, 0, time.UTC),
			wantErr: UnimplementedDecodingError{
				fromType: reflect.TypeFor[int](),
				toType:   reflect.TypeFor[time.Time](),
				data:     int(1720515257),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := goctiHookFunc(test.args.fromType, test.args.toType, test.args.data)
			if test.wantErr != nil {
				require.ErrorContains(t, err, test.wantErr.Error())

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}
