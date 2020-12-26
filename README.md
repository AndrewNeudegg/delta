# delta
![build](https://github.com/AndrewNeudegg/delta/workflows/build/badge.svg) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta?ref=badge_shield) [![codecov](https://codecov.io/gh/AndrewNeudegg/delta/branch/main/graph/badge.svg?token=PZNGIZGN2V)](https://codecov.io/gh/AndrewNeudegg/delta) [![Go Report Card](https://goreportcard.com/badge/github.com/andrewneudegg/delta)](https://goreportcard.com/report/github.com/andrewneudegg/delta) ![coverage](https://github.com/AndrewNeudegg/delta/workflows/coverage/badge.svg)


An easy to understand and pluggable Kubernetes Native eventing system.

## Overview

This repository is a 'starter-pack' with a few good default recipes. The specific recipe that you will need
will be a combination of the avalibale ingredients with some additions and some removals.



## Sink

To sink events you must have an external source of event data that is capable of HTTP posting data.

Please see the [sink](cmd/sink/README.md) specific documentation for more information.

## Bridge

To bridge events you must have an external source of event data that can be consumed.

Please see the [bridge](cmd/bridge/README.md) specific documentation for more information.

## Table of Contents

- [delta](#delta)
  - [Overview](#overview)
  - [Sink](#sink)
  - [Bridge](#bridge)
  - [Table of Contents](#table-of-contents)
  - [License](#license)


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta?ref=badge_large)