# tapir-mgmt

CLI utility primarily to manage select TAPIR Core tasks like
global configuration, status reports from TAPIR Edge, etc.
For this the config details needed to connect to eg. TAPIR-SLOGGER are
located in the config files in the /etc/dnstapir/ directory.

For some uses, `tapir-mgmt` is used in "standalone" mode. 

`tapir-mgmt` has a large number of commands and subcommands. The entire
set of commands is structured as a tree with the root in the
`tapi-mgmt` command.  All commands, regardless of where in the tree of
commands they are located, have online help via the flag `-h`. I.e. to
get help on the `tapir-mgmt slogger ping` command, run:

```
tapir-mgmt slogger ping -h
Send an API ping request to TAPIR-SLOGGER and present the response

Usage:
  tapir-mgmt slogger ping [flags]

Flags:
  -c, --count int   #pings to send
  -h, --help        help for ping
  -n, --newapi      use new api client

Global Flags:
      --config string   config file (default is /etc/dnstapir/tapir-pop.yaml)
  -d, --debug           Debugging output
  -H, --headers         Show column headers
      --tls             Use a TLS connection to TAPIR-POP (default true)
  -v, --verbose         Verbose mode
```

The flag `-h` also lists all subcommands underneath the command in question.

The `tapir-mgmt` command has a number of subcommands, each of which is a command group. The command groups are:

```
tapir-mgmt -h                
CLI  utility used to interact with TAPIR-SLOGGER, i.e. the TAPIR Status Logger, among
other tasks

Usage:
  tapir-mgmt [command]

Available Commands:
  api         request a TAPIR-POP api summary
  completion  Generate the autocompletion script for the specified shell
  debug       Prefix command to various debug tools; do not use in production
  help        Help about any command
  keyupload   Upload a public key to a TAPIR Core
  mqtt        Prefix command, not usable directly
  slogger     Prefix command to TAPIR-Slogger, only usable in TAPIR Core, not in TAPIR Edge

Flags:
      --config string   config file (default is /etc/dnstapir/tapir-mgmt.yaml)
  -d, --debug           Debugging output
  -H, --headers         Show column headers
  -h, --help            help for tapir-cli
      --tls             Use a TLS connection to TAPIR-SLOGGER (default true)
  -v, --verbose         Verbose mode

Use "tapir-mgmt [command] --help" for more information about a command.
```

Some of the commands are only there as debugging tools. They are not intended for use in production. 
