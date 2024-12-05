from dataclasses import dataclass
from pathlib import Path
from string import Template

from pycti import OpenCTIApiClient

from gocti_type_generator import utils
from gocti_type_generator.type import (
    GRAPHQL_TYPES_PKG_NAME,
    MAX_NESTED_TYPE_DEPTH,
    InputValue,
    Kind,
    Type,
)


@dataclass
class GraphQLQueryArg:
    name: str
    type_name: str  # GraphQL type name
    alias: str = ""
    resolved_type: Type | None = None


@dataclass
class GraphQLQuery:
    """Represents any GraphQL query (both queries and mutations) than can be made on the server.
    The returned Type and the argument Types are not resolved at creation to avoid unnecessary computations.
    """

    name: str
    type_name: str  # GraphQL type name
    args: dict[str, GraphQLQueryArg]  # Argument name: GraphQL type name
    kind: str  # "query" or "mutation"
    content: str = ""
    deprecated: bool = False
    resolved_type: Type | None = None

    def must_embed(self, min_depth_embed: int) -> bool:
        return self.write().count("\n") > min_depth_embed

    def write(self, offset: int = 0) -> str:
        """Returns the properly formatted GraphQL query string, with all arguments set as variables."""

        content = (" " + self.content) if self.content != "" else ""
        query_args = dict(sorted(self.args.items()))

        query_template = Template(
            self.kind + " ($vars) {\n\t$name ($args)" + content + "\n}"
        )
        vars = ("\n" + " " * (len(self.kind) + 2)).join(
            [f"${k}: {v.type_name}" for k, v in query_args.items()]
        )
        args = ("\n\t" + " " * (len(self.name) + 2)).join(
            [f"{k}: ${k}" for k in query_args]
        )

        query = query_template.substitute(
            name=self.name,
            vars=vars,
            args=args,
        )

        # Add offset
        offset: str = "\t" * offset
        query = offset + query.strip().replace("\n", "\n" + offset)

        return utils.norm_string(query)


def list_graphql_queries(client: OpenCTIApiClient) -> dict[str, GraphQLQuery]:
    """Retrieves a list of all available queries and mutations from the server."""

    output: dict[str, GraphQLQuery] = {}

    # Nested type query.
    temp = "name, kind, ofType { $type_name_graphql }"
    type_name_graphql = "name, kind"
    for _ in range(MAX_NESTED_TYPE_DEPTH):
        type_name_graphql = Template(temp).substitute(
            type_name_graphql=type_name_graphql,
        )

    # Same operation for both queries and mutations.
    for kind in ["query", "mutation"]:
        # Prepare the query.
        graphql_query = Template(
            """query { __schema { ${kind}Type { fields {
                name, isDeprecated, type { $type_name_graphql }, args { name, type{ $type_name_graphql } }
            }}}}"""
        ).substitute(kind=kind, type_name_graphql=type_name_graphql)

        queries = client.query(graphql_query)["data"]["__schema"][kind + "Type"][
            "fields"
        ]

        # Process data.
        for query in queries:
            name = query["name"]
            deprecated = query["isDeprecated"]
            return_type = Type.name_from_dict(query["type"])
            args: dict[str, GraphQLQueryArg] = {}
            for arg in query["args"]:
                arg_name = arg["name"]
                arg_type = Type.name_from_dict(arg["type"])
                args[arg_name] = GraphQLQueryArg(name=arg_name, type_name=arg_type)
            output[name] = GraphQLQuery(
                name=name,
                kind=kind,
                type_name=return_type,
                args=args,
                deprecated=deprecated,
            )

    for o in output.values():
        if o.deprecated:
            print(f"WARNING: {o.name} is deprecated!")

    return output


class UnknownQueryException(Exception):
    pass


DEFAULT_PKG_NAME = "entity"

ENTITY_FILE_GO_TMPL = """
package $pkg

import (
$imports
)

$entity_declaration

$additional_sections
""".strip()

