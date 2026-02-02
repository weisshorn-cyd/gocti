from copy import copy
from dataclasses import dataclass
from enum import Enum
from string import Template

from pycti import OpenCTIApiClient
from rich import print

from gocti_type_generator import utils


class Kind(Enum):
    UNSET = "UNSET"
    # Most Common
    ENUM = "ENUM"
    SCALAR = "SCALAR"
    OBJECT = "OBJECT"
    INPUT_OBJECT = "INPUT_OBJECT"
    # Interfaces
    UNION = "UNION"
    INTERFACE = "INTERFACE"
    # Type Modifiers
    LIST = "LIST"
    NON_NULL = "NON_NULL"


@dataclass
class Field:
    name: str
    opencti_type: "Type"
    args: list["InputValue"]
    description: str = ""

    @classmethod
    def from_dict(cls, client: OpenCTIApiClient, data: dict[str, any]) -> "Field":
        return Field(
            name=data["name"],
            description=data["description"],
            opencti_type=Type.from_dict(client, data["type"]),
            args=[InputValue.from_dict(client, arg) for arg in data["args"]],
        )


@dataclass
class EnumValue:
    name: str
    description: str = ""

    @classmethod
    def from_dict(cls, data: dict[str, any]) -> "EnumValue":
        return EnumValue(name=data["name"], description=data["description"])


@dataclass
class InputValue:
    name: str
    opencti_type: "Type"
    description: str = ""
    default_value: str = ""

    @classmethod
    def from_dict(cls, client: OpenCTIApiClient, data: dict[str, any]) -> "InputValue":
        return InputValue(
            name=data["name"],
            description=data["description"],
            opencti_type=Type.from_dict(client, data["type"]),
            default_value=data["defaultValue"],
        )


MAX_NESTED_TYPE_DEPTH = 5
GRAPHQL_TYPES_PKG_NAME = "graphql"

ENUM_GO_TMPL = """
type $TypeName string

const(
    $EnumValues
)
""".strip()

STRUCT_GO_TMPL = """
type $TypeName struct{
    $Fields
}
""".strip()

INTERFACE_IMPL_GO_TMPL = """
func ($t $TypeName) Implementations() []reflect.Type {
    return []reflect.Type{
        $reflectImplTypes
    }
}

func ($t $TypeName) Remainder() map[string]any {
    return $t.Remain
}""".strip()

NULL_MARSHALL_GO_TMPL = """
func ($t $TypeName) MarshalJSON() ([]byte, error) {
	if reflect.ValueOf($t).IsZero() {
		return []byte("null"), nil
	}

	type temp$TypeName $TypeName

	//nolint:wrapcheck // To avoid complete implementation
	return json.Marshal((*temp$TypeName)(&$t))
}""".strip()


