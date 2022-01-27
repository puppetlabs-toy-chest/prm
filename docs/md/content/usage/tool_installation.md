---
title: "Installing Tools"
description: "Guidance for installing tools to PRM."
category: narrative
tags:
  - tools
  - usage
  - install-tool
  - update-tool
  - list-tools
weight: 10
---

This document explains how you can install tools to execute with PRM.

By default, tools are installed  to the `tools` folder alongside PRM unless the `toolpath` flag is specified with an alternate location.

### Local archive

Tool packages can be installed locally using the `prm install` command.

For example:

```bash
prm install ~/my-tool-1.2.3.tar.gz
```

That command will install the tool contained in `my-tool-1.2.3.tar.gz` to the default location.

### Remote archive

Tool packages stored remotely can be automatically downloaded and extracted with `prm install` so long as you know the URL to where the archive is.

For example:

```bash
prm install https://packages.mycompany.com/prm/my-tool-1.2.3.tar.gz
```

This command will attempt to download the PRM tool from the specified url and then install it like any other locally available PRM tool archive.

### Remote git repository

**Git** must be installed for this feature to work. The git repository must contain only one tool and must be structured with the `prm-config.yml` file and the `content` directory in the root directory of the repository.

For example:

```bash
prm install --git-uri https://github.com/myorg/myawesometool
```

This command will attempt to clone the PRM tool from the git repository at the specified URI and then install it to the default tool location.

### Force tool installation

Adding the `-f` or the `--force` flag to the `prm install` command will forcefully install/overwrite a tool if there is a tool installed with the same `author`, `name` and `version`.

### Updating tools

While there is no `update` command, newer versions of tools can be installed like any other tool package.

Currently, only the latest version of a selected tool is executable; the ability to select an older version of the tool to execute will be added in the future.

### List installed tools

Installed tools can be listed by running the command `prm exec --list`, displayed in the following format:

![prm tool list screenshot](../../images/exec-list-tools.png)

The `--toolpath` flag can also be added to list tools installed in an alternate location.
