package list_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weisshorn-cyd/gocti/list"
)

//nolint:gochecknoglobals // Test parameters
var (
	// Test filters.
	filterEmpty = list.Filter{}
	filterFull  = list.Filter{
		Mode:     list.FilterModeAnd,
		Key:      []string{"k1", "k2"},
		Operator: list.FilterOperatorEq,
		Values:   []any{"v1", 2, false},
	}
	filterNoMode = list.Filter{
		Key:      []string{"k1", "k2"},
		Operator: list.FilterOperatorEq,
		Values:   []any{"v1", 2, false},
	}
	filterNoOperator = list.Filter{
		Mode:   list.FilterModeAnd,
		Key:    []string{"k1", "k2"},
		Values: []any{"v1", 2, false},
	}
	filterNoValues = list.Filter{
		Mode:     list.FilterModeAnd,
		Key:      []string{"k1", "k2"},
		Operator: list.FilterOperatorEq,
		Values:   []any{},
	}

	// Test filterGroups.
	filterGroupEmpty = list.FilterGroup{}
	filterGroupFull  = list.FilterGroup{
		Mode:         list.FilterModeAnd,
		Filters:      []list.Filter{filterFull},
		FilterGroups: []list.FilterGroup{},
	}
	filterGroupFullFull = list.FilterGroup{
		Mode:         list.FilterModeAnd,
		Filters:      []list.Filter{filterFull},
		FilterGroups: []list.FilterGroup{filterGroupFull},
	}
	filterGroupNoMode = list.FilterGroup{
		Filters:      []list.Filter{filterFull},
		FilterGroups: []list.FilterGroup{},
	}
	filterGroupFilterEmpty = list.FilterGroup{
		Mode:         list.FilterModeAnd,
		Filters:      []list.Filter{filterEmpty},
		FilterGroups: []list.FilterGroup{},
	}
	filterGroupFilterGroupNoMode = list.FilterGroup{
		Mode:         list.FilterModeAnd,
		Filters:      []list.Filter{filterFull},
		FilterGroups: []list.FilterGroup{filterGroupNoMode},
	}
)

func TestFilter_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		filter  list.Filter
		wantErr []error
	}{
		{
			name:    "valid",
			filter:  filterFull,
			wantErr: nil,
		}, {
			name:    "no mode",
			filter:  filterNoMode,
			wantErr: []error{list.ErrEmptyFilterMode},
		}, {
			name:    "no operator",
			filter:  filterNoOperator,
			wantErr: []error{list.ErrEmptyFilterOperator},
		}, {
			name:    "no values",
			filter:  filterNoValues,
			wantErr: []error{list.ErrEmptyFilterValues},
		}, {
			name:   "empty",
			filter: filterEmpty,
			wantErr: []error{
				list.ErrEmptyFilterMode,
				list.ErrEmptyFilterOperator,
				list.ErrEmptyFilterValues,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.filter.IsValid()
			for _, erri := range test.wantErr {
				require.ErrorIs(t, err, erri)
			}

			if test.wantErr == nil {
				require.NoError(t, err)
			}
		})
	}
}

func TestFilterGroup_IsValid(t *testing.T) {
	tests := []struct {
		name        string
		filterGroup list.FilterGroup
		wantErr     []error
	}{
		{
			name:        "valid",
			filterGroup: filterGroupFull,
			wantErr:     nil,
		}, {
			name:        "full depth",
			filterGroup: filterGroupFullFull,
			wantErr:     nil,
		}, {
			name:        "no mode",
			filterGroup: filterGroupNoMode,
			wantErr:     []error{list.ErrEmptyFilterMode},
		}, {
			name:        "missing internal fields",
			filterGroup: filterGroupFilterGroupNoMode,
			wantErr:     []error{list.ErrEmptyFilterMode},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.filterGroup.IsValid()
			for _, erri := range test.wantErr {
				require.ErrorIs(t, err, erri)
			}

			if test.wantErr == nil {
				require.NoError(t, err)
			}
		})
	}
}

