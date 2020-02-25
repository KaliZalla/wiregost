## Completers

The `completers` package contains all completers used in the console. It works as follows:
A main/root completer is registered during console instantiation, and provides completion only
for root commands.

Each time a root command needs further completion (for its subcommands, arguments, or data_service objects),
the root completer calls specialized functions contained in their appropriate source files.

----
#### Main menu

* `completer.go`        - Main console completer, calling other specialized completer functions
* `workspace.go`        - Workspace subcommands and objects
* `hosts.go`            - Hosts subcommands and objects
* `module.go`           - Completes paths for all modules
* `stack.go`            - Completes module paths for stack-loaded modules
* `option.go`           - Completes module options
* `profiles.go`         - Completes implant profiles names 
* `user.go`             - Completes user option fields 
* `server.go`           - Completes server configs for server switching 

#### Implant menu
* `agent-help.go`       - Completes agent help commands 
