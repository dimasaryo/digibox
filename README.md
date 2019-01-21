# Digibox
digitalocean development box cli

## Installation

If you don't have glide, get dependencies from `glide.yaml` using go get
```
$ go get github.com/dimasaryo/digibox
$ cd $GOPATH/src/github.com/dimasaryo/digibox
$ glide install
$ go install
```
If you have added <<your GOBIN to your path, you should able to execute `digibox` command
```
$ digibox --help
NAME:
   digibox - Digitalocean remote development server cli

USAGE:
   digibox [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     start, s  start a remote development server `NAME`
     stop      stop remote development server `NAME`
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help                                                           
   --version, -v  print the version
```

## Usage

### Start Development Box
```
$ digibox start <name>
```

### Stop Development Box
```
$ digibox stop <name>
```
