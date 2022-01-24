---
title: "Installation"
description: "Steps to install PRM on Windows, macOS and Linux."
category: narrative
tags:
  - fundamentals
  - install
weight: 10
---

While PRM is in early release, we provide an archive and a simple script to unpack it.
When we move closer to a full release we will add a platform specific installer.
Use the `install.[ps1|sh]` script, depending upon your OS:

### Bash

```bash
curl -L https://pup.pt/prm/install.sh | sh
```

### PowerShell

```powershell
iex "&{ $(irm 'https://pup.pt/prm/install.ps1'); Install-Prm }"
```

This will install the latest release of PRM to `~/.puppetlabs/prm`.

> **Warning!**
>
> If you do not use the install script and are extracting the archive yourself, be sure to use the fully qualified path to `~/.puppetlabs/prm` on *nix or `$HOME/.puppetlabs/prm` on Windows when you set your `PATH` environment variable.

A version of the product, with telemetry functionality disabled, is available too.
See [here](#installing-telemetry-free-version) for instructions on how to install it.

### Setting up Tab Completion

After installation, we'd highly recommend setting up tab completion for your shell to ensure the best possible experience.

PRM has built in tab completion support for the following shells: `bash`, `zsh`, `fish` and `powershell`.

To view the install instructions, run `prm completion --help` and follow the instructions for your shell.

## Installing Telemetry Free Version

As of `0.5.0`, we have been gathering telemetry data to provide insights in to how our products are being used.

The following data is collected:

- Version of application in use
- OS / platform of the device
- What commands have been invoked (including command args)
- Any errors that occurred when running the application

We understand that there will be some users who prefer to have no telemetry data sent.
For those users, we offer a version of PRM with the telemetry functionality disabled.

To install:

### Bash

```bash
curl -L https://pup.pt/prm/install.sh | sh -s -- --no-telemetry
```

### PowerShell

```powershell
iex "&{ $(irm 'https://pup.pt/prm/install.ps1'); Install-Prm -NoTelemetry }"
```

This will install the latest release of PRM, without telemetry functionality, to `~/.puppetlabs/prm`.
