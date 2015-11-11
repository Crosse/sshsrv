package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	log "bitbucket.org/crosse3/gosimplelogger"
)

const serviceName = "ssh"

var (
	sshPath string
	whatif  bool
)

func GetSSHEndpoint(hostname string) (target string, port uint16, err error) {
	cname, addrs, err := net.LookupSRV(serviceName, "tcp", hostname)
	if err != nil {
		return
	}
	log.Printf("Retrieved record for %v", cname)
	log.Printf("Found %d records", len(addrs))

	// "The returned records are sorted by priority and randomized
	// by weight within a priority", so return details for the first
	// one in the list.

	// The target DNS names are fully-specified with the root ("."),
	// so trim that off.
	target = strings.TrimRight(addrs[0].Target, ".")
	port = addrs[0].Port

	return
}

func init() {
	var err error
	sshPath, err = exec.LookPath("ssh")
	if err != nil {
		log.Fatal("Could not find ssh!")
	}

	var verbose = flag.Bool("v", false, "enable verbose logging")
	flag.BoolVar(&whatif, "whatif", false, "perform everything but the actual connection")
	flag.Parse()

	if *verbose {
		log.LogLevel = log.LogVerbose
	}
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	host := flag.Args()[0]
	sshArgs := flag.Args()[1:]
	targetHost, targetPort, err := GetSSHEndpoint(host)
	if err != nil {
		log.Fatal(err)
	}

	log.Verbosef("Target for %v is %v:%v", host, targetHost, targetPort)

	if whatif {
		return
	}

	//args := []string{sshPath}
	args := []string{}
	args = append(args, fmt.Sprintf("-p %d", targetPort))
	args = append(args, targetHost)
	args = append(args, sshArgs...)
	cmd := exec.Command(sshPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Verbosef("Connecting to %v:%v", targetHost, targetPort)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Wait()
}