class Type:
    """Represents a GraphQL Type and its attributes.
    Instatiation should be made using the static method `from_name`."""

    verbose: bool = False

    # Dict cache of known types.
    known_types: dict[str, "Type"] = {}

    # Mapping between GraphQL and Golang types.
    # Can also be used to override any type.
    scalar_map: dict[str, str] = {
        "Any": "any",
        "Int": "int",
        "Float": "float64",
        "Boolean": "bool",
        "String": "string",
        "ID": "string",
        "StixId": "string",
        "StixRef": "string",
        "DateTime": "*time.Time",
        "Upload": "[]byte",
        "name_String_NotNull_minLength_2": "string",
        "context_String_NotNull_minLength_2": "string",
        "name_String_NotNull_minLength_1": "string",
        "algorithm_String_NotNull_minLength_3": "string",
        "hash_String_NotNull_minLength_5": "string",
        "source_name_String_NotNull_minLength_2": "string",
        "content_String_NotNull_minLength_2": "string",
        "user_email_String_NotNull_minLength_5_format_email": "string",
        "JSON": "string",
        "VocabularyCategory": "string",  # Enum with too long names for golangci-lint
    }

    # Will replace the given field's type of the given 'Type' with a custom type name.
    # Format: [TypeName].[field_type].[field_name] : [CustomTypeName]
    type_exceptions: dict[str, str] = {
        "IndicatorAddInput.input_field.pattern_type": "PatternType",
        "IndicatorAddInput.input_field.indicator_types": "[IndicatorType]",
        "IndicatorAddInput.input_field.x_mitre_platforms": "[Platform]",
        "ReportAddInput.input_field.report_types": "[ReportType]",
        "StixCyberObservableAddInput.input_field.type": "StixCyberObservableType",
        "GroupingAddInput.input_field.context": "GroupingContext",
        "OpinionAddInput.input_field.opinion": "OpinionType",
    }

    # Order in which interface fields are resolved and displayed in a type declaration.
    interface_resolve_order: list[str] = [
        "InternalObject",
        "BasicObject",
        "BasicRelationship",
        "StixObject",
        "StixCoreObject",
        "StixDomainObject",
        "StixMetaObject",
        "StixRelationship",
        "Container",
        "Case",
        "Location",
        "Identity",
        "ThreatActor",
        "StixCyberObservable",
        "HashedObservable",
    ]

    # List of fields that should be ignored in the Golang Type declaration.
    # If the type is an interface, its implementations will ignore it as well.
    ignored_fields_list: dict[str, list[str]] = {
        "InternalObject": [],
        "BasicObject": [],
        "BasicRelationship": ["fromRole", "toRole"],
        "StixObject": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "x_opencti_inferences",
        ],
        "StixCoreObject": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "toStix",
            "pendingFiles",
            "exportFiles",
            "editContext",
            "connectors",
            "jobs",
        ],
        "StixDomainObject": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "toStix",
            "pendingFiles",
            "exportFiles",
            "editContext",
            "connectors",
            "jobs",
            "x_opencti_graph_data",
            "avatar",
        ],
        "StixMetaObject": ["x_opencti_stix_ids", "is_inferred"],
        "StixRelationship": [
            "fromRole",
            "toRole",
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "x_opencti_inferences",
            "toStix",
        ],
        "Container": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "x_opencti_graph_data",
            "objectAssignee",
            "avatar",
            "authorized_members",
            "currentUserAccessRight",
            "relatedContainers",
        ],
        "Case": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "x_opencti_inferences",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "toStix",
            "pendingFiles",
            "exportFiles",
            "editContext",
            "connectors",
            "jobs",
            "x_opencti_graph_data",
            "objectAssignee",
            "avatar",
            "authorized_members",
            "currentUserAccessRight",
            "relatedContainers",
            "objectParticipant",
            "content_mapping",
        ],
        "Location": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "toStix",
            "pendingFiles",
            "exportFiles",
            "editContext",
            "connectors",
            "jobs",
            "x_opencti_graph_data",
            "objectAssignee",
            "avatar",
        ],
        "Identity": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "toStix",
            "pendingFiles",
            "exportFiles",
            "editContext",
            "connectors",
            "jobs",
            "x_opencti_graph_data",
            "objectAssignee",
            "avatar",
        ],
        "ThreatActor": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "x_opencti_inferences",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "toStix",
            "pendingFiles",
            "exportFiles",
            "editContext",
            "connectors",
            "jobs",
            "x_opencti_graph_data",
            "objectAssignee",
            "avatar",
        ],
        "StixCyberObservable": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "toStix",
            "pendingFiles",
            "exportFiles",
            "editContext",
            "connectors",
            "jobs",
        ],
        "HashedObservable": [
            "representative",
            "x_opencti_stix_ids",
            "is_inferred",
            "numberOfConnectedElement",
            "containersNumber",
            "containers",
            "reports",
            "notes",
            "opinions",
            "observedData",
            "groupings",
            "cases",
            "stixCoreRelationships",
            "stixCoreObjectsDistribution",
            "stixCoreRelationshipsDistribution",
            "toStix",
            "pendingFiles",
            "exportFiles",
            "editContext",
            "connectors",
            "jobs",
        ],
    }

    @classmethod
    def reset_known_types(cls) -> None:
        Type.known_types.clear()

    @classmethod
    def from_name(
        cls, client: OpenCTIApiClient, name: str, custom: bool = False
    ) -> "Type":
        """Retrieves the Type from the schema by name. The shorthand name can be used."""
        name, is_list, is_required, is_content_required = Type._extract_type_name(name)

        if name in Type.known_types:
            new_type = copy(Type.known_types[name])
        else:
            new_type = copy(Type(client=client, name=name, custom=custom))

        new_type.is_list = is_list
        new_type.is_required = is_required
        new_type.is_content_required = is_content_required

        return new_type

    @classmethod
    def from_dict(cls, client: OpenCTIApiClient, data: dict[str, any]) -> "Type":
        """Extracting the Type from a GraphQL nested Type map.
        The expected structure is: `{name, kind, ofType {name, kind, ofType {...}}}`"""

        names = utils.search_dict(data, "name")
        kinds = utils.search_dict(data, "kind")

        if len(names) != 1:
            raise ValueError(f"{names} is not unique")
        type_name = names[0]

        out = Type.from_name(client, name=type_name)

        # Go through Type modifier kinds
        for kind in kinds:
            if kind == Kind.NON_NULL.value:
                if out.is_list:
                    out.is_content_required = True
                else:
                    out.is_required = True
            elif kind == Kind.LIST.value:
                out.is_list = True

        return out

    @classmethod
    def name_from_dict(cls, data: dict[str, any]) -> str:
        """Extracting the Type name from a GraphQL nested Type map.
        The expected structure is: `{name, kind, ofType  {...}}`"""

        names = utils.search_dict(data, "name")
        kinds = utils.search_dict(data, "kind")

        if len(names) != 1:
            raise ValueError(f"{names} is not unique")
        type_name = names[0]

        # Go through Type modifier kinds
        for kind in kinds[::-1]:
            if kind == Kind.NON_NULL.value:
                type_name += "!"
            elif kind == Kind.LIST.value:
                type_name = f"[{type_name}]"

        return type_name

    @classmethod
    def query_schema(cls, client: OpenCTIApiClient, type_name: str) -> dict[str, any]:
        type_graphql = """query { __type (name: "$name"){
           name, kind, description,
           fields { $field_graphql },
           interfaces { $type_name_graphql }
           possibleTypes { $type_name_graphql }
           enumValues { $enum_value_graphql }
           inputFields { $input_value_graphql }
        }}"""

        # Data required to get the complete name of a Type: Types are nested to indicate lists, required fields, etc...
        temp = "name, kind, ofType { $type_name_graphql }"
        type_name_graphql = "name, kind"
        for _ in range(MAX_NESTED_TYPE_DEPTH):
            type_name_graphql = Template(temp).substitute(
                type_name_graphql=type_name_graphql,
            )

        # Data required to describe an Input Value.
        input_value_graphql = (
            "name, description, type { $type_name_graphql }, defaultValue"
        )
        input_value_graphql = Template(input_value_graphql).substitute(
            type_name_graphql=type_name_graphql,
        )

        # Data required to describe a Field.
        field_graphql = "name, description, type { $type_name_graphql }, args { $input_value_graphql }"
        field_graphql = Template(field_graphql).substitute(
            type_name_graphql=type_name_graphql,
            input_value_graphql=input_value_graphql,
        )

        # Data required to describe an Enum Value.
        enum_value_graphql = "name, description"

        query = Template(type_graphql).substitute(
            name=type_name,
            type_name_graphql=type_name_graphql,
            field_graphql=field_graphql,
            enum_value_graphql=enum_value_graphql,
            input_value_graphql=input_value_graphql,
        )

        schema_data = client.query(query)["data"]["__type"]
        if schema_data is None:
            raise ValueError(f"Unknown type '{type_name}'")

        return schema_data

    @classmethod
    def _extract_type_name(cls, name: str) -> tuple[str, bool, bool, bool]:
        """Removes the ! and [] shorthands from a Type name."""
        is_list, is_required, is_content_required = False, False, False

        if name.endswith("!"):
            is_required = True
            name = name[:-1]
        if name.startswith("["):
            is_list = True
            name = name[1:-1]
        if name.endswith("!"):
            is_content_required = True
            name = name[:-1]

        return name, is_list, is_required, is_content_required

    def __init__(
        self, client: OpenCTIApiClient, name: str, custom: bool = False
    ) -> None:
        # For all server interactions
        self.client: OpenCTIApiClient = client
        self.pkg = GRAPHQL_TYPES_PKG_NAME

        # __init__ is called only with unmodified Type names (no ! or [] shorthand)
        if "!" in name or "[" in name or "]" in name:
            raise ValueError(f"{name} is a modified Type name")

        # Create only once with __init__
        if name in Type.known_types:
            raise KeyError(f"Type {name} already resolved")
        Type.known_types[name] = self

        # General Type attributes
        self.name: str = name
        self.is_custom: bool = custom
        self.kind: Kind = Kind.UNSET
        self.description: str = ""

        # Different field attributes
        self.fields: dict[str, Field] = {}
        self.interfaces: dict[str, Type] = {}
        self.possible_types: dict[str, Type] = {}
        self.enum_values: dict[str, EnumValue] = {}
        self.input_fields: dict[str, InputValue] = {}
        self.custom_fields: dict[str, Field | InputValue] = {}

        # GraphQL type attributes (instance-specific)
        self.is_list: bool = False
        self.is_required: bool = False
        self.is_content_required: bool = False  # List content required: e.g. [Type!]

        # Intercept custom implementations before querying the server.
        if custom:
            Type.known_types.pop(name)
            return

        # Query the server schema
        schema_data = Type.query_schema(self.client, name)
        self.kind = (
            Kind(schema_data["kind"])
            if self.name not in self.scalar_map
            else Kind.SCALAR
        )
        self.description = schema_data["description"]

        # Scan field attributes
        self.__init_fields__(schema_data)
        self.__init_interfaces__(schema_data)
        self.__init_possible_types__(schema_data)
        self.__init_enum_values__(schema_data)
        self.__init_input_fields__(schema_data)

    def __str__(self) -> str:
        return self.__repr__()

    def __repr__(self) -> str:
        return self.graphql_type(full=True) + "(Type)"

    def __init_fields__(self, schema_data: dict[str, any]) -> None:
        fields = schema_data["fields"]
        if fields is None:
            return

        to_add = [Field.from_dict(self.client, field) for field in fields]
        self.add_fields(to_add)

    def __init_enum_values__(self, schema_data: dict[str, any]) -> None:
        enum_values = schema_data["enumValues"]
        if enum_values is None:
            return

        to_add = [EnumValue.from_dict(enum) for enum in enum_values]
        self.add_enum_values(to_add)

    def __init_interfaces__(self, schema_data: dict[str, any]) -> None:
        interfaces = schema_data["interfaces"]
        if interfaces is None:
            return

        to_add = [Type.from_dict(self.client, interface) for interface in interfaces]
        self.add_interfaces(to_add)

    def __init_possible_types__(self, schema_data: dict[str, any]) -> None:
        possible_types = schema_data["possibleTypes"]
        if possible_types is None:
            return

        to_add = [
            Type.from_dict(self.client, possible_type)
            for possible_type in possible_types
        ]
        self.add_possible_types(to_add)

    def __init_input_fields__(self, schema_data: dict[str, any]) -> None:
        input_fields = schema_data["inputFields"]
        if input_fields is None:
            return

        to_add = [InputValue.from_dict(self.client, field) for field in input_fields]
        self.add_input_fields(to_add)

    def add_fields(self, fields: list[Field]) -> None:
        if len(fields) == 0:
            return

        # Apply an appropriate kind if unset (typically for custom Types).
        if self.kind == Kind.UNSET:
            self.kind = Kind.OBJECT

        # Check for type exceptions
        for field in fields:
            field_id = f"{self.name}.field.{field.name}"

            # Resolve exceptions
            if field_id in self.type_exceptions:
                field.opencti_type = Type.from_name(
                    self.client, self.type_exceptions[field_id], custom=True
                )
                if Type.verbose:
                    print(
                        "INFO",
                        f"Replaced {field_id} with custom type {field.opencti_type}",
                    )

            self.fields[field.name] = field

    def add_enum_values(self, enum_values: list[str | EnumValue]) -> None:
        if len(enum_values) == 0:
            return

        # Apply an appropriate kind if unset (typically for custom Types).
        if self.kind == Kind.UNSET:
            self.kind = Kind.ENUM

        for value in enum_values:
            if isinstance(value, str):
                value = EnumValue(name=value)
            self.enum_values[value.name] = value

    def add_interfaces(self, interfaces: list["Type"]) -> None:
        if len(interfaces) == 0:
            return

        for interface in interfaces:
            self.interfaces[interface.name] = interface

    def add_possible_types(self, possible_types: list["Type"]) -> None:
        if len(possible_types) == 0:
            return

        # Apply an appropriate kind if unset (typically for custom Types).
        if self.kind == Kind.UNSET:
            self.kind = Kind.INTERFACE

        for possible_type in possible_types:
            self.possible_types[possible_type.name] = possible_type

    def add_input_fields(self, input_fields: list[InputValue]) -> None:
        if len(input_fields) == 0:
            return

        # Apply an appropriate kind if unset (typically for custom Types).
        if self.kind == Kind.UNSET:
            self.kind = Kind.INPUT_OBJECT

        # Check for type exceptions
        for field in input_fields:
            field_id = f"{self.name}.input_field.{field.name}"

            # Resolve exceptions
            if field_id in self.type_exceptions:
                field.opencti_type = Type.from_name(
                    self.client, self.type_exceptions[field_id], custom=True
                )
                if Type.verbose:
                    print(
                        "INFO",
                        f"Replaced {field_id} with custom type {field.opencti_type}",
                    )

            self.input_fields[field.name] = field

    def add_custom_fields(self, custom_fields: list[Field | InputValue]) -> None:
        if len(custom_fields) == 0:
            return

        # Apply an appropriate kind if unset (typically for custom Types).
        if self.kind == Kind.UNSET:
            self.kind = Kind.OBJECT

        for field in custom_fields:
            self.custom_fields[field.name] = field

    def print(self) -> None:
        txt = [f"Type: {self.graphql_type(full=True)} ({self.kind.value})\n"]
        if self.interfaces:
            txt += ["\nInterfaces:", self.interfaces]
        if self.fields:
            txt += ["\nFields:", self.fields]
        if self.input_fields:
            txt += ["\nInput Fields:", self.input_fields]
        if self.custom_fields:
            txt += ["\nCustom Fields:", self.custom_fields]
        if self.enum_values:
            txt += ["\nEnum Values:", self.enum_values]

        print(*txt)

    def resolve_imports(self) -> dict[str, bool]:
        """Returns all Golang imports required to print this Type."""
        go_deps: dict[str, bool] = {
            f'"github.com/weisshorn-cyd/gocti/{self.pkg}"': True
        }

        for f in self.printed_fields():
            go_deps[f'"github.com/weisshorn-cyd/gocti/{f.opencti_type.pkg}"'] = True
            if f.opencti_type.name == "DateTime":
                go_deps['"time"'] = True

        if self.kind in [Kind.INTERFACE, Kind.UNION]:
            go_deps['"reflect"'] = True

        if self.is_custom:
            go_deps['"github.com/goccy/go-json"'] = True

        return go_deps

    def resolve_type_dependencies(
        self, recursive: bool = False, processed: list[str] | None = None
    ) -> dict[str, "Type"]:
        """Returns a dict with all other Type dependencies required to print this Type."""

        output: dict[str, "Type"] = {}
        printed_fields = self.printed_fields()

        for field in printed_fields:
            if field.opencti_type.is_custom:
                continue
            output[field.opencti_type.name] = field.opencti_type
        for impl in self.possible_types.values():
            if impl.is_custom:
                continue
            output[impl.name] = impl

        if recursive:
            if processed is None:
                processed = []
            if self.name not in processed:
                processed.append(self.name)

                for field_type in output.values():
                    output = output | field_type.resolve_type_dependencies(
                        True, processed
                    )

        output = dict(sorted(output.items()))

        return output

    def printed_fields(self) -> list[Field | InputValue]:
        """Returns the list of all fields of this Type to be printed."""

        output: list[Field | InputValue] = []
        for l in self.printed_fields_by_interfaces().values():
            output += l

        return output

    def printed_fields_by_interfaces(self) -> dict[str, list[Field | InputValue]]:
        """Returns a dict with all fields of this Type to be printed for each Interface / Type source.
        It helps organise the Golang Type definition in sections."""

        fields = self.fields | self.input_fields
        fields_to_print: dict[str, list[Field | InputValue]] = {}
        impl_interfaces = dict(
            sorted(
                self.interfaces.items(),
                key=lambda x: self.interface_resolve_order.index(x[0])
                if x[0] in self.interface_resolve_order
                else 1000,
            )
        )

        # Attribute / reject fields based on interfaces
        for ki, vi in impl_interfaces.items():
            fields_to_print[ki] = []
            for kf, vf in vi.fields.items():
                if kf in fields:
                    fields.pop(kf)
                    if (
                        ki not in self.ignored_fields_list
                        or kf not in self.ignored_fields_list[ki]
                    ):
                        fields_to_print[ki].append(vf)
            if len(fields_to_print[ki]) == 0:
                fields_to_print.pop(ki)

        # Attribute remaining fields
        fields_to_print[self.name] = [
            f
            for f in fields.values()
            if self.name not in self.ignored_fields_list
            or f.name not in self.ignored_fields_list[self.name]
        ]

        # Add custom fields
        if len(self.custom_fields) > 0:
            fields_to_print["Custom"] = list(self.custom_fields.values())

        return fields_to_print

    def graphql_type(self, full: bool = False) -> str:
        """Returns the Type's GraphQL string representation.
        The required '!' and list '[]' identifiers are added as needed if 'full' is True.
        """
        representation = self.name

        if not full:
            return representation
        if self.is_list:
            representation = (
                f"[{representation + ('!' if self.is_content_required else '')}]"
            )
        if self.is_required:
            representation += "!"
        return representation

    def go_type(self, target_pkg: str, full: bool = False) -> str:
        """Returns the Type's Golang representation.
        The list '[]' identifier prefix is added as needed if 'full' is True.
        Standard computation can be overridden using 'Type.scalar_map'"""
        if self.name in Type.scalar_map:
            go = Type.scalar_map[self.name]
        else:
            go = utils.go_name(self.name)

        if self.kind != Kind.SCALAR and target_pkg != self.pkg:
            go = self.pkg + "." + go

        if self.is_list and full:
            go = "[]" + go

        return go

    def go_write(
        self,
        tags: list[str] = None,
        with_null_marshal: bool = False,
        add_dependencies: bool = False,
    ) -> str:
        """Returns the formatted text required to declare the Type in Golang."""

        go_type_declaration = ""
        if tags is None:
            tags = []

        # Declaration depends on the kind
        match self.kind:
            # Do not write scalars
            case Kind.SCALAR:
                return ""

            # Enums are written as constants of an aliased string type
            case Kind.ENUM:
                go_type_declaration = self._go_write_enum()

            # Input Objects and Objects are written as structs.
            # Interfaces and Unions are structs as well, but have an extra field and 2 methods
            # attached to implement the GraphQLInterface golang interface.
            case Kind.INPUT_OBJECT | Kind.OBJECT | Kind.INTERFACE | Kind.UNION:
                go_type_declaration = self._go_write_type(tags=tags)

            # Unimplemented error
            case _:
                raise ValueError(f"go_write() undefined for {self.name} ({self.kind})")

        # Add the field and special methods for interfaces / unions.
        if self.kind in [Kind.INTERFACE, Kind.UNION] and "gocti" in tags:
            go_type_declaration += "\n\n" + Template(INTERFACE_IMPL_GO_TMPL).substitute(
                t=self.go_type(self.pkg)[0].lower(),
                TypeName=self.go_type(self.pkg),
                reflectImplTypes="\n\t\t".join(
                    [
                        f"reflect.TypeFor[{impl.go_type(target_pkg=self.pkg)}](),"
                        for impl in self.possible_types.values()
                    ]
                ),
            )

        # Add null marshalling
        if with_null_marshal:
            go_type_declaration += "\n\n" + Template(NULL_MARSHALL_GO_TMPL).substitute(
                t=self.go_type(self.pkg)[0].lower(),
                TypeName=self.go_type(self.pkg),
            )

        # Add dependencies
        if add_dependencies:
            declarations = [go_type_declaration]
            declarations += [
                dependency.go_write(
                    tags=tags,
                    with_null_marshal=with_null_marshal,
                    add_dependencies=False,
                )
                for dependency in self.resolve_type_dependencies(
                    recursive=True
                ).values()
            ]
            declarations = list(filter(None, declarations))
            go_type_declaration = "\n\n".join(declarations)

        return utils.norm_string(go_type_declaration)

    def _go_write_enum(self):
        """ENUMs are written as constants of an aliased string type."""

        enum_value_template = Template('${TypeName}$Suffix $TypeName = "$value"')
        enum_values_declarations = [
            enum_value_template.substitute(
                TypeName=self.go_type(self.pkg),
                Suffix=utils.go_name(enum_value.name),
                value=enum_value.name,
            )
            for enum_value in self.enum_values.values()
        ]

        return Template(ENUM_GO_TMPL).substitute(
            TypeName=self.go_type(self.pkg),
            EnumValues="\n\t".join(enum_values_declarations),
        )

    def _go_write_type(self, tags: list[str] = None):
        """For INPUT_OBJECTs, OBJECTs, INTERFACEs, and UNIONs are written as structs.
        INTERFACEs and UNIONs have an additionnal field to store extra field data of GraphQL Type implementations."""
        if tags is None:
            tags = []

        def tags_tmpl(tags: list[str], value: str) -> str:
            return (
                (
                    " `"
                    + " ".join(
                        [
                            f'{tag}:"{value},omitempty"'
                            if tag == "json" and value != "-"
                            else f'{tag}:"{value}"'
                            for tag in tags
                        ]
                    )
                    + "`"
                )
                if len(tags) > 0
                else ""
            )

        struct_field_template = Template(f"$GoName $GoType{tags_tmpl(tags, '$name')}")

        # Struct fields are separated in groups belonging to the implemented interfaces.
        printable_fields = self.printed_fields_by_interfaces()

        fields = []
        for interface in printable_fields:
            fields += ["", f"// {interface}"] + [
                struct_field_template.substitute(
                    GoName=utils.go_name(field.name),
                    GoType=("*" if field.opencti_type.name == self.name else "")
                    + field.opencti_type.go_type(self.pkg, full=True),
                    name=field.name if interface != "Custom" else "-",
                )
                for field in printable_fields[interface]
            ]

        # Add a 'Remain' field to interfaces
        if self.kind in [Kind.INTERFACE, Kind.UNION] and "gocti" in tags:
            fields += [
                "",
                "// For storing interface implementations' data",
                'Remain map[string]any `gocti:",remain"`',
            ]

        output = Template(STRUCT_GO_TMPL).substitute(
            TypeName=self.go_type(self.pkg),
            Fields="\n\t".join(fields),
        )

        return (
            self._go_write_docstring() + "\n" + output if not self.is_custom else output
        )

    def _go_write_docstring(self) -> str:
        if self.kind in [Kind.OBJECT, Kind.INPUT_OBJECT]:
            return "\n".join(
                [
                    f"// {self.name} represents a GraphQL {self.kind.name}",
                    "// Some fields from the OpenCTI schema may be missing",
                    "// (See the examples for ways to expand an existing type)",
                ]
            )
        elif self.kind in [Kind.INTERFACE, Kind.UNION]:
            return "\n".join(
                [
                    f"// {self.name} represents a GraphQL {self.kind.name}",
                    "// Some fields from the OpenCTI schema may be missing",
                    "// (See the examples for ways to expand an existing type)",
                    "//",
                    "// Available implementations are:",
                ]
                + [f"// - [{t}]" for t in self.possible_types]
                + [
                    "//",
                    "// (See the examples for ways to decode an interface field into an implementation)",
                ]
            )
        else:
            return ""
