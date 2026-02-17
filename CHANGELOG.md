# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.64.0] - 2026-02-17

### Changed
- Support OpenCTI version 6.9.18

## [0.63.0] - 2026-02-16

### Changed
- Support OpenCTI version 6.9.17

## [0.62.0] - 2026-02-04

### Changed
- Support OpenCTI version 6.9.15

## [0.61.0] - 2026-01-29

### Changed
- Support OpenCTI version 6.9.8
- Healthcheck for elastichsearch in docker compose

### Fixed
- Missing type `Any` in mapping of generator
- Wrong resolve order in code generation

## [0.60.0] - 2026-01-14

### Changed
- Support OpenCTI version 6.9.7

## [0.59.0] - 2026-01-12

### Changed
- Support OpenCTI version 6.9.6

## [0.58.0] - 2025-12-17

### Changed
- Support OpenCTI version 6.9.1 - 6.9.4

## [0.57.0] - 2025-12-16

### Changed
- Support OpenCTI version 6.9.0

## [0.56.0] - 2025-12-12

### Changed
- Support OpenCTI version 6.8.17

## [0.55.0] - 2025-12-10

### Changed
- Support OpenCTI version 6.8.16

## [0.54.0] - 2025-12-03

### Changed
- Support OpenCTI version 6.8.15

## [0.53.0] - 2025-11-27

### Changed
- Support OpenCTI version 6.8.14

## [0.52.0] - 2025-11-24

### Changed
- Support OpenCTI version 6.8.13

### Fixed
- Add missing status scope for `set_status_in_workflow` query.

## [0.51.0] - 2025-11-19

### Changed
- Support OpenCTI version 6.8.12

## [0.50.0] - 2025-11-17

### Changed
- Support OpenCTI version 6.8.11

## [0.49.0] - 2025-11-03

### Changed
- Support OpenCTI version 6.8.10

## [0.48.0] - 2025-10-30

### Changed
- Support OpenCTI version 6.8.9

## [0.47.0] - 2025-10-27

### Changed
- Support OpenCTI version 6.8.7 - 6.8.8
- Bump Go to version 1.25.1
- Bump golangci-lint to 2.5.0

## [0.46.0] - 2025-10-03

### Changed
- Support OpenCTI version 6.8.2 - 6.8.6

## [0.45.0] - 2025-10-01

### Changed
- Support OpenCTI version 6.8.1

## [0.44.0] - 2025-09-30

### Changed
- Support OpenCTI version 6.8.0

## [0.43.0] - 2025-09-11

### Changed
- Support OpenCTI version 6.7.19 - 6.7.20

## [0.42.0] - 2025-09-10

### Changed
- Support OpenCTI version 6.7.18

## [0.41.0] - 2025-09-05

### Changed
- Support OpenCTI version 6.7.17

## [0.40.0] - 2025-09-03

### Changed
- Support OpenCTI version 6.7.16

## [0.39.0] - 2025-08-29

### Changed
- Support OpenCTI version 6.7.15

## [0.38.0] - 2025-08-27

### Changed
- Support OpenCTI version 6.7.14

## [0.37.0] - 2025-08-20

### Changed
- Support OpenCTI version 6.7.12

## [0.36.0] - 2025-08-13

### Changed
- Support OpenCTI version 6.7.10 - 6.7.11

## [0.35.0] - 2025-07-31

### Changed
- Support OpenCTI version 6.7.8 - 6.7.9

## [0.34.0] - 2025-07-23

### Changed
- Support OpenCTI version 6.7.6 - 6.7.7

## [0.33.0] - 2025-07-18

### Changed
- Support OpenCTI version 6.7.5

## [0.32.0] - 2025-07-10

### Changed
- Support OpenCTI version 6.7.3 - 6.7.4

### Fixed
- Auto update of CHANGELOG when having a single OpenCTI version

## [0.31.0] - 2025-07-01

### Changed
- Support OpenCTI version 6.7.1 - 6.7.2
- Auto-update CI does not error on Go linting failure
- Disable `funcorder` linter

## [0.30.0] - 2025-06-25

### Changed
- Support OpenCTI version 6.7.0
- Bump Go to version 1.24.4
- Update golangci-lint to v2
- Update range version in CHANGELOG when creating a new release

