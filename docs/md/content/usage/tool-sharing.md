---
title: "Sharing Tools"
description: "Learn how to share PRM tools."
category: narrative
tags:
  - tools
  - usage
---

After you've written your own tool you may wish to share it with other members of your team or the wider Puppet community. Work is underway to improve this initial functionality.

### prm build

This command will attempt to package the current working directory. You can change the directory to pack by providing `--sourcedir`.

``` bash
prm build [--sourcedir <dir>][--targetdir <dir>]
```

The `build` command will ensure that the directory that you are attempting to package will produce a valid Puppet Runtime Manager tool by looking for a `prm-config.yml` and a `content` directory.

The resulting `tar.gz` package will be created by default in `$cwd/pkg`. You can change the directory the package is created in by providing `--targetdir`.

> NOTE: When using `--sourcedir` use either the full path or the relative path without the leading `./` (`.\` on Windows).
> If you use the leading `./` it will incorrectly tar the project.
