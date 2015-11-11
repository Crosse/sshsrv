package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

const serviceName = "ssh"

var (
	sshPath string
	verbose bool
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

	const usageVerbose = "enable verbose logging"
	flag.BoolVar(&verbose, "v", false, usageVerbose)
	flag.BoolVar(&verbose, "verbose", false, usageVerbose+" (shorthand)")
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	host := flag.Args()[0]
	sshArgs := flag.Args()[1:]
	targetHost, targetPort, err := GetSSHEndpoint(host)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Target for %v is %v:%v", host, targetHost, targetPort)

	//args := []string{sshPath}
	args := []string{}
	args = append(args, fmt.Sprintf("-p %d", targetPort))
	args = append(args, targetHost)
	args = append(args, sshArgs...)
	if verbose {
		args = append(args, "-v")
	}
	cmd := exec.Command(sshPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Wait()
}
