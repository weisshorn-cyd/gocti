package list

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-json"
)

var (
	ErrEmptyFilterMode     = errors.New("a filter mode must be set")
	ErrEmptyFilterOperator = errors.New("a filter operator must be set")
	ErrEmptyFilterValues   = errors.New("at least one value must be set")
)

// OrderingMode defines the various ordering modes available when retrieving multiple entities.
type OrderingMode string

const (
	OrderModeAsc  OrderingMode = "asc"
	OrderModeDesc OrderingMode = "desc"
)

// FilterMode defines the various filter combination modes available.
type FilterMode string

const (
	FilterModeAnd FilterMode = "and"
	FilterModeOr  FilterMode = "or"
)

// FilterOperator defines the various filtering operators available.
type FilterOperator string

const (
	FilterOperatorEq            FilterOperator = "eq"
	FilterOperatorNotEq         FilterOperator = "not_eq"
	FilterOperatorLt            FilterOperator = "lt"
	FilterOperatorLte           FilterOperator = "lte"
	FilterOperatorGt            FilterOperator = "gt"
	FilterOperatorGte           FilterOperator = "gte"
	FilterOperatorMatch         FilterOperator = "match"
	FilterOperatorWildcard      FilterOperator = "wildcard"
	FilterOperatorContains      FilterOperator = "contains"
	FilterOperatorNotContains   FilterOperator = "not_contains"
	FilterOperatorEndsWith      FilterOperator = "ends_with"
	FilterOperatorNotEndsWith   FilterOperator = "not_ends_with"
	FilterOperatorStartsWith    FilterOperator = "starts_with"
	FilterOperatorNotStartsWith FilterOperator = "not_starts_with"
	FilterOperatorScript        FilterOperator = "script"
	FilterOperatorNil           FilterOperator = "nil"
	FilterOperatorNotNil        FilterOperator = "not_nil"
	FilterOperatorSearch        FilterOperator = "search"
)

type Filter struct {
	Key      []string       `json:"key"`
	Values   []any          `json:"values"`
	Operator FilterOperator `json:"operator"`
	Mode     FilterMode     `json:"mode"`
}

type FilterGroup struct {
	Mode         FilterMode    `json:"mode"`
	Filters      []Filter      `json:"filters"`
	FilterGroups []FilterGroup `json:"filterGroups"`
}

// IsValid returns any error related to missing Filter content.
func (filter Filter) IsValid() error {
	errs := []error{}

	if filter.Mode == "" {
		errs = append(errs, ErrEmptyFilterMode)
	}

	if filter.Operator == "" {
		errs = append(errs, ErrEmptyFilterOperator)
	}

	if len(filter.Values) == 0 {
		errs = append(errs, ErrEmptyFilterValues)
	}

	return errors.Join(errs...)
}

// IsValid returns any error related to missing FilterGroup content.
func (filterGroup FilterGroup) IsValid() error {
	errs := []error{}

	if filterGroup.Mode == "" {
		errs = append(errs, ErrEmptyFilterMode)
	}

	for _, filter := range filterGroup.Filters {
		errs = append(errs, filter.IsValid())
	}

	for _, filterG := range filterGroup.FilterGroups {
		errs = append(errs, filterG.IsValid())
	}

	return errors.Join(errs...)
}

// MarshalJSON checks for Filter missing content and returns its json representation.
func (filter Filter) MarshalJSON() ([]byte, error) {
	// Check for empty fields
	emptyKey := len(filter.Key) == 0
	emptyValues := len(filter.Values) == 0
	emptyOperator := filter.Operator == ""
	emptyMode := filter.Mode == ""

	// Return 'null' if filter is empty
	if emptyKey && emptyValues && emptyOperator && emptyMode {
		return []byte("null"), nil // Return 'null' if filter is empty
	}

	// Check Filter content
	if err := filter.IsValid(); err != nil {
		return nil, fmt.Errorf("filter has invalid content: %w", err)
	}

	// Marshalling key slice
	key, err := marshalList(filter.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal key '%s': %w", filter.Key, err)
	}

	// Marshalling values slice
	values, err := marshalList(filter.Values)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal values '%s': %w", filter.Values, err)
	}

	// Marshalling operator
	operator := `"` + filter.Operator + `"`

	// Marshalling mode
	mode := `"` + filter.Mode + `"`

	return []byte(fmt.Sprintf(`{"key":%s,"values":%s,"operator":%s,"mode":%s}`, key, values, operator, mode)), nil
}

