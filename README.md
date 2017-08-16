# ccruncher - simple Cloud Controller log prettifier

*ccruncher* is a simple CLI tool to parse an existing CF Cloud Controller log and output it as a more readable yaml document.  The document is organized by apps and, within each app, by requests.  This allows a user to see what all the requests were for a given app, and within a given request, what all of the log messages were.  The idea is to make the CC log more usable/valuable for non-SMEs.

## Installation

```
git clone https://github.com/ScarletTanager/ccruncher
cd ccruncher
go build cmd/main.go
```

## Use

```
ccruncher cloud_controller_ng.log
```

The output is stored in the current working directory as `ccLog.yml`.

## Future ideas

- Interleaving of multiple cc log files
- Human timestamps
- Support for other CF entity types
    - Organizations first
    - Spaces?
    - What else?
- Storage of the parsed output as noSQL (e.g. etcd) documents for ease of querying/retrieval
