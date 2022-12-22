# CLI for datarhei Core

A CLI for connection to instances of datarhei Core. Currently only Core versions 16+ are supported.

## Installation

```
go install github.com/datarhei/core-cli-go@latest
```

## Usage

You see all commands when you call the binary without any arguments:

```
$ corecli                                                                                                                                git:main cmd:0
A CLI for the datarhei Core.

Usage:
  corecli [command]

Available Commands:
  config      Config related commands
  core        Core related commands
  editor      Editor related commands
  fs          Filesystem related commands
  help        Help about any command
  log         Logging related commands
  metadata    Metadata related commands
  metrics     Metrics related commands
  process     Process related commands
  skills      FFmpeg skills related commands
  srt         SRT related commands

Flags:
      --config string   config file (default is $HOME/.corecli.json)
  -h, --help            help for corecli

Use "corecli [command] --help" for more information about a command.
```

## Quick start

You first have to add a core to the list of known cores and give it the name `mycore`:

```
corecli core add mycore http://127.0.0.1:8080 --username admin --password datarhei
```

Now you can list all known cores:

```
$ corecli core list
┌───────-─┬────────────────────────────┬─────────┬──────┬────┐
│ NAME    │ HOST                       │ VERSION │ NAME │ ID │
├───────-─┼────────────────────────────┼─────────┼──────┼────┤
│ *mycore │ http://127.0.0.1:8080      │         │      │    │
└──────-──┴────────────────────────────┴─────────┴──────┴──-─┘
```

The star denotes the currently selected core that all other commands refer to.

Try to connect to the core by issuing any command, e.g.

```
$ corecli core about
raspy-lake-6361 16.11.0 (darwin/amd64) 51c6f26c-2352-4403-95d3-a008d9edd81d @ http://127.0.0.1:8080
{
  "app": "datarhei-core",
  "auths": [],
  "created_at": "2022-12-22T07:17:28+01:00",
  "id": "51c6f26c-2352-4403-95d3-a008d9edd81d",
  "name": "raspy-lake-6361",
  "uptime_seconds": 11,
  "version": {
    "arch": "darwin/amd64",
    "build_date": "",
    "compiler": "go1.19.4",
    "number": "16.11.0",
    "repository_branch": "",
    "repository_commit": ""
  }
}
```

## Adding a process

Print out a template for the process configuration:

```
corecli process template > process.json
```

Now modify the process configuration to your needs and create to the process on the core:

```
corecli process add < process.json
```

The moment the process is created on the core you can edit it on-the-fly (assume that the process has the ID `example`):

```
corecli process edit example
```

This will open the process config in the system's default editor. If you want to use a different editor, specify it with e.g. `core editor set /usr/bin/nano`). After saving the process config, it will replace the one currently on the core.

Use the process commands in order to control the process:

```
corecli process start example
```

List all processes on the core:

```
corecli process list
```

Get details for one specific process:

```
corecli process show example
```

Remove the process:

```
corecli process delete example
```