READ_SECTION_GO_TMPL = """
// Implementing the [api.ReadableEntity] interface.
$embedded_query_segment
func ($e $Entity) ReadQueryString(customAttributes string) string {
	return fmt.Sprintf(
		$graphql_query,
		customAttributes,
	)
}

func ($e $Entity) ReadResponseField() string { return "$query_name" }
""".strip()
LIST_SECTION_GO_TMPL = """
// Implementing the [api.ListableEntity] interface.
$embedded_query_segment
func ($e $Entity) ListQueryString(customAttributes string) string {
	return fmt.Sprintf(
		$graphql_query,
		customAttributes,
	)
}

func ($e $Entity) ListResponseField() string { return "$query_name" }
""".strip()
CREATE_SECTION_GO_TMPL = """
// Implementing the [api.CreatableEntity] interface.
$embedded_query_segment
func ($e $Entity) CreateQueryString(customAttributes string) string {
	return fmt.Sprintf(
		$graphql_query,
		customAttributes,
	)
}

func ($e $Entity) CreateResponseField() string { return "$query_name" }

$input_type_segment
""".strip()
DELETE_SECTION_GO_TMPL = """
// Implementing the [api.DeletableEntity] interface.
$embedded_query_segment
func ($e $Entity) DeleteQueryString() string {
	return $graphql_query
}

func ($e $Entity) DeleteResponseField() string { return "$query_name" }
""".strip()

READ_HELPER_GO_TMPL = """
func (client *OpenCTIAPIClient) Read$Entity(
	ctx context.Context,
	customAttributes, id string,
) ($pkg.$Entity, error) {
	return api.StructuredRead[$pkg.$Entity, $pkg.$Entity](
        ctx, client, customAttributes, id,
    )
}""".strip()
LIST_HELPER_GO_TMPL = (
    """
// Available [list.Option] parameters:$options
func (client *OpenCTIAPIClient) List$Entities(
	ctx context.Context,
	customAttributes string,
	getAll bool,
	pageInfo *"""
    + GRAPHQL_TYPES_PKG_NAME
    + """.PageInfo,
	opts ...list.Option,
) ([]$pkg.$Entity, error) {
	return api.StructuredList[$pkg.$Entity, $pkg.$Entity](
        ctx, client, customAttributes, getAll, pageInfo, opts...,
    )
}""".strip()
)
CREATE_HELPER_GO_TMPL = """
func (client *OpenCTIAPIClient) Create$Entity(
	ctx context.Context,
	customAttributes string,
	input $pkg.$EntityAddInput,
) ($pkg.$Entity, error) {
	return api.StructuredCreate[$pkg.$Entity, $pkg.$Entity](
        ctx, client, customAttributes, input,
    )
}""".strip()
DELETE_HELPER_GO_TMPL = """
func (client *OpenCTIAPIClient) Delete$Entity(
	ctx context.Context,
	id string,
) (string, error) {
	return api.Delete[$pkg.$Entity](ctx, client, id)
}""".strip()

INPUT_TYPE_SEGMENT_STANDARD_GO_TMPL = """
$type_def

func (input $input_type_name) Input() (map[string]any, error) {
	return map[string]any{
		$func_fields
	}, nil
}""".strip()
INPUT_TYPE_SEGMENT_CUSTOM_GO_TMPL = """
$type_def

func (input $input_type_name) Input() (map[string]any, error) {
	output := map[string]any{}

	inputByte, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("cannot encode input: %w", err)
	}

	err = json.Unmarshal(inputByte, &output)
	if err != nil {
		return nil, fmt.Errorf("cannot decode input: %w", err)
	}

	return output, nil
}""".strip()

READ_QUERY_CONTENT = "{%s}"
LIST_QUERY_CONTENT = """
    {
        edges {
            node {%s}
        }
        pageInfo {
            startCursor
            endCursor
            hasNextPage
            hasPreviousPage
            globalCount
        }
    }
""".strip()
CREATE_QUERY_CONTENT = "{%s}"
DELETE_VIA_EDIT_QUERY_CONTENT = """
    {
        delete
    }
""".strip()
DELETE_QUERY_CONTENT = ""


