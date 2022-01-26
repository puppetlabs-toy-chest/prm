---
title: "Anatomy of a PRM Tool"
description: "How a PRM tool is structured"
category: concept
tags:
  - tools
  - development
weight: 20
---

## Structure

A PRM tool can be packaged as a `.tar.gz` or hosted in a Github repository.
Tools are how PRM knows to take defined actions against a file or folder; linting, testing, exporting, etc.

A tool can be:

- a Ruby gem,
- an executable binary,
- or a container.

## Composition

At a minimum, a PRM tool must contain a `prm-config.yml` in the root directory.

Optionally, you may include a `content` directory, also at the root.
The `content` directory can contain scripts, binaries, config files, etc. required to supplement the functionality of the tool outside the parameters of the configuration file.

When a tool is used with the Docker backend, everything in the `content` directory is mounted to `/tmp` in the container;
e.g. `content/myfile.sh` will be mounted to `/tmp/myfile.sh`.

### Required Parameters

All tools **must** include mandatory metadata in the `prm-config.yml` file which enumerate it for PRM to understand what it is.
They must include:

- the author (alphanumeric characters, dashes, and underscores only),
- a short id for the tool (alphanumeric characters, dashes, and underscores only),
- a display name (an arbitrary string which _can_ include spaces),
- a semantic version,
- and a url to the _upstream_ project the tool enables users to leverage (e.g., rubocop's landing page)

For example:

```yaml
plugin:
  author: fantastic-me
  id: awesome-tool
  display: Awesome Tool
  version: 0.1.0
  upstream_project_url: https://github.com/fantastic-me/awesome-me
```

Tools are executed by namespace, always in the form `author/id`.
To execute the example tool above, you would run the command `prm exec fantastic-me/awesome-tool`

### Common Parameters

In addition to the `plugin` configuration items, there is a `common` section that contains parameters that change the behavior of the tool regardless of underlying implementation details.

`can_validate`
: Set this to `true` if the tool is meant to validate a target code directory.
: Defaults to `false`

<!-- Force a break between definitions -->

`help_arg`
: Set this to the argument the tool expects to be passed to display its help info.
: Defaults to `--help`.

<!-- Uncomment when these when implemented
`interleave_stdout_err`
: Should the stdout & stderr be interleaved in to one stream, as opposed to separate ones?
: Defaults to `false`.

`needs_write_access`
: Will the execution of this tool require RW permissions against the target code dir?
: Defaults to `false`.
-->

`sucess_exit_code`
: Set this to the integer that the tool will exit with if it runs successfully.
: Defaults to `0`.

<!-- Force a break between definitions -->

`default_args`
: Set this to an array of arguments to pass to the tool by default; these can be overridden with the `toolArgs` flag when executed.
: No default value.

<!-- Force a break between definitions -->

`env`
: Set this to a map of environment variable names and their values to be set automatically prior to tool execution.
: No default value.

<!-- Uncomment when this is implemented
`output_mode`
: The output mode (`junit|json|yaml`)
: No default value.
-->

`requires_git`
: Set this to `true` if the tool requires a `git` binary.
: No default value.

<!-- Force a break between definitions -->

`use_script`
: Set this to the name of a script that will be invoked to execute the tool.
: The name **should not** include the file extension (e.g. `runme.sh` -> `use_script: runme`).
: The script **should** be present in the `content` directory and have execute permissions set.
: The script **should** have execute permissions set prior to packaging the tool (ie `chmod +x runme.sh`).
: No default value.

#### Example Common Parameters

```yaml
common:
  can_validate: true
  help_arg: '-h'
  success_exit_code: 2
  default_args: ['--include-tempfiles', '--fail-fast']
  env:
    - "TARGET_VERSION": "1.2.3"
    - "CONFIG_FILE": "/code/config.yaml"
  requires_git: true
  use_script: "collate_files_and_run"
```

This tool, as configured:

1. Can be used with `prm validate`
2. Will pass `-h` to retrieve the help information for the underlying tool
3. Will treat _only_ executions with an exit code of `2` as successful
4. Will by default append `--include-tempfiles --fail-fast` when executed
5. Will set the `TARGET_VERSION` environment variable to `1.2.3` and the `CONFIG_FILE` environment variable to `/code/config.yaml` in the execution context
6. Does require `git` to be installed/available
7. Will use the `collate_files_and_run.sh` in the `content` directory to execute the tool.

### Gem Tools

All `gem` based tools must declare which gems must be installed and the name of the executable to run.

For example:

```yaml
gem:
  name: [amazing_gem, dependency_gem, another_dependency_gem]
  executable: amazing_gem
```

By default, any gems specified will always attempt to resolve and use the _latest_ released version.

If there are compatibility concerns, version pins by Ruby version can be enumerated under the `compatibility` section.

For example:

```yaml
gem:
  name: [amazing_gem, dependency_gem, another_dependency_gem]
  executable: amazing_gem
  compatibility:
  - 2.6:
    - "amazing_gem": "2.15.0"
  - 2.7:
    - "amazing_gem": "~> 3.0"
```

This will install the latest versions of `dependency_gem` and `another_dependency_gem` regardless of Ruby version.

For any version of Ruby other than `2.6` and `2.7`, the tool will expect to use the latest version of `amazing_gem` as well.
For Ruby 2.6, `amazing_gem` will only ever be used at version `2.15.0`.
For Ruby 2.7, `amazing_gem` will always use the latest version greater than or equal to `3.0` and less than `4.0`.

### Binary Tools

All `binary` based tools must declare the name of the binary which is to be run and installation steps for acquiring that binary on particular platforms.

For example:

```yaml
binary:
  name: my_tool
  install_steps:
    windows: |
      # Windows specific install steps
      choco install my_tool
    linux: |
      # Linux specific install steps
      rpm install my_tool
    darwin: |
      # MacOS specific install steps
      brew install my_tool
```

This configuration will expect to find `my_tool` on the `PATH` and calls out installation steps for Windows, Linux, and MacOS.

> **Note:** At this time, only Linux docker containers are supported.
> Future releases will support other platforms and backends.

### Container Tools

All `container` based tools must declare the image name and tag to be used.
The images must be available or retrievable by Docker on the machine running PRM.

For example:

```yaml
container:
  name: myorg/myimage
  tag: latest
```

This configuration will execute in the `myorg/myimage:latest` container.
