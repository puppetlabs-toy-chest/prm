---
title: "Using the Explain Subsystem"
description: "How to use the Explain command to discover documentation at the commandline."
category: narrative
tags:
  - fundamentals
weight: 20
---

The Explain subsystem is a built-in way to surface PRMs full documentation at the commandline with pretty formatting.
While our documentation is published online, sometimes it is useful to read it in the terminal you're working from.

There are three broad categories of documentation used in this project: reference, narrative, and concept.
The Explain subsystem is used to surface narrative and concept docs.

Reference docs are surfaced using the built-in help subsystem by calling `prm help <command>` or `prm <command> --help`.

Narrative docs are walkthroughs or guides which take you through a specific action or problem.
Concept docs are used to explain broader or deeper topics, like configuration and telemetry.

This guide will walk you through how to look for and read narrative and concept documents from the terminal.

## Searching for Topics

If you run `prm explain` without any additional parameters, it will list all available topics into your terminal in a table format.
This will include the name, description, category, and any tags associated with each topic.

You can also filter by category using the `category` and `tag` flags, including using them together.
Moreover, PRM already knows the available tags and categories, so if you're using shell completion, you can tab or step through them.

For example, the command below will display any topics which are part of the `narrative` category and include the `install` tag:

```sh
prm explain --category narrative --tag install
```

You can also use shell completion to auto-complete the list of available topics directly.

## Displaying a Single Topic

Once you have found the name of the topic you want to read, either by searching the list of topics or auto-completing in the terminal, you can run `prm explain <name>` to display it.

For example, to read this documentation:

```sh
prm explain explain-subsystem
```

Which will render this document to your terminal with syntax highlighting.

## Documentation Versions

The version of the documentation available in the terminal is always the documents as they existed when that version was shipped.
For updated docs, you will need to read the documentation online, either in [Github](https://github.com/puppetlabs/prm/tree/main/docs/md/content) or on our [docs site](https://puppetlabs.github.io/devx/prm/).