### Fixed
- Unittests for OpenCTI 6.7.0 new fields
- Example in README

## [0.29.0] - 2025-05-23

### Changed
- Support OpenCTI version 6.6.12 - 6.6.18

### Fixed
- Validate PyCTI version is last of supported OpenCTI version range
- Only check for graphql changes if there is a new version of OpenCTI
- Update README with new version when there are graphql changes

## [0.28.0] - 2025-05-12

### Changed
- Support OpenCTI version 6.6.10
- Only create a new release if there are some graphql changes

## [0.27.0] - 2025-04-29

### Changed
- Support OpenCTI version 6.6.8

## [0.26.0] - 2025-04-25

### Changed
- Support OpenCTI version 6.6.7

## [0.25.0] - 2025-04-23

### Changed
- Support OpenCTI version 6.6.6

## [0.24.0] - 2025-04-18

### Changed
- Support OpenCTI version 6.6.5

## [0.23.0] - 2025-04-16

### Changed
- Support OpenCTI version 6.6.4

## [0.22.0] - 2025-04-13

### Changed
- Support OpenCTI version 6.6.3

## [0.21.0] - 2025-04-09

### Changed
- Support OpenCTI version 6.6.1

## [0.20.0] - 2025-04-08

### Changed
- Support OpenCTI version 6.5.11

## [0.19.0] - 2025-03-31

### Changed
- Support OpenCTI version 6.5.10

## [0.18.0] - 2025-03-24

### Changed
- Support OpenCTI version 6.5.9
- Pin actions version to hash

### Fixed
- Wrong action hash

## [0.17.0] - 2025-03-19

### Changed
- Support OpenCTI version 6.5.8

## [0.16.0] - 2025-03-18

### Changed
- Support OpenCTI version 6.5.7

## [0.15.0] - 2025-03-12

### Changed
- Support OpenCTI version 6.5.6

## [0.14.0] - 2025-03-10

### Changed
- Support OpenCTI version 6.5.5

## [0.13.0] - 2025-03-06

### Changed
- Support OpenCTI version 6.5.4

## [0.12.0] - 2025-02-25

### Changed
- Support OpenCTI version 6.5.3

## [0.11.0] - 2025-02-18

### Changed
- Support OpenCTI version 6.5.2

## [0.10.0] - 2025-02-17

### Changed
- Support OpenCTI version 6.5.1

## [0.9.0] - 2025-02-11

### Changed
- Support OpenCTI version 6.5.0

### Fixed
- Auto-update workflow uses a branch name unique to the OpenCTI version it
  updates for

## [0.8.0] - 2025-01-30

### Changed
- Support OpenCTI version 6.4.10

## [0.7.0] - 2025-01-28

### Changed
- Support OpenCTI version 6.4.9

### Fixed
- Auto-update workflow correctly no longer opens a pull request if OpenCTI is
  already at the latest version

## [0.6.0] - 2025-01-24

### Added
- CI actions workflow to auto-update for new OpenCTI versions

### Changed
- Support OpenCTI version 6.4.8
- Linting exception for long struct tags
- Regenerate GoCTI
- Run workflows only when there are relevant changes
- Bump Go to version 1.23.5

### Fixed
- The `generate` Makefile target now correctly formats the generated code
- Specify golangci-lint config file in workflow

## [0.5.0] - 2024-01-17

### Changed
- Support OpenCTI version 6.4.7

## [0.4.0] - 2024-01-13

### Added
- CI actions validate the generator version

### Fixed
- Accept a max confidence level of 0
- Format Python code according to 2025 style guide

### Changed
- Support OpenCTI version 6.4.6

## [0.3.0] - 2025-01-07

### Added
- Query to unassign a group from a user
- Query to unassign a role from a group
- Query to unassign a marking definition from a group
- Query to unassign a capability from a role
- Query to unset a status from a workflow

### Changed
- Support OpenCTI version 6.4.5

## [0.2.0] - 2024-12-18

### Changed
- Support OpenCTI version 6.4.4

## [0.1.0] - 2024-12-05

### Added
- Initial version
- Support OpenCTI version 6.3.13