class Entity:
    """Represents an OpenCTI entity, with its fields and related queries."""

    # Paths relative to the entity's folder
    default_properties_folder: str = "default_properties"
    list_queries_folder: str = "list_queries"
    read_queries_folder: str = "read_queries"
    create_queries_folder: str = "create_queries"
    delete_queries_folder: str = "delete_queries"

    # Overrides on standard query names
    query_name_exceptions: dict[str, str] = {
        "stixNestedRefRelationship.read": "stixRefRelationship",
        "stixNestedRefRelationship.create": "stixRefRelationshipAdd",
        "stixNestedRefRelationship.edit": "stixRefRelationshipEdit",
        "threatActorGroup.list": "threatActorsGroup",
        "threatActorIndividual.list": "threatActorsIndividuals",
    }

    # Overrides on standard entity default properties
    custom_default_properties: dict[str, list[str | list[str]]] = {
        "capability": ["id", "entity_type"],
        "group": [
            "id",
            "entity_type",
            "name",
            "default_assignation",
            "auto_new_marking",
            "group_confidence_level {\n    max_confidence\n}",
            "created_at",
            "updated_at",
        ],
        "role": ["id", "entity_type", "name", "created_at", "updated_at"],
        "user": ["id", "entity_type", "name", "user_email", "api_token"],
        "statusTemplate": ["id", "name", "color", "usages"],
        "caseTemplate": [
            "id",
            "name",
            "description",
            "tasks",
            ["edges", ["node", ["id", "name"]], "pageInfo", ["globalCount"]],
        ],
        "taskTemplate": ["id", "standard_id", "name", "description"],
        "subType": [
            "id",
            "label",
            "workflowEnabled",
            "statuses",
            ["id", "order", "template", ["name", "color", "id"]],
        ],
    }

    # List of queries expected to be unavailable and should not raise an error.
    query_name_expected_unavailable: list[str] = [
        "stix.type",
        "stix.list",
        "stix.create",
        "stixCoreObject.create",
        "stixNestedRefRelationship.type",
        "stixObjectOrStixRelationship.list",
        "stixObjectOrStixRelationship.create",
        "stixObjectOrStixRelationship.delete",
        "threatActor.create",
        "threatActor.delete",
        "capability.read",
        "capability.create",
        "capability.delete",
        "subType.create",
    ]

    def __init__(
        self,
        client: OpenCTIApiClient,
        graphql_queries: dict[str, GraphQLQuery],
        name: str,
        pkg: str = DEFAULT_PKG_NAME,
        min_depth_embed: int = 6,
    ):
        # For all server interactions
        self.client: OpenCTIApiClient = client
        self.graphql_queries: dict[str, GraphQLQuery] = graphql_queries

        # Quick name references
        self.Name: str  # CamelCased name
        self.name: str  # camelCased name
        self.Names: str  # CamelCased plural name
        self.names: str  # camelCased plural name
        self.na_me: str  # snake_cased name
        self.n: str  # first letter
        self.__init_names__(name=name)

        self.pkg: str = pkg
        self.filename: str = f"opencti_{self.na_me}.go"

        self.opencti_type: Type | None = (
            Type.from_name(self.client, self.Name)
            if f"{self.name}.type" not in self.query_name_expected_unavailable
            else None
        )

        # Miscelaneous
        self.min_depth_embed: int = min_depth_embed

        # Resolve default properties
        self.has_default_properties: bool = False
        self.default_properties: str
        self.embedded_default_properties_filename: str
        self.embedded_default_properties_var: str
        self.__init_default_properties__()

        # Resolve 'READ' query data
        self.has_embedded_read: bool = False
        self.read_query: GraphQLQuery | None = None
        self.embedded_read_query_filename: str
        self.embedded_read_query_var: str
        self.__init_read_data__()

        # Resolve 'LIST' query data
        self.has_embedded_list: bool = False
        self.list_query: GraphQLQuery | None = None
        self.embedded_list_query_filename: str
        self.embedded_list_query_var: str
        self.__init_list_data__()

        # Resolve 'CREATE' query data
        self.has_embedded_create: bool = False
        self.create_query: GraphQLQuery | None = None
        self.embedded_create_query_filename: str
        self.embedded_create_query_var: str
        self.__init_create_data__()

        # Resolve Input Type
        self.input_type: Type | None = None
        self.__init_input_type__()

        # Resolve 'DELETE' query data
        self.has_embedded_delete: bool = False
        self.delete_use_edit: bool = False
        self.delete_query: GraphQLQuery | None = None
        self.embedded_delete_query_filename: str
        self.embedded_delete_query_var: str
        self.__init_delete_data__()

    def __str__(self) -> str:
        return self.__repr__()

    def __repr__(self) -> str:
        return self.Name + "(Entity)"

    def __init_names__(self, name: str):
        # Standard names
        self.Name = utils.go_name(name)
        self.name = self.Name[0].lower() + self.Name[1:]
        # Dealing with plurals
        self.Names = self.Name + "s"
        if self.Name[-1] == "y":
            self.Names = self.Name[:-1] + "ies"
        if self.Name[-2:] == "is":
            self.Names = self.Name[:-2] + "es"
        if self.Name == "CourseOfAction":
            self.Names = "CoursesOfAction"
        self.names = self.Names[0].lower() + self.Names[1:]
        # Miscellaneous names
        self.na_me = "_".join([a[0].lower() + a[1:] for a in utils.split_caps(name)])
        self.n = name[0].lower()

    def __init_default_properties__(self):
        # Retrieve properties from pycti
        self.default_properties = self._default_properties()
        if self.default_properties != "":
            self.has_default_properties = True
        else:
            # Should never be false at the moment, as the minimum properties is 'id'
            self.has_default_properties = False
        # Embedding parameters
        self.embedded_default_properties_filename = (
            self.na_me + "_default_properties.txt"
        )
        self.embedded_default_properties_var = self.name + "DefaultProperties"

    def __init_read_data__(self):
        # Setup variables
        read_query_name = self.query_name_exceptions.get(f"{self.name}.read", self.name)
        self.embedded_read_query_filename = self.na_me + "_read_query.txt"
        self.embedded_read_query_var = self.name + "ReadQueryString"

        if read_query_name in self.graphql_queries:
            self.read_query = self.graphql_queries[read_query_name]
            self.read_query.content = READ_QUERY_CONTENT
        elif f"{self.name}.read" in self.query_name_expected_unavailable:
            return
        else:
            raise UnknownQueryException(f"Unknown query: {read_query_name}")

        # Testing embedding necessity
        if self.read_query.must_embed(self.min_depth_embed):
            self.has_embedded_read = True

    def __init_list_data__(self):
        # Setup variables
        list_query_name = self.query_name_exceptions.get(
            f"{self.name}.list", self.names
        )
        self.embedded_list_query_filename = self.na_me + "_list_query.txt"
        self.embedded_list_query_var = self.name + "ListQueryString"

        if list_query_name in self.graphql_queries:
            self.list_query = self.graphql_queries[list_query_name]
            self.list_query.content = LIST_QUERY_CONTENT
        elif f"{self.name}.list" in self.query_name_expected_unavailable:
            return
        else:
            raise UnknownQueryException(f"Unknown query: {list_query_name}")

        # Testing embedding necessity
        if self.list_query.must_embed(self.min_depth_embed):
            self.has_embedded_list = True

    def __init_create_data__(self):
        # Setup variables
        create_query_name = self.query_name_exceptions.get(
            f"{self.name}.create", self.name + "Add"
        )
        self.embedded_create_query_filename = self.na_me + "_create_query.txt"
        self.embedded_create_query_var = self.name + "CreateQueryString"

        if create_query_name in self.graphql_queries:
            self.create_query = self.graphql_queries[create_query_name]
            self.create_query.content = CREATE_QUERY_CONTENT
        elif f"{self.name}.create" in self.query_name_expected_unavailable:
            return
        else:
            raise UnknownQueryException(f"Unknown mutation: {create_query_name}")

        # Testing embedding necessity
        if self.create_query.must_embed(self.min_depth_embed):
            self.has_embedded_create = True

    def __init_input_type__(self):
        if self.create_query is None:
            return

        input_type_name = self.Name + "AddInput"
        input_fields = self.create_query.args

        # Check if input format is regular (ex:label) or not (ex:stixCyberObservable)
        if "input" in input_fields:
            arg = input_fields["input"]
            if input_type_name not in arg.type_name:
                input_type_name = arg.type_name.replace("!", "")
            type_ = Type.from_name(self.client, arg.type_name)
            is_regular = True
        else:
            type_ = Type.from_name(self.client, input_type_name, custom=True)
            is_regular = False
        type_.pkg = self.pkg

        # Populate additional fields
        for arg in input_fields.values():
            if arg.name == "input":
                continue
            field = InputValue(
                name=arg.name, opencti_type=Type.from_name(self.client, arg.type_name)
            )
            if is_regular:
                type_.add_custom_fields([field])
            else:
                type_.add_input_fields([field])

        self.input_type = type_

    def __init_delete_data__(self):
        # Setup variables
        delete_query_name = self.query_name_exceptions.get(
            f"{self.name}.delete", self.name + "Delete"
        )
        self.embedded_delete_query_filename = self.na_me + "_delete_query.txt"
        self.embedded_delete_query_var = self.name + "DeleteQueryString"

        # Some entities don't have a 'DELETE' mutation, but use 'EDIT' instead.
        edit_query_name = self.query_name_exceptions.get(
            f"{self.name}.edit", self.name + "Edit"
        )

        # Check for 'EDIT' first, as it is preferred.
        if edit_query_name in self.graphql_queries:
            self.delete_query = self.graphql_queries[edit_query_name]
            self.delete_query.content = DELETE_VIA_EDIT_QUERY_CONTENT
            self.delete_use_edit = True
        elif delete_query_name in self.graphql_queries:
            self.delete_query = self.graphql_queries[delete_query_name]
            self.delete_query.content = DELETE_QUERY_CONTENT
            self.delete_use_edit = False
        elif f"{self.name}.delete" in self.query_name_expected_unavailable:
            return
        else:
            raise UnknownQueryException(
                f"Unknown mutations: {edit_query_name} / {delete_query_name}"
            )

        # Testing embedding necessity
        if self.delete_query.must_embed(self.min_depth_embed):
            self.has_embedded_delete = True

    def _default_properties(self) -> str:
        # Get properties from pycti, if available
        try:
            props: str = eval(f"self.client.{self.na_me}.properties_with_files")
        except AttributeError:
            try:
                props: str = eval(f"self.client.{self.na_me}.properties")
            except AttributeError:

                def resolve_properties(properties: str | list[str], depth: int = -1):
                    output = ""
                    if isinstance(properties, str):
                        output = "\t" * depth + properties
                    else:
                        output += " {" if depth >= 0 else ""

                        for p in properties:
                            if isinstance(p, str):
                                output += f"\n{resolve_properties(p,depth+1)}"
                            else:
                                output += f"{resolve_properties(p,depth+1)}"

                        output.replace("\n", "\n" + "\t" * (depth + 1))
                        output += "\n" + "\t" * depth + "}" if depth >= 0 else ""

                    return utils.norm_string(output)

                return (
                    "id"
                    if self.name not in self.custom_default_properties
                    else resolve_properties(
                        self.custom_default_properties[self.name]
                    ).strip()
                )

        return utils.norm_string(utils.format_graphql_attributes(props))

    def resolve_imports(self) -> dict[str, bool]:
        go_deps: dict[str, bool] = {}

        if self.opencti_type:
            go_deps[f'"github.com/weisshorn-cyd/gocti/{GRAPHQL_TYPES_PKG_NAME}"'] = True
        if self.input_type:
            go_deps = go_deps | self.input_type.resolve_imports()

        if self.read_query or self.list_query or self.create_query or self.delete_query:
            go_deps['"fmt"'] = True

        if (
            self.has_default_properties
            or self.has_embedded_read
            or self.has_embedded_list
            or self.has_embedded_create
            or self.has_embedded_delete
        ):
            go_deps['_ "embed"'] = True

        return go_deps

    def resolve_type_dependencies(self) -> dict[str, Type]:
        """Returns a dict with all Type dependencies required to print this Entity."""
        output: dict[str, Type] = {}

        if self.input_type is not None:
            output = output | self.input_type.resolve_type_dependencies(recursive=True)

        if self.opencti_type is not None:
            output[self.opencti_type.name] = self.opencti_type
            output = output | self.opencti_type.resolve_type_dependencies(
                recursive=True
            )

        return output

    def get_file_txt(self, write_folder: Path | None = None) -> str:
        main_declaration = self._entity_declaration()

        sections = "\n\n".join(
            [
                self._default_properties_section(write_folder=write_folder),
                self._listable_entity_section(write_folder=write_folder),
                self._readable_entity_section(write_folder=write_folder),
                self._creatable_entity_section(write_folder=write_folder),
                self._deletable_entity_section(write_folder=write_folder),
            ]
        )

        golang_deps = self.resolve_imports()
        golang_deps.pop(f'"github.com/weisshorn-cyd/gocti/{self.pkg}"', "")

        imports = [[imp for imp in golang_deps if imp != '_ "embed"']]
        if '_ "embed"' in golang_deps:
            imports += [['_ "embed"']]

        return utils.norm_string(
            Template(ENTITY_FILE_GO_TMPL).substitute(
                pkg=self.pkg,
                imports=utils.format_go_imports(imports),
                entity_declaration=main_declaration,
                additional_sections=sections,
            )
        )

    def _entity_declaration(self) -> str:
        return (
            "type " + self.Name + " struct {}"
            if self.opencti_type is None or self.opencti_type.kind in [Kind.SCALAR]
            else "type "
            + self.Name
            + " struct {\n\t"
            + self.opencti_type.go_type(self.pkg)
            + ' `gocti:",squash"`'
            + "\n}"
        )

    def _default_properties_section(self, write_folder: Path | None = None) -> str:
        # Write file for embedded default properties
        if write_folder is not None:
            with open(
                Path(
                    write_folder,
                    self.default_properties_folder,
                    self.embedded_default_properties_filename,
                ),
                "w",
            ) as f:
                f.write(self.default_properties + "\n")

        # Resolve embedded segment
        embedded_properties = self._embedded_variable_segment(
            folder=self.default_properties_folder,
            filename=self.embedded_default_properties_filename,
            variable=self.embedded_default_properties_var,
        )

        # Resolve Golang section
        return Template(
            embedded_properties
            + "\n\nfunc ($e $Entity) DefaultProperties() string {\n"
            + "\treturn $embedded_properties_var\n}"
        ).substitute(
            e=self.n,
            Entity=self.Name,
            embedded_properties_var=self.embedded_default_properties_var,
        )

    def _readable_entity_section(self, write_folder: Path | None = None) -> str:
        # Exit if unavailable
        if self.read_query is None:
            return ""

        must_embed = self.read_query.must_embed(self.min_depth_embed)

        # Write file for embedded read query
        if write_folder:
            filepath = Path(
                write_folder,
                self.read_queries_folder,
                self.embedded_read_query_filename,
            )

            # Add or remove file
            if must_embed:
                with open(filepath, "w") as f:
                    f.write(self.read_query.write() + "\n")
            elif filepath.exists():
                Path.unlink(filepath)

        # Resolve embedded segment
        embedded_read = (
            f"\n{self._embedded_variable_segment(
                folder=self.read_queries_folder,
                filename=self.embedded_read_query_filename,
                variable=self.embedded_read_query_var,
            )}\n"
            if must_embed
            else ""
        )

        # Resolve Golang section
        return Template(READ_SECTION_GO_TMPL).substitute(
            e=self.n,
            Entity=self.Name,
            query_name=self.read_query.name,
            embedded_query_segment=embedded_read,
            graphql_query=(
                self.embedded_read_query_var
                if must_embed
                else f"`{self.read_query.write(offset=2).strip()}`"
            ),
        )

    def _listable_entity_section(self, write_folder: Path | None = None) -> str:
        # Exit if unavailable
        if self.list_query is None:
            return ""

        must_embed = self.list_query.must_embed(self.min_depth_embed)

        # Write file for embedded list query
        if write_folder:
            filepath = Path(
                write_folder,
                self.list_queries_folder,
                self.embedded_list_query_filename,
            )

            # Add or remove file
            if must_embed:
                with open(filepath, "w") as f:
                    f.write(self.list_query.write() + "\n")
            elif filepath.exists():
                Path.unlink(filepath)

        # Resolve embedded segment
        embedded_list = (
            f"\n{self._embedded_variable_segment(
                folder=self.list_queries_folder,
                filename=self.embedded_list_query_filename,
                variable=self.embedded_list_query_var,
            )}\n"
            if must_embed
            else ""
        )

        # Resolve Golang section
        return Template(LIST_SECTION_GO_TMPL).substitute(
            e=self.n,
            Entity=self.Name,
            query_name=self.list_query.name,
            embedded_query_segment=embedded_list,
            graphql_query=(
                self.embedded_list_query_var
                if must_embed
                else f"`{self.list_query.write(offset=2).strip()}`"
            ),
        )

    def _creatable_entity_section(self, write_folder: Path | None = None) -> str:
        # Exit if unavailable
        if self.create_query is None:
            return ""

        must_embed = self.create_query.must_embed(self.min_depth_embed)

        # Write file for embedded create query
        if write_folder:
            filepath = Path(
                write_folder,
                self.create_queries_folder,
                self.embedded_create_query_filename,
            )

            # Add or remove file
            if must_embed:
                with open(filepath, "w") as f:
                    f.write(self.create_query.write() + "\n")
            elif filepath.exists():
                Path.unlink(filepath)

        # Resolve embedded segment
        embedded_create = (
            f"\n{self._embedded_variable_segment(
                folder=self.create_queries_folder,
                filename=self.embedded_create_query_filename,
                variable=self.embedded_create_query_var,
            )}\n"
            if must_embed
            else ""
        )

        # Resolve Golang section
        return Template(CREATE_SECTION_GO_TMPL).substitute(
            e=self.n,
            Entity=self.Name,
            query_name=self.create_query.name,
            embedded_query_segment=embedded_create,
            graphql_query=(
                self.embedded_create_query_var
                if must_embed
                else f"`{self.create_query.write(offset=2).strip()}`"
            ),
            input_type_segment=self._input_type_segment(),
        )

    def _deletable_entity_section(self, write_folder: Path | None = None) -> str:
        # Exit if unavailable
        if self.delete_query is None:
            return ""

        must_embed = self.delete_query.must_embed(self.min_depth_embed)

        # Write file for embedded delete query
        if write_folder:
            filepath = Path(
                write_folder,
                self.delete_queries_folder,
                self.embedded_delete_query_filename,
            )

            # Add or remove file
            if must_embed:
                with open(filepath, "w") as f:
                    f.write(self.delete_query.write() + "\n")
            elif filepath.exists():
                Path.unlink(filepath)

        # Resolve embedded segment
        embedded_delete = (
            f"\n{self._embedded_variable_segment(
                folder=self.delete_queries_folder,
                filename=self.embedded_delete_query_filename,
                variable=self.embedded_delete_query_var,
            )}\n"
            if must_embed
            else ""
        )

        # Resolve Golang section
        return Template(DELETE_SECTION_GO_TMPL).substitute(
            e=self.n,
            Entity=self.Name,
            query_name=self.delete_query.name,
            embedded_query_segment=embedded_delete,
            graphql_query=(
                self.embedded_delete_query_var
                if must_embed
                else f"`{self.delete_query.write(offset=3).strip()}`"
            ),
        )

    def _embedded_variable_segment(
        self, folder: str, filename: str, variable: str
    ) -> str:
        return Template(
            "//go:embed $folder/$filename\nvar $variable string"
        ).substitute(
            folder=folder,
            filename=filename,
            variable=variable,
        )

    def _input_type_segment(self) -> str:
        # Exit if unavailable
        if self.input_type is None:
            return ""

        type_def_str = self.input_type.go_write(
            ["gocti", "json"],
            with_null_marshal=False,
            add_dependencies=False,
        )

        # Dealing with custom input types
        if not self.input_type.is_custom:
            input_type_segment = INPUT_TYPE_SEGMENT_STANDARD_GO_TMPL
            func_fields_str = "\n\t\t".join(
                ['"input": input,']
                + [
                    # Additional input fields are stored as custom fields.
                    f'"{field.name}": input.{utils.go_name(field.name)},'
                    for field in self.input_type.custom_fields.values()
                ]
            )

        else:
            input_type_segment = INPUT_TYPE_SEGMENT_CUSTOM_GO_TMPL
            func_fields_str = "\n\t\t".join(
                [
                    f'"{field.name}": input.{utils.go_name(field.name)},'
                    for field in self.input_type.input_fields.values()
                ]
            )

        return Template(input_type_segment).substitute(
            type_def=type_def_str,
            input_type_name=self.input_type.name,
            func_fields=func_fields_str,
        )

    def read_helper(self) -> str:
        return Template(READ_HELPER_GO_TMPL).substitute(
            pkg=self.pkg,
            Entity=self.Name,
        )

    def list_helper(self) -> str:
        options = [
            f"//\t- [list.With{utils.go_name(arg.alias if arg.alias != "" else arg.name)}]"
            for arg in dict(sorted(self.list_query.args.items())).values()
        ]
        if len(options) > 0:
            options = [""] + options

        return Template(LIST_HELPER_GO_TMPL).substitute(
            pkg=self.pkg,
            Entities=self.Names,
            Entity=self.Name,
            options="\n".join(options),
        )

    def create_helper(self) -> str:
        return Template(CREATE_HELPER_GO_TMPL).substitute(
            pkg=self.pkg,
            Entity=self.Name,
            EntityAddInput=self.input_type.name,
        )

    def delete_helper(self) -> str:
        return Template(DELETE_HELPER_GO_TMPL).substitute(
            pkg=self.pkg,
            Entity=self.Name,
        )
