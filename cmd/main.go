package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xorrior/xpcutil/pkg/xpc"
)

type XpcMan struct {
	xpc.Emitter
	conn        xpc.XPC
	logger      *log.Logger
	servicename string
}

func New(service string, privileged int) *XpcMan {
	x := &XpcMan{logger: log.New(os.Stdout, fmt.Sprintf("%s>", service), log.Lshortfile|log.LstdFlags), Emitter: xpc.Emitter{}}
	x.Emitter.Init()
	x.conn = xpc.XpcConnect(service, x, privileged)
	return x
}

func (x *XpcMan) XpcManLog(msg string) {
	if x.logger != nil {
		x.logger.Println(msg)
	}
}

func (x *XpcMan) HandleXpcEvent(event xpc.Dict, err error) {
	if err != nil {
		x.XpcManLog(fmt.Sprintf("error: %s", err.Error()))
		if event == nil {
			return
		}
	}

	// marshal the xpc.Dict object to raw, indented json and print it
	raw, err := json.Marshal(event)
	if err != nil {
		x.XpcManLog(fmt.Sprintf("error: %s", err.Error()))
		return
	}

	x.XpcManLog(fmt.Sprintf("%s\n", string(raw)))
	return
}

// Author: @xorrior, @raff
func main() {
	// Main function

	// Declare command line args
	command := flag.String("command", "", "Command to execute. Avalable commands:\n list: List all available services \n send: Send data to a target service \n start: start a service\n stop: stop a service\n submit: Submit a launchd job \n load: Load a plist file with launchd \n unload: Unload a plist file \n procinfo: Obtain process information for a given pid \n status: Obtain status information about a given service \n listen: Create an xpc service and listen for connections [Not Implemented]")
	serviceName := flag.String("service", "", "Service name/ Bundle ID for xpc service to target.\n When using the submit sub command, the serviceName will be used for the label.")
	program := flag.String("program", "", "command to execute with the submit sub command")
	keepalive := flag.Bool("keepalive", false, "keep the program alive. Used with the submit sub command")
	file := flag.String("file", "", "Path to the plist file for the load/unload commands")
	pid := flag.Int("pid", 0, "Target process ID for process info")
	data := flag.String("data", "", "Base64 encoded json data to send to the target service")
	privileged := flag.Bool("privileged", false, "If set to true the XPC_CONNECTION_MACH_SERVICE_PRIVILEGED flag will be used when connecting to the service")

	flag.Parse()

	switch *command {
	case "list":
		// List the available services buy sending an xpc message to launchd
		if len(*serviceName) == 0 {
			response := xpc.XpcLaunchList("")
			response = response.(xpc.Dict)
			raw, err := json.MarshalIndent(response, "", "	")
			if err != nil {
				fmt.Println("Error serializing golang object ", err.Error())
			}
			fmt.Printf("%s\n", string(raw))
			//fmt.Printf("Result: \n %+v\n", response)
		} else {
			response := xpc.XpcLaunchList(*serviceName)
			response = response.(xpc.Dict)
			raw, err := json.MarshalIndent(response, "", "	")
			if err != nil {
				fmt.Println("Error serializing golang object ", err.Error())
			}
			fmt.Printf("%s\n", string(raw))
			//fmt.Printf("Result: \n %+v\n", response)
		}
		break
	case "start":
		if len(*serviceName) == 0 {
			fmt.Println("missing service name")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else {
			response := xpc.XpcLaunchControl(*serviceName, 1)
			response = response.(xpc.Dict)
			raw, err := json.MarshalIndent(response, "", "	")
			if err != nil {
				fmt.Println("Error serializing golang object ", err.Error())
			}
			fmt.Printf("%s\n", string(raw))
		}
		break
	case "stop":
		if len(*serviceName) == 0 {
			fmt.Println("missing service name")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else {
			response := xpc.XpcLaunchControl(*serviceName, 0)
			response = response.(xpc.Dict)
			raw, err := json.MarshalIndent(response, "", "	")
			if err != nil {
				fmt.Println("Error serializing golang object ", err.Error())
			}
			fmt.Printf("%s\n", string(raw))
		}
		break
	case "load":
		if len(*file) == 0 {
			fmt.Println("Missing file argument")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else {
			response := xpc.XpcLaunchLoadPlist(*file)
			response = response.(xpc.Dict)
			raw, err := json.MarshalIndent(response, "", "	")
			if err != nil {
				fmt.Println("Error serializing golang object ", err.Error())
			}
			fmt.Printf("%s\n", string(raw))
		}
		break
	case "unload":
		if len(*file) == 0 {
			fmt.Println("Missing file argument")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else {
			response := xpc.XpcLaunchUnloadPlist(*file)
			response = response.(xpc.Dict)
			raw, err := json.MarshalIndent(response, "", "	")
			if err != nil {
				fmt.Println("Error serializing golang object ", err.Error())
			}
			fmt.Printf("%s\n", string(raw))
		}
		break
	case "status":
		if len(*serviceName) == 0 {
			fmt.Println("missing service name")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else {
			response := xpc.XpcLaunchStatus(*serviceName)
			response = response.(xpc.Dict)
			raw, err := json.MarshalIndent(response, "", "	")
			if err != nil {
				fmt.Println("Error serializing golang object ", err.Error())
			}
			fmt.Printf("%s\n", string(raw))
		}
		break
	case "procinfo":
		if *pid == 0 {
			fmt.Println("missing target pid")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else {
			response := xpc.XpcLaunchProcInfo(*pid)
			fmt.Println(response)
		}
		break
	case "send":
		if len(*data) == 0 || len(*serviceName) == 0 {
			fmt.Println("Missing data or service name to send")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else {
			base64DecodedSendData, err := base64.StdEncoding.DecodeString(*data)
			if err != nil {
				fmt.Println("Error decoding data: ", err.Error())
				break
			}

			data := xpc.Dict{}
			err = json.Unmarshal(base64DecodedSendData, &data)
			if err != nil {
				fmt.Println("Error in Unmarshal to deserialize data: ", err.Error())
				break
			}

			var m *XpcMan
			if *privileged {
				m = New(*serviceName, 1)
			} else {
				m = New(*serviceName, 0)
			}

			m.logger.Println("Sending data to xpcservice:  ", *serviceName)
			m.conn.Send(data, false)
			break
		}
		break
	case "submit":
		if len(*serviceName) == 0 {
			fmt.Println("missing service name")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else if len(*program) == 0 {
			fmt.Println("missing program")
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		} else {
			var k int
			if *keepalive {
				k = 1
			} else {
				k = 0
			}
			response := xpc.XpcLaunchSubmit(*serviceName, *program, k)
			response = response.(xpc.Dict)
			raw, err := json.MarshalIndent(response, "", "	")
			if err != nil {
				fmt.Println("Error serializing golang object ", err.Error())
			}
			fmt.Printf("%s\n", string(raw))
		}
		break
	default:
		fmt.Println("Missing command")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		break
	}
}
