# Sensu Go Elasticsearch metric handler plugin
[![Bonsai Asset Badge](https://img.shields.io/badge/CHANGEME-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/CHANGEME/CHANGEME) [![TravisCI Build Status](https://travis-ci.org/CHANGEME/sensu-go-elasticsearch.svg?branch=master)](https://travis-ci.org/CHANGEME/sensu-go-elasticsearch)

TODO: Description.

## Installation

Download the latest version of the sensu-go-elasticsearch from [releases][1],
or create an executable script from this source.

From the local path of the sensu-go-elasticsearch repository:

```
go build -o /usr/local/bin/sensu-go-elasticsearch main.go
```

## Configuration

Example Sensu Go definition:

```json
{
    "api_version": "core/v2",
    "type": "CHANGEME",
    "metadata": {
        "namespace": "default",
        "name": "CHANGEME"
    },
    "spec": {
        "...": "..."
    }
}
```

## Usage Examples

Help:

```
The Sensu Go CHANGEME for x

Usage:
  sensu-go-elasticsearch [flags]

Flags:
  -f, --foo string   example
  -h, --help         help for sensu-go-elasticsearch
```

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/CHANGEME/sensu-go-elasticsearch/releases
