# Delta


![build](https://github.com/AndrewNeudegg/delta/workflows/build/badge.svg) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta?ref=badge_shield) [![codecov](https://codecov.io/gh/AndrewNeudegg/delta/branch/main/graph/badge.svg?token=PZNGIZGN2V)](https://codecov.io/gh/AndrewNeudegg/delta) [![Go Report Card](https://goreportcard.com/badge/github.com/andrewneudegg/delta)](https://goreportcard.com/report/github.com/andrewneudegg/delta) ![coverage](https://github.com/AndrewNeudegg/delta/workflows/coverage/badge.svg) ![LOC](https://sloc.xyz/github/andrewneudegg/delta)


An easy to understand and pluggable eventing system.

## Overview

This repository is a 'starter-pack' with a few good default recipes. The specific recipe that you will need
will be a combination of the avalibale ingredients with some additions and some removals.

## Table of Contents

- [Delta](#delta)
  - [Overview](#overview)
  - [Table of Contents](#table-of-contents)
  - [Components](#components)
    - [Sink](#sink)
    - [Bridge](#bridge)
    - [Distributor](#distributor)
  - [License](#license)

## Components

### Sink

To sink events you must have an external source of event data that is capable of HTTP posting data.

Please see the [sink](cmd/sink/README.md) specific documentation for more information.

### Bridge

To bridge events you must have an external source of event data that can be consumed.

Please see the [bridge](cmd/bridge/README.md) specific documentation for more information.

### Distributor

Once you have a source of events, either a sink or a bridge, you can distribute the events to your target applications.

Please see the [distributor](cmd/distributor/README.md) specific documentation for more information.


## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta?ref=badge_large)