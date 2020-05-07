# xpcutil
Golang Tool to interact with Launchd and other services with XPC

Original source code for `xpc_wrapper.c` and `xpc_wrapper.h` can be found here: http://newosxbook.com/articles/jlaunchctl.html

To build: `go build -o xpcutil cmd/main.go`
### Commands

- list: List all available services
- send: Send an XPC message to a target service
- start: Start a service
- stop: Stop a service
- submit: Submit a launchd job (WIP)
- load: Load a plist file with launchd
- unload: Unload a plist file
- procinfo: Obtain process information for a given pid
- status: Obtain status information about a given service
- listen: Create an xpc service and listen for connections (Not implemented)

### Arguments
- data: Base64 encoded json data to send to a target service
- file: Path to the plist file for load/unload commands
- keepalive: Set the keepalive flag to true for the submit command
- pid: Target process id for the procinfo command
- privileged:  If set to true the XPC_CONNECTION_MACH_SERVICE_PRIVILEGED flag will be used when connecting to the service (default false)
- program: Command to execute with the submit command
- service: Target XPC service name for the start, stop, and status commands. When using the submit command, this will be the label