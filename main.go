package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	log "github.com/crosse/gosimplelogger"
)

const (
	serviceName = "ssh"
	defaultPort = 22
)

// GetSSHEndpoint tries to determine how to connect to a particular host
// via SSH.  GetSSHEndpoint first attempts to discover the endpoint via
// DNS SRV records of the form "_ssh._tcp.<hostname>".  If found,
// GetSSHEndpoint will return the target host and port to connect to
// instead of the "bare" DNS hostname.  If no SRV record is found, it
// will simply return the hostname and default SSH port (22).
func GetSSHEndpoint(hostname string) (target string, port uint16, err error) {
	cname, srvAddrs, err := net.LookupSRV(serviceName, "tcp", hostname)
	if err != nil {
		if _, ok := err.(*net.DNSError); ok {
			// DNS-related error. This is okay for now;
			// we'll just fall back to trying the hostname
			// as-is instead.
			log.Verbosef("error: %v", err)
			err = nil
		} else {
			// Non-DNS error.  Probably want to stop now.
			log.Fatal(err)
		}
	}

	if len(srvAddrs) > 0 {
		log.Verbosef("Found %d SRV record(s) for %v", len(srvAddrs), cname)

		for i, r := range srvAddrs {
			log.Verbosef("Record %d:\t%d %d %d %s", i, r.Priority, r.Weight, r.Port, r.Target)
		}

		// "The returned records are sorted by priority and randomized
		// by weight within a priority", so return details for the first
		// one in the list.

		// The target DNS names are fully-specified with the root ("."),
		// so trim that off.
		target = strings.TrimRight(srvAddrs[0].Target, ".")
		port = srvAddrs[0].Port
	} else {
		log.Verbosef("No SRV record found for %v", hostname)
		target = hostname
		port = defaultPort
	}

	// TODO: Extend this to return the entire list in priority-order
	return
}

func usage() {
	fmt.Println("Usage: check the ssh man page.")
}

// parseArgs parses the command line arguments to find and return:
// a) arguments to the ssh(1) command;
// b) "user@host";
// c) the command to run on the remote server, if any.
//
// Doing it this way just seemed easier than using the "flags"
// package and creating flag variables for each of the 44 options,
// because honestly, we don't care about any of them.
func parseArgs(args []string) (sshArgs []string, host string, sshCommand []string) {
	const (
		// Taken from the ssh(1) man page.  These two will need
		// to be updated if/when other options are added to
		// ssh(1).
		sshBoolArgs  = "1246AaCfGgKkMNnqsTtVvXxYy"
		sshParamArgs = "bcDEeFIiLlmOopQRSWw"
	)

	var argPos int

	for argPos = 1; argPos < len(os.Args); argPos++ {
		arg := os.Args[argPos]
		if len(arg) > 1 && arg[0] == '-' {
			if strings.Contains(sshBoolArgs, string(arg[1])) {
				sshArgs = append(sshArgs, arg)
				if arg[1] == 'v' {
					log.LogLevel = log.LogVerbose
				}
			} else if strings.Contains(sshParamArgs, string(arg[1])) {
				sshArgs = append(sshArgs, fmt.Sprintf("%v %v", arg, os.Args[argPos+1]))
				argPos++
			}
		} else {
			if strings.Contains(arg, "@") {
				x := strings.Split(arg, "@")
				sshArgs = append(sshArgs, fmt.Sprintf("-l %v", x[0]))
				host = x[1]
			} else {
				host = arg
			}
			argPos++
			break
		}
	}
	if len(os.Args[argPos:]) > 0 {
		sshCommand = os.Args[argPos:]
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		log.Fatal("Could not find ssh!")
	}

	sshArgs, host, sshCommand := parseArgs(os.Args[1:])

	targetHost, targetPort, err := GetSSHEndpoint(host)
	if err != nil {
		log.Fatal(err)
	}

	log.Verbosef("Target for %v is %v:%v", host, targetHost, targetPort)

	args := []string{}
	args = append(args, fmt.Sprintf("-p %d", targetPort))
	args = append(args, sshArgs...)
	args = append(args, targetHost)
	args = append(args, sshCommand...)
	log.Verbosef("command: %v %v", sshPath, strings.Join(args, " "))

	cmd := exec.Command(sshPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Verbosef("Connecting to %v:%v", targetHost, targetPort)
	if err = cmd.Run(); err != nil {
		if err, ok := err.(*exec.ExitError); !ok {
			log.Fatal(err)
		}
	}
}
