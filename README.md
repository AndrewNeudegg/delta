# Delta


![build](https://github.com/AndrewNeudegg/delta/workflows/build/badge.svg) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta?ref=badge_shield) [![codecov](https://codecov.io/gh/AndrewNeudegg/delta/branch/main/graph/badge.svg?token=PZNGIZGN2V)](https://codecov.io/gh/AndrewNeudegg/delta) [![Go Report Card](https://goreportcard.com/badge/github.com/andrewneudegg/delta)](https://goreportcard.com/report/github.com/andrewneudegg/delta) ![coverage](https://github.com/AndrewNeudegg/delta/workflows/coverage/badge.svg) ![LOC](https://sloc.xyz/github/andrewneudegg/delta)


An easy to understand and pluggable eventing system.

## Overview

Delta is a run anywhere binary that allows you to easily dictate the flow of events in your system. To achieve this
delta uses configuration files to specify where events come from, what should happen to them during
processing, and where they should be sent.

The specific configuration recipe that you will need will be a combination of the avalibale ingredients with some additions and some removals. For some inspiration, take a look in the [`./recipes`](https://github.com/AndrewNeudegg/delta/tree/main/recipies/simulation) directory.

To get the latest binaries head over to the [releases](https://github.com/AndrewNeudegg/delta/releases) page or pull
the docker image `docker run -it andrewneudegg/delta`.

## Table of Contents

- [Delta](#delta)
  - [Overview](#overview)
  - [Table of Contents](#table-of-contents)
  - [TL;DR (Example Recipes)](#tldr-example-recipes)
    - [HTTP Server](#http-server)
    - [Encryption](#encryption)
    - [Performance](#performance)
  - [License](#license)

## TL;DR (Example Recipes)

### HTTP Server

```yaml
applicationSettings: {}
pipeline:
  # This first pipeline generates and emits http events.
  - id: pipelines/fipfo
    config:
      input:
        - id: utilities/generators/v1
          config:
            interval: 10s
            numberEvents: 10000
            numberCollections: 1
      output:
        - id: http/v1
          config:
            targetAddress: http://localhost:8080

  # This second pipeline consumes those events and writes to stdout.
  - id: pipelines/fipfo
    config:
      input:
        - id: http/v1
          config:
            listenAddress: :8080
            maxBodySize: 1000000 # 1mb
      output:
        - id: utilities/performance/v1
          config:
            sampleWindow: 60s
          nodes:
            - id: utilities/console/v1
```

### Encryption

```yaml
applicationSettings: {}
pipeline:
  - id: pipelines/fipfo
    config:
      input:
        # The crypto resource wraps the generator resource.
        - id: utilities/crypto/v1
          config:
            mode: encrypt
            password: iddX0DQKGMl7LszqdDKUL6aFVvMGAtwd
          nodes:
            - id: utilities/generators/v1
              config:
                interval: 1s
                numberEvents: 1
                numberCollections: 1
      output:
        # Decrypt the events before writing the output to console.
        - id: utilities/crypto/v1
          config:
            mode: decrypt
            password: iddX0DQKGMl7LszqdDKUL6aFVvMGAtwd
          nodes:
            - id: utilities/console/v1
```

### Performance

```yaml
applicationSettings: {}
pipeline:
  - id: pipelines/fipfo
    config:
      input:
        # Wrap a resource with a performance measuring resource.
        - id: utilities/performance/v1
          config:
            sampleWindow: 10s
          nodes:
            - id: utilities/generators/v1
              config:
                interval: 1s
                numberEvents: 1000
                numberCollections: 1000
      output:
        - id: utilities/console/v1
```

## License

[MIT LICENSE](./LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FAndrewNeudegg%2Fdelta?ref=badge_large)
