# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org).
## [0.2.1](https://github.com/puppetlabs/prm/tree/0.2.1) - 2022-05-31

[Full Changelog](https://github.com/puppetlabs/prm/compare/0.2.0...0.2.1)

### Added

- (MAINT) Better error outputting to terminal [#174](https://github.com/puppetlabs/prm/pull/174) ([petergmurphy](https://github.com/petergmurphy))

### Fixed

- (GH-171) Better handle docker not running error [#173](https://github.com/puppetlabs/prm/pull/173) ([petergmurphy](https://github.com/petergmurphy))

- (GH-163) Fix `--codedir` flag default [#165](https://github.com/puppetlabs/prm/pull/165) ([petergmurphy](https://github.com/petergmurphy))

## [0.2.0](https://github.com/puppetlabs/prm/tree/0.2.0) - 2022-05-05

[Full Changelog](https://github.com/puppetlabs/prm/compare/0.1.0...0.2.0)

### Added

- (GH-130) Implement multi-tool validate [#148](https://github.com/puppetlabs/prm/pull/148) ([petergmurphy](https://github.com/petergmurphy))

- (GH-129) Tool output handling [#145](https://github.com/puppetlabs/prm/pull/145) ([petergmurphy](https://github.com/petergmurphy))

- (GH-341) Add build command [#123](https://github.com/puppetlabs/prm/pull/123) ([michaeltlombardi](https://github.com/michaeltlombardi))

### Fixed

- (GH-138) Ensure cache directory exists [#160](https://github.com/puppetlabs/prm/pull/160) ([petergmurphy](https://github.com/petergmurphy))

- (GH-137) Fix tool listing order [#159](https://github.com/puppetlabs/prm/pull/159) ([petergmurphy](https://github.com/petergmurphy))

- (GH-139) Use puppet-agent 7.15.0 [#154](https://github.com/puppetlabs/prm/pull/154) ([chelnak](https://github.com/chelnak))

## [0.1.0](https://github.com/puppetlabs/prm/tree/0.1.0) - 2022-01-31

[Full Changelog](https://github.com/puppetlabs/prm/compare/9d7ff6acd89f19cd24492acc9058c4c5752bcbe5...0.1.0)

### Added

- (GH-85) Add `alwaysBuild` flag and functionality [#104](https://github.com/puppetlabs/prm/pull/104) ([petergmurphy](https://github.com/petergmurphy))

- (GH-54) - Implement a timeout flag [#86](https://github.com/puppetlabs/prm/pull/86) ([david22swan](https://github.com/david22swan))

- (GH-11) Add explain command and docs package [#83](https://github.com/puppetlabs/prm/pull/83) ([petergmurphy](https://github.com/petergmurphy))

- (GH-63) Allow additional tool args to be set within validate.yml  [#64](https://github.com/puppetlabs/prm/pull/64) ([da-ar](https://github.com/da-ar))

- (GH-13) Implement PRM install [#42](https://github.com/puppetlabs/prm/pull/42) ([petergmurphy](https://github.com/petergmurphy))

- (GH-14) Initial status implementation [#40](https://github.com/puppetlabs/prm/pull/40) ([michaeltlombardi](https://github.com/michaeltlombardi))

- (GH-8) PRM exec [#33](https://github.com/puppetlabs/prm/pull/33) ([da-ar](https://github.com/da-ar))

- (GH-15) Implement prm set|get backend [#32](https://github.com/puppetlabs/prm/pull/32) ([sanfrancrisko](https://github.com/sanfrancrisko))

- (GH-9) Add EPP validator tool example [#29](https://github.com/puppetlabs/prm/pull/29) ([michaeltlombardi](https://github.com/michaeltlombardi))

- (GH-16) Add `prm (set|get) puppet` commands [#22](https://github.com/puppetlabs/prm/pull/22) ([sanfrancrisko](https://github.com/sanfrancrisko))

- (GH-18) Add docker backend status method [#21](https://github.com/puppetlabs/prm/pull/21) ([michaeltlombardi](https://github.com/michaeltlombardi))

- (GH-265) Scaffold tool & related structs [#2](https://github.com/puppetlabs/prm/pull/2) ([michaeltlombardi](https://github.com/michaeltlombardi))

### Fixed

- (GH-95) Add installing and updating tools docs [#113](https://github.com/puppetlabs/prm/pull/113) ([petergmurphy](https://github.com/petergmurphy))

- (GH-97) Document explain subsystem [#112](https://github.com/puppetlabs/prm/pull/112) ([michaeltlombardi](https://github.com/michaeltlombardi))

- (GH-88) Draft install docs [#108](https://github.com/puppetlabs/prm/pull/108) ([michaeltlombardi](https://github.com/michaeltlombardi))

- (GH-94) Draft quickstart documentation [#106](https://github.com/puppetlabs/prm/pull/106) ([michaeltlombardi](https://github.com/michaeltlombardi))

- (GH-93) Draft Anatomy of a Tool concept document [#72](https://github.com/puppetlabs/prm/pull/72) ([sanfrancrisko](https://github.com/sanfrancrisko))

- (GH-36) Fix NPE when puppet version unset in cfg on disk [#37](https://github.com/puppetlabs/prm/pull/37) ([sanfrancrisko](https://github.com/sanfrancrisko))