// MarshalJSON checks for FilterGroup missing content and returns its json representation.
func (filterGroup FilterGroup) MarshalJSON() ([]byte, error) {
	// Check for empty fields
	emptyMode := filterGroup.Mode == ""
	emptyFilters := len(filterGroup.Filters) == 0
	emptyFilterGroups := len(filterGroup.FilterGroups) == 0

	// Return 'null' if FilterGroup is empty
	if emptyMode && emptyFilters && emptyFilterGroups {
		return []byte("null"), nil
	}

	// Check FilterGroup content
	if err := filterGroup.IsValid(); err != nil {
		return nil, fmt.Errorf("filterGroup has invalid content: %w", err)
	}

	// Marshalling mode
	mode := `"` + filterGroup.Mode + `"`

	// Marshalling filters slice
	filters, err := marshalList(filterGroup.Filters)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal filters '%s': %w", filterGroup.Filters, err)
	}

	// Marshalling filterGroups slice
	filterGroups, err := marshalList(filterGroup.FilterGroups)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal filterGroups '%s': %w", filterGroup.FilterGroups, err)
	}

	return []byte(fmt.Sprintf(`{"mode":%s,"filters":%s,"filterGroups":%s}`, mode, filters, filterGroups)), nil
}

func marshalList[T any](list []T) ([]byte, error) {
	// If empty, return empty list (instead of 'null' if not initialized)
	if len(list) == 0 {
		return []byte("[]"), nil
	}

	//nolint:wrapcheck // Helper function
	return json.Marshal(list)
}

// ToGraphQL returns the appropriate representation of a Filter for direct use in a 'customAttributes' field.
// In the GraphQL query, the field names and typed values are not quoted, but strings are.
func (filter Filter) ToGraphQL() (string, error) {
	if err := filter.IsValid(); err != nil {
		return "", fmt.Errorf("filter has invalid content: %w", err)
	}

	key, err := marshalList(filter.Key)
	if err != nil {
		return "", fmt.Errorf("cannot encode filter values to GraphQL: %w", err)
	}

	values, err := marshalList(filter.Values)
	if err != nil {
		return "", fmt.Errorf("cannot encode filter values to GraphQL: %w", err)
	}

	return fmt.Sprintf(
		`{mode:%s, key:%s, operator:%s, values:%s}`,
		filter.Mode,
		key,
		filter.Operator,
		values,
	), nil
}

// ToGraphQL returns the appropriate representation of a FilterGroup for direct use in a 'customAttributes' field.
func (filterGroup FilterGroup) ToGraphQL() (string, error) {
	if err := filterGroup.IsValid(); err != nil {
		return "", fmt.Errorf("filterGroup has invalid content: %w", err)
	}

	filters := make([]string, len(filterGroup.Filters))

	for i, filt := range filterGroup.Filters {
		filtStr, err := filt.ToGraphQL()
		if err != nil {
			return "", fmt.Errorf("cannot encode filter to GraphQL: %w", err)
		}

		filters[i] = filtStr
	}

	filterGroups := make([]string, len(filterGroup.FilterGroups))

	for i, filtG := range filterGroup.FilterGroups {
		filtGStr, err := filtG.ToGraphQL()
		if err != nil {
			return "", fmt.Errorf("cannot encode filterGroup to GraphQL: %w", err)
		}

		filterGroups[i] = filtGStr
	}

	return fmt.Sprintf(
		`{mode:%s, filters:%s, filterGroups:%s}`,
		filterGroup.Mode,
		"["+strings.Join(filters, ",")+"]",
		"["+strings.Join(filterGroups, ",")+"]",
	), nil
}
