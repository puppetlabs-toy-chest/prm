---
title: "Telemetry"
description: "An overview of PRM's telemetry."
category: concept
tags:
  - meta
draft: false
weight: 10
---

By default, PRM collects some telemetry during each run.
This information is non-identifying and does not include any personal or organizational details.

We use it to see how PRM is being used so we can better understand issues, features, and what to prioritize as we move forward.
We emit this telemetry information to [Honeycomb](https://honeycomb.io) for introspection.
We keep the telemetry for 60 days, after which it ages out.

## What We Gather

Each PRM run reports the following data:

- The UTC timestamp for when the run started
- How long the command took to run
- The operating system architecture (x64, ARM, etc)
- The operating system type (darwin, linux, windows)
- A uniquely generated internal identifier for the machine running the command
- Whether or not the command ran in a CI environment
- The arguments passed to PRM in this run

This data allows us to see broadly how many machines the program is running on, what the breakdown of operating system/system architecture is, how performant the program is, and which commands are most used.

When there is an error, PRM reports the error information as well - the stack trace information and the error message displayed to a user.
This lets us investigate issues to see how common they are and improve our testing and error handling.

## Running PRM Without Telemetry

We know that some users do not want to report telemetry;
for those users, it is possible to install a "notel" version of the binary which does not include the telemetry bits at all.
This binary is identical to the standard binary save that it does not have any mechanisms for recording or sending telemetry data.

For those users, an alternate install is available.

To install:

### Bash

```bash
curl -L https://pup.pt/pdkgo/install.sh | sh -s -- --no-telemetry
```

### PowerShell

```powershell
iex "&{ $(irm 'https://pup.pt/pdkgo/install.ps1'); Install-Pct -NoTelemetry }"
```
