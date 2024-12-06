def search_dict(data: dict, key: str) -> list[any]:
    """Recursively searches a dictionary for a given key and returns all values
    that match the key, from lowest to deepest depth if on a single leaf."""
    out = []
    for k, v in data.items():
        if k == key and v is not None:
            out.append(v)
        elif isinstance(v, dict):
            out += search_dict(data=v, key=key)
    return out


def go_name(name: str) -> str:  # noqa: F901
    """Converts any name (typically GraphQL types and fields) to a capitalized
    version following the Golang conventions."""
    if is_go_name(name):
        return name

    words = []
    for s in name.split("_"):
        if len(s) == 0:
            continue
        for a in s.split("-"):
            for n in split_caps(a):
                if n.lower().endswith("url"):
                    n = n[:-3] + "URL"
                if n.lower().endswith("id") and not n.lower().endswith("valid"):
                    n = n[:-2] + "ID"
                if n.lower().endswith("ids"):
                    n = n[:-3] + "IDs"
                if n.lower().endswith("http"):
                    n = n[:-4] + "HTTP"
                if n.lower().endswith("vm"):
                    n = n[:-2] + "VM"
                words.append(n[0].capitalize() + n[1:])
    name = "".join(words)

    return name.replace("?", "").replace("!", "").replace("[", "").replace("]", "")


def is_go_name(name: str) -> bool:
    """Checks whether a string is a Go-style name."""
    if not name[0].isupper():
        return False
    if "_" in name or "-" in name:
        return False
    if "url" in name.lower() and "URL" not in name:
        return False
    if "id" in name.lower() and "valid" not in name.lower() and "ID" not in name:
        return False
    if "http" in name.lower() and "HTTP" not in name:
        return False
    if "vm" in name.lower() and "VM" not in name:
        return False
    if "!" in name or "[" in name or "?" in name:
        return False

    return True


def split_caps(txt: str) -> list[str]:
    """Splits a single word into its sub-words. Keeps blocks of capital letters
    together, e.g. 'stixID' -> ['stix', 'ID']"""
    out = [txt[0]]
    for c in txt[1:]:
        if c.isupper():
            if out[-1][-1].isupper():
                out[-1] += c
            else:
                out.append(c)
        else:
            if out[-1].isupper() and len(out[-1]) > 1:
                out.append(c)
            else:
                out[-1] += c
    return out


def norm_string(txt: str) -> str:
    """Replaces tabs with 4 spaces."""
    return txt.replace("\t", "    ")


def format_go_imports(import_groups: list[list[str]]) -> str:
    """Formats groups of Go imports."""
    imports = len(import_groups) * [""]
    for i, group in enumerate(import_groups):
        imports[i] = (
            "\t" + "\n\t".join(go_import.strip() for go_import in group).strip()
        )
    return "\t" + "\n\n".join(group for group in imports).strip()


def format_graphql_attributes(attr: str, offset: int = 0) -> str:
    attr = attr.replace("\t", "")
    attr = attr.replace("  ", "")

    output = ""

    for line in attr.split("\n"):
        line = line.strip()

        if len(line) == 0:
            continue
        elif line[-1] == "{":
            output += "\t" * offset + line + "\n"
            offset += 1
        elif line[-1] == "}":
            offset -= 1
            output += "\t" * offset + line + "\n"
        else:
            output += "\t" * offset + line + "\n"

    return output.strip("\n")
