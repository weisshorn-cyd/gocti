# GoCTI Type Generator

This package holds helpers for generating Entity files and required types for use within GoCTI.
It makes use of the official OpenCTI Python Client ([pycti](https://github.com/OpenCTI-Platform/client-python)) to retrieve information from the GraphQL schema of a running OpenCTI instance.

### Usage

To generate the GoCTI files, run the [main.py](gocti_type_generator/main.py) file. You must provide one of the following "cmd" values:

- "entities": Will generate the GoCTI Entity files and dependent types.
- "helpers": Will scrape the current GoCTI files and generate the main helper files:
  - [read.go](../../read.go)
  - [list.go](../../list.go)
  - [create.go](../../create.go)
  - [delete.go](../../delete.go)

It also accepts the following additional arguments:

| Parameter      | Default Value                       | Description                                       |
| -------------- | ----------------------------------- | ------------------------------------------------- |
| cmd            | n/a                                 | Required argument. See above.                     |
| -e, --entities | n/a                                 | Scrapes the GoCTI repo for the names of all the entities to generate.        |
| -p, --pkg      | "entity"                            | The target gocti package for the generated files. |
| --env          | n/a                                 | Path to the `.env` config file.                   |

### Configuration

The target GoCTI repo should be provided.
Additionally, the generator also requires an OpenCTI url and token to retrieve schema information.
Those parameters can be passed either via command line arguments or environment variables.
A sample [.env](./.env.sample) file is provided containing the proper variable declarations.

| Cmd Line Argument | Environment Variable | Default Value         | Description                                                  |
| ----------------- | -------------------- | --------------------- | ------------------------------------------------------------ |
| -u, --url         | OPENCTI_URL          | http://localhost:8080 | Url of a running OpenCTI instance.                           |
| -t, --token       | OPENCTI_TOKEN        | n/a                   | Token used to retrieve the GraphQL schema.                   |
| -r, --repo        | GOCTI_REPO           | .                     | Path to the GoCTI repo in which the files will be generated. |