func TestFilter_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		filter  list.Filter
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			filter:  filterFull,
			want:    `{"key":["k1","k2"],"values":["v1",2,false],"operator":"eq","mode":"and"}`,
			wantErr: false,
		}, {
			name:    "empty",
			filter:  filterEmpty,
			want:    "null",
			wantErr: false,
		}, {
			name:    "no mode",
			filter:  filterNoMode,
			want:    "",
			wantErr: true,
		}, {
			name:    "no operator",
			filter:  filterNoOperator,
			want:    "",
			wantErr: true,
		}, {
			name:    "no values",
			filter:  filterNoValues,
			want:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.filter.MarshalJSON()
			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestFilterGroup_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		filterGroup list.FilterGroup
		want        string
		wantErr     bool
	}{
		{
			name:        "valid",
			filterGroup: filterGroupFull,
			want: `{"mode":"and","filters":[{"key":["k1","k2"],"values":["v1",2,false],` +
				`"operator":"eq","mode":"and"}],"filterGroups":[]}`,
			wantErr: false,
		}, {
			name:        "full depth",
			filterGroup: filterGroupFullFull,
			want: `{"mode":"and","filters":[{"key":["k1","k2"],"values":["v1",2,false],` +
				`"operator":"eq","mode":"and"}],"filterGroups":[{"mode":"and","filters":[{"key":["k1","k2"],` +
				`"values":["v1",2,false],"operator":"eq","mode":"and"}],"filterGroups":[]}]}`,
			wantErr: false,
		}, {
			name:        "empty",
			filterGroup: filterGroupEmpty,
			want:        "null",
			wantErr:     false,
		}, {
			name:        "no mode",
			filterGroup: filterGroupNoMode,
			want:        "",
			wantErr:     true,
		}, {
			name:        "missing internal filter fields",
			filterGroup: filterGroupFilterEmpty,
			want:        "",
			wantErr:     true,
		}, {
			name:        "missing internal filterGroup fields",
			filterGroup: filterGroupFilterGroupNoMode,
			want:        "",
			wantErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.filterGroup.MarshalJSON()
			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestFilter_ToGraphQL(t *testing.T) {
	tests := []struct {
		name    string
		filter  list.Filter
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			filter:  filterFull,
			want:    `{mode:and, key:["k1","k2"], operator:eq, values:["v1",2,false]}`,
			wantErr: false,
		}, {
			name:    "empty",
			filter:  filterEmpty,
			want:    ``,
			wantErr: true,
		}, {
			name:    "no mode",
			filter:  filterNoMode,
			want:    ``,
			wantErr: true,
		}, {
			name:    "no operator",
			filter:  filterNoOperator,
			want:    ``,
			wantErr: true,
		}, {
			name:    "no values",
			filter:  filterNoValues,
			want:    ``,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.filter.ToGraphQL()
			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestFilterGroup_ToGraphQL(t *testing.T) {
	tests := []struct {
		name        string
		filterGroup list.FilterGroup
		want        string
		wantErr     bool
	}{
		{
			name:        "valid",
			filterGroup: filterGroupFull,
			want: `{mode:and, filters:[{mode:and, key:["k1","k2"], ` +
				`operator:eq, values:["v1",2,false]}], filterGroups:[]}`,
			wantErr: false,
		}, {
			name:        "full depth",
			filterGroup: filterGroupFullFull,
			want: `{mode:and, filters:[{mode:and, key:["k1","k2"], ` +
				`operator:eq, values:["v1",2,false]}], filterGroups:[{mode:and, filters:[{mode:and, key:["k1","k2"], ` +
				`operator:eq, values:["v1",2,false]}], filterGroups:[]}]}`,
			wantErr: false,
		}, {
			name:        "empty",
			filterGroup: filterGroupEmpty,
			want:        ``,
			wantErr:     true,
		}, {
			name:        "no mode",
			filterGroup: filterGroupNoMode,
			want:        ``,
			wantErr:     true,
		}, {
			name:        "missing internal filter fields",
			filterGroup: filterGroupFilterEmpty,
			want:        ``,
			wantErr:     true,
		}, {
			name:        "missing internal filterGroup fields",
			filterGroup: filterGroupFilterGroupNoMode,
			want:        ``,
			wantErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.filterGroup.ToGraphQL()
			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}
