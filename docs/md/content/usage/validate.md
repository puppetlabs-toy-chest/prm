---
title: "Validating Puppet Content"
description: "Guidance for validating code with PRM."
category: narrative
tags:
- tools
- usage
- validation
weight: 10
---

In this guide you'll learn how the validate command works and how to use
it to quickly validate your Puppet content.

#### What is a validator
A validator is a tool which can read and validate a codebase of Puppet content.

The configuration of a validator is almost identical to that of a normal
tool. The only difference is that validators contain the `can_validate: true`
property in their `prm-config.yml` configuration file.

All validation tools in a tool directory can be listed with the following command:

```bash
prm validate --list
```

### How to use the validate command

#### Single tool validation

It is very easy to validate with a single tool. The command below shows an
outline of the core arguments and flags necessary for content validation:

```bash
prm validate {author/tool} --codedir {location/of/code} --toolArgs="{any, tool, args}"
```

Here is an example of a command composed to validate `epp` content against the current directory,
specifically, `dir1` and `dir2` as specified by the `toolArgs` flag.

```bash
prm validate puppetlabs/epp --codedir . --toolArgs="dir1, dir2"
```

When the command is executed PRM will validate with the `epp` tool using the appropriate tool
arguments; the tool arguments in this case will instruct epp to validate content in the specified directories.



#### Multi-tool validation

PRM supports running multiple validation tools simultaneously. That's where the
`validate.yml` file comes in; it allows you to arrange validators into groups.

After grouping validators in a `validate.yml` file, you can easily run a group of validators
from a single command. See below for usage of multi-tool validation.

##### Example of a `validate.yml` file

```yaml
groups:
  - id: "ci"
    tools:
      - name: puppetlabs/epp
        args: [--example, args]
      - name: puppetlabs/parser
      - name: puppetlabs/puppet-lint
      - name: puppetlabs/onceover
  - id: "quick_validate"
    tools:
      - name: puppetlabs/epp
      - name: puppetlabs/rubocop
      - name: puppetlabs/parser
      - name: puppetlabs/r10k
  - id: "syntax_validation"
    tools:
      - name: puppetlabs/puppet-syntax
      - name: puppetlabs/metadata-json-lint
      - name: puppetlabs/puppet-lint
```

{{% alert title="Note" color="primary" %}}
The `validate.yml` file should be stored in the **root** of the Puppet content codebase
which is to be validated.
{{% /alert %}}

##### `group` flag

The `--group {string}` flag can be used to specify which group of validators to run, e.g.

```bash
prm validate --codedir . --group syntax_validation
```

When the command is executed PRM will validate with the `syntax_validation` group of validators.

##### `workerCount` flag

The `--workerCount {int}` flag can be used to specify how many validators will run simultaneously; e.g.

```bash
prm validate --codedir . --group syntax_validation --workerCount 2
```

When the command is executed PRM will validate with the `syntax_validation` group of validators,
running a maximum of 2 validators concurrently.

##### `serial` flag

The `--serial` flag can be used to specify that only one validator should be run a time; e.g.

```bash
prm validate --codedir . --group syntax_validation --serial
```

When the command is executed PRM will validate with the `syntax_validation` group of validators,
running one validator at a time.




#### Viewing validation results

PRM can currently output validation results to the terminal or to a
log file.

{{% alert title="Note" color="primary" %}}
When outputting to the terminal, PRM will only output `stderr`, if there
is any; whereas both `stdout` and `stderr` are outputted to log files. This
is to prevent the terminal from becoming congested with verbose and
irrelevant text.
{{% /alert %}}


##### Default output locations

See below the default output locations for each type of validation:

| Single tool | Multi-tool |
|-------------|------------|
| terminal    | file       |

##### `resultsView` flag

The `--resultsView {terminal|file}` can be used to specify where you want the results of validation to output
to.

```bash
prm validate puppetlabs/epp --codedir . --resultsView file
```

The results of this validation will be sent to a log file, overriding the default output for single tool
validation, which is the terminal.

#### Example of multi-tool results outputted to the terminal:

```bash
$ prm validate --codedir . --group syntax_validation --resultsView terminal
3:23PM INF Found tool group: syntax_validation 
3:23PM INF Validating with the puppet-syntax tool
3:23PM INF Validating with the puppet-lint tool
3:23PM INF Validating with the metadata-json-lint tool
3:24PM ERR puppet-syntax:
---> syntax:manifests
---> syntax:templates
ERRORS:
Syntax error at 'Kernel' (file: templates/motd.epp, line: 5, column: 1)

      TOOL NAME      | VALIDATION EXIT CODE  
---------------------+-----------------------
  puppet-syntax      |                    1  
  metadata-json-lint |                    0  
  puppet-lint        |                    0  
3:24PM ERR Validation returned 1 error
```

#### Example of multi-tool results outputted to the log files:

```bash
# No need to specify the --resultsView flag here as it defaults to `file` for mulit-tool validation
$ prm validate --codedir . --group syntax_validation
3:49PM INF Found tool group: syntax_validation 
3:49PM INF Validating with the metadata-json-lint tool
3:49PM INF Validating with the puppet-syntax tool
3:49PM INF Validating with the puppet-lint tool

      TOOL NAME      | VALIDATION EXIT CODE |                                 FILE LOCATION
---------------------+----------------------+--------------------------------------------------------------------------------
  puppet-syntax      |                    1 | .prm-validate/syntax_validation/puppet-syntax_2022_April_26_16-49-59.log
  metadata-json-lint |                    0 | .prm-validate/syntax_validation/metadata-json-lint_2022_April_26_16-49-59.log
  puppet-lint        |                    0 | .prm-validate/syntax_validation/puppet-lint_2022_April_26_16-49-59.log
3:49PM ERR Validation returned 1 error
```

{{% alert title="Note" color="primary" %}}
Support for formatting of the validation results will be implemented in a future 
release. E.g. JSON or JUNIT
{{% /alert %}}


