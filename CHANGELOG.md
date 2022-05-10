# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- [(GH-163)](https://github.com/puppetlabs/pct/issues/341) `--codedir {string}` flag is omitted, use working directory

## [0.2.0]

### Added

- [(GH-341)](https://github.com/puppetlabs/pct/issues/341) `prm build` command to package tools for distribution
- [(GH-130)](https://github.com/puppetlabs/prm/issues/130) Add support for multi-tool validation
- [(GH-129)](https://github.com/puppetlabs/prm/issues/129) Tool output handling to either `terminal` or `file`
- [(GH-125)](https://github.com/puppetlabs/prm/issues/125) Autocompletion for validation tools
- [(GH-127)](https://github.com/puppetlabs/prm/issues/127) `prm validate` command minimal implementation on Docker backend

### Fixes

- [(GH-139)](https://github.com/puppetlabs/prm/issues/139) Fix invalid default Puppet agent Docker image version
- [(GH-138)](https://github.com/puppetlabs/prm/issues/138) Ensures the `cache dir` is created if using default value, on `exec` and `validate` commands
- [(GH-137)](https://github.com/puppetlabs/prm/issues/137) List tools in alphabetical order, instead of random order

## [0.1.0]

### Added

- [(GH-11)](https://github.com/puppetlabs/prm/issues/11) `prm explain` command to surface documentation in the terminal
- [(GH-13)](https://github.com/puppetlabs/prm/issues/13) `prm install` command to add new tools or versions of tools
- [(GH-14)](https://github.com/puppetlabs/prm/issues/14) `prm status` command to verify status of the backend
- [(GH-15)](https://github.com/puppetlabs/prm/issues/15) `prm get|set backend` commands to enable setting the backend from the terminal
- [(GH-16)](https://github.com/puppetlabs/prm/issues/16) `prm get|set puppet` commands to enable setting the Puppet runtime version from the terminal
- [(GH-35)](https://github.com/puppetlabs/prm/issues/35) `prm exec` command to run tools in a Puppet Runtime

[Unreleased]: https://github.com/puppetlabs/prm/compare/0.2.0..main
[0.2.0]: https://github.com/puppetlabs/prm/releases/tag/0.2.0
[0.1.0]: https://github.com/puppetlabs/prm/releases/tag/0.1.0
