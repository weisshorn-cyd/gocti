package api

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	inProgress = "PROCESSING"
	maxDepth   = 10
)

// GraphQLInterface is a type that represents a GraphQL Interface or Union.
type GraphQLInterface interface {
	// The [Implementations] method should return all types that implement this GraphQL interface.
	Implementations() []reflect.Type
	// The [Remainder] method should return the additional data to be stored in an implementation's
	// specific fields. In GraphQL, Interfaces (and Unions) have less fields than their
	// implementations. The additional data can be saved separately when decoding a GraphQL
	// response into a Golang struct by adding a extra field using the tag `gocti:",remain"`
	// The [Remainder] method allows [DecodeInterface] to retrieve this extra data for conversion.
	Remainder() map[string]any
}

type attributeParser struct {
	inProgress string
	maxDepth   int
	comments   bool
}

func newAttributeParser() attributeParser {
	return attributeParser{
		inProgress: inProgress,
		maxDepth:   maxDepth,
		comments:   true,
	}
}

type ParseOption func(p *attributeParser)

// WithMaxDepth sets the maximum depth at which the parser stops
// recursively exploring struct fields.
func WithMaxDepth(maxDepth int) ParseOption {
	return func(p *attributeParser) {
		p.maxDepth = maxDepth
	}
}

// WithComments sets whether the parser should signal ignored recursion loops
// or exceeded maximum depth fields with a comment.
func WithComments(comments bool) ParseOption {
	return func(p *attributeParser) {
		p.comments = comments
	}
}

// ParseAttributes parses a struct type to build a GraphQL query segment matching its fields.
// The provided type must have a structure compatible with the GraphQL schema.
// Only fields marked with the "gocti" struct tag will be parsed.
// Field arguments can be added after the name tag (e.g. `gocti:"name,(id: \"0123456789\")"`).
// They should be added in second position, after the name, and before any decoding-specific
// tag value, like "remain", "omitempty", or "squash".
//
// WARNING: Using ParseAttributes with large structs pointing to many other structs
// (typically non-truncated OpenCTI entities) can quickly lead to a very large returned string.
func ParseAttributes[Entity any](options ...ParseOption) string {
	parser := newAttributeParser()
	for _, opt := range options {
		opt(&parser)
	}

	entityType := reflect.TypeFor[Entity]()
	attributes := parser.parseAttributes(entityType, 0, map[reflect.Type]string{})

	// Triming leading and trailing curly braces for query integration.
	attributes = strings.Trim(attributes, "{}")
	attributes = strings.Trim(attributes, "\n")

	return FormatAttributes(attributes, 0)
}

func FormatAttributes(attr string, offset int) string {
	attr = strings.ReplaceAll(attr, "\t", "")
	attr = strings.ReplaceAll(attr, "  ", "")

	if offset < 0 {
		offset = 0
	}

	output := ""

	lines := strings.Split(attr, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "\t ")

		switch {
		default:
			output += strings.Repeat("\t", offset) + line + "\n"
		case len(line) > 0 && line[len(line)-1] == '{':
			output += strings.Repeat("\t", offset) + line + "\n"
			offset++
		case len(line) > 0 && line[len(line)-1] == '}':
			offset--
			output += strings.Repeat("\t", offset) + line + "\n"
		}
	}

	output = strings.Trim(output, "\n")

	return output
}

// parseAttributes is the internal, recursive, parser for [ParseAttributes].
func (p attributeParser) parseAttributes(entityType reflect.Type, depth int, processed map[reflect.Type]string) string {
	// Parse structs only.
	if entityType.Kind() != reflect.Struct {
		return ""
	}

	// Manage state.
	alreadyProcessed, exists := processed[entityType]

	switch {
	case !exists:
		// Create new entry in cache, signaling beginning of processing for this type.
		processed[entityType] = p.inProgress
	case exists && alreadyProcessed != p.inProgress:
		// Return stored result.
		return alreadyProcessed
	}

	// List of all fields for this type.
	graphQLFields := []string{"{"}
	graphQLCommentFields := []string{}

	// Struct field analysis (recursive).
	for i := range entityType.NumField() {
		field := entityType.Field(i)
		graphQLField := p.parseFieldAttributes(field, depth, processed)

		if graphQLField == "" {
			continue
		}

		if graphQLField[0] == '#' {
			graphQLCommentFields = append(graphQLCommentFields, graphQLField)

			continue
		}

		graphQLFields = append(graphQLFields, graphQLField)
	}

	// Check if Type is a GraphQL Interface, then add implementation attributes.
	if graphQLInterface, ok := reflect.New(entityType).Interface().(GraphQLInterface); ok {
		for _, impl := range graphQLInterface.Implementations() {
			implAttr := p.parseAttributes(impl, depth+1, processed)
			graphQLFields = append(graphQLFields, "... on "+impl.Name()+" "+implAttr)
		}
	}

	if len(graphQLFields) == 1 {
		processed[entityType] = "" // Store the result in temp cache.

		return ""
	}

	attr := strings.Join(append(graphQLFields, graphQLCommentFields...), "\n") + "\n}"

	processed[entityType] = attr // Store the result in temp cache.

	return attr
}

// parseFieldAttributes parses a single struct field for [parseAttributes].
func (p attributeParser) parseFieldAttributes(
	field reflect.StructField,
	depth int,
	processed map[reflect.Type]string,
) string {
	graphQLField := parseFieldTag(field.Tag.Get("gocti"))

	if graphQLField == "" {
		return ""
	}

	// Go Deeper in fields that are structs themselves or slices.
	var parseFurther reflect.Type

	//nolint:exhaustive // Default case does its job.
	switch field.Type.Kind() {
	case reflect.Struct:
		parseFurther = field.Type
	case reflect.Slice, reflect.Pointer:
		parseFurther = field.Type.Elem()
	default:
	}

	switch {
	case parseFurther == nil:
		return graphQLField
	case depth > p.maxDepth:
		if p.comments {
			return fmt.Sprintf("# Ignored '%s' (%s): maximum depth exceeded", graphQLField, parseFurther.Name())
		}

		return ""
	case processed[parseFurther] == p.inProgress:
		if p.comments {
			return fmt.Sprintf("# Ignored '%s' (%s): recursion loop", graphQLField, parseFurther.Name())
		}

		return ""
	default:
		attr := p.parseAttributes(parseFurther, depth+1, processed)
		if attr != "" {
			return graphQLField + " " + attr
		}

		return ""
	}
}

// parseFieldTag returns the appropriate graphQLField, including arguments if present.
func parseFieldTag(goctiTag string) string {
	goctiTagValues := strings.Split(goctiTag, ",")
	goctiFieldName := goctiTagValues[0]
	goctiFieldArgs := ""

	if goctiFieldName == "" || goctiFieldName == "-" {
		return ""
	}

	if len(goctiTagValues) > 1 {
		goctiFieldArgs = goctiTagValues[1]
	}

	// Some reserved tag values (see [https://pkg.go.dev/github.com/go-viper/mapstructure/v2])
	// trigger specific behaviours when decoding. Those should not be added as field arguments.
	if goctiFieldArgs == "remain" || goctiFieldArgs == "omitempty" || goctiFieldArgs == "squash" {
		goctiFieldArgs = ""
	}

	graphQLField := goctiFieldName

	if goctiFieldArgs != "" {
		graphQLField += " " + goctiFieldArgs
	}

	return graphQLField
}
