---
title: "Readme"
description: "An overview of the PCT program."
category: concept
tags:
  - meta
draft: false
---

# Puppet Runtime Manager

- [Puppet Runtime Manager](#puppet-runtime-manager)
  - [Overview](#overview)
  - [Installing](#installing)
    - [Unix Systems](#unix-systems)
    - [Windows Systems](#windows-systems)
  - [Requesting a feature](#requesting-a-feature)
  - [Reporting Problems](#reporting-problems)
  - [Installing Telemetry Free Version](#installing-telemetry-free-version)
    - [Unix Systems](#unix-systems-1)
    - [Windows Systems](#windows-systems-1)

## Overview

The Puppet Runtime Manager (PRM) is a tool for validating Puppet content and for running arbitrary development/maintenance tasks against that content.

> :warning: PRM is currently in an EXPERIMENTAL phase and feedback is encouraged via [prm/discussions](https://github.com/puppetlabs/prm/discussions) and starting a `feedback` post.

## Installing

While PRM is in early release, we provide an archive and a simple script to unpack it.
When we move closer to a full release we will add a platform specific installer.
Use the `install.[ps1|sh]` script, depending upon your OS:

### Unix Systems

```bash
curl -L https://pup.pt/prm/install.sh | sh
```

### Windows Systems

```ps
iex "&{ $(irm https://pup.pt/prm/install.ps1); Install-Prm }"
```

This will install the latest release of PRM to `~/.puppetlabs/prm`.

<!-- This gif needs to be created once the scripts are callable -->
<!-- ![install_prm](docs/_resources/install_and_export_path.gif) -->

> :warning: If you do not use the install script and are extracting the archive yourself, be sure to use the fully qualified path to `~/.puppetlabs/prm` on *nix or `$HOME/.puppetlabs/prm` on Windows when you set your `PATH` environment variable.

A version of the product, with telemetry functionality disabled, is available too.
See [here](#installing-telemetry-free-version) for instructions on how to install it.

## Requesting a feature

Open a new feature request in our [Github discussion](https://github.com/puppetlabs/prm/discussions/new) page.

## Reporting Problems

If you're having trouble with the experimental PRM tool, please follow these instructions
to file an issue on our GitHub repository: https://github.com/puppetlabs/prm/issues/new

Make sure to fill in the information that is requested in the issue template as it
will help us investigate the problem more quickly.

## Installing Telemetry Free Version

We gather telemetry data to provide insights into how our products are being used.

The following data is collected:

- Version of application in use
- OS / platform of the device
- What commands have been invoked (including command args)
- Any errors that occurred when running the application

We understand that there will be some users who prefer to have no telemetry data sent.
For those users, we offer a version of PRM with the telemetry functionality disabled.

To install:
### Unix Systems

```bash
curl -L https://pup.pt/prm/install.sh | sh -s -- --no-telemetry
```

### Windows Systems

```ps
iex "&{ $(irm https://pup.pt/prm/install.ps1); Install-Prm -NoTelemetry }"
```

This will install the latest release of PRM, without telemetry functionality, to `~/.puppetlabs/prm`.
