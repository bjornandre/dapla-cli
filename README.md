# dapla-cli

The dapla cli is a command-line application users can use to interract with the da(ta)pla(form). The command 
has several sub-commands.

```
# dapla --help
The dapla command is a collection of utilities you can use with the dapla platform.

Usage:
  dapla [command]

Available Commands:
  completion  Generate completion script
  help        Help about any command
  ls          List information about the dataset(s) under PATH
  rm          Remove the dataset(s) under PATH

Flags:
  -h, --help            help for dapla
      --jupyter         fetch the Bearer token from jupyter
  -s, --server string   set URI of the API server
      --token string    set the Bearer token to use to authenticate with the server

Use "dapla [command] --help" for more information about a command.
```

## Installation

The command is already installed in the dapla jupyterlab environement. To install the command locally extract the content of the release archive on your computer and alias the `dapla-cli` executable to `dapla`.

## Authentication

In order to be able to communicate with the API servers one need to provide an authentication methods and the API server URI. 

The flags `--jupyter` can be used when the dapla command runs inside the container. In this case the application will try to retrieve the authentication token by itself: 

`# dapla --jupyter --server "https://server-api/"`

Alternatively one can provide an authentication token manually using the `--token` flag:

`# dapla --token "my.jwt.token" --server "https://server-api/"`

## Commands

### ls (list)

```
$ dapla ls --help 
Usage:
  dapla ls [PATH]... [flags]

Flags:
  -l, --       use a long listing format
  -h, --help   help for ls

$ dapla ls /
/felles/
/kilde/
/produkt/
/raw/
/skatt/
/tmp/
/user/
```

### rm (remove)

### completion
