# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
