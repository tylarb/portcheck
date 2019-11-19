/*
Released under MIT license, copyright 2019 Tyler Ramer
*/

package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "portcheck [flags] [SERVER:PORT]...",
	Short: "portcheck",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		ports, err = parsePorts(inputPorts)
		if err != nil {
			return err
		}
		if !tcp && !udp {
			return fmt.Errorf("At least one of tcp (-t) or udp (-u) flag must be specified")
		}

		return runPortTest(args)
	},
}

var verbose bool

var inputPorts string
var tcp bool
var udp bool
var ports []int

// Size of the payload, in bytes.
var payloadSize int

var errPortFmt = fmt.Errorf("Invalid port range, must be list of comma separated ranges like 1,2-4,10-10")

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputPorts, "ports", "p", "", "list of ports to check, like 1,2-4,10")
	rootCmd.PersistentFlags().BoolVarP(&tcp, "tcp", "t", false, "check tcp ports")
	rootCmd.PersistentFlags().BoolVarP(&udp, "udp", "u", false, "check udp ports")
	rootCmd.PersistentFlags().IntVarP(&payloadSize, "payloadsize", "s", 1300, "Size of the TCP/UDP payload, useful for testing MTU related issues")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Debug logs")
}

func main() {
	rootCmd.Execute()
}

func runPortTest(hosts []string) error {
	if tcp {
		for _, host := range hosts {
			addr, err := net.ResolveTCPAddr("tcp", host)
			if err != nil {
				fmt.Printf("Unable to resolve TCP host %s, skipping\n", host)
			}
			_ = tcpPortTest(addr)
		}
	}

	if udp {
		for _, host := range hosts {
			addr, err := net.ResolveUDPAddr("udp", host)
			if err != nil {
				fmt.Printf("Unable to resolve UDP host %s, skipping\n", host)
			}
			_ = udpPortTest(addr)
		}
	}

	return nil
}

func tcpPortTest(addr *net.TCPAddr) error {
	success := false
	for _, port := range ports {
		local := net.TCPAddr{Port: port}
		c, err := net.DialTCP("tcp", &local, addr)
		if err != nil {
			fmt.Printf("FAILED: local TCP port %d failed to access %s\n", port, addr.String())
			continue
		}

		header := []byte(fmt.Sprintf("Message sent from %s\n", local.String()))
		footer := []byte("end of TCP message")
		b := composeByteMessage(header, footer, payloadSize)

		length, err := c.Write(b)
		if err != nil {
			fmt.Printf("error writing to connection: %q\n", err)
		}
		if verbose {
			fmt.Printf("DEBUG: length of packets written: %d\n", length)
		}
		c.Close()
		success = true
	}

	if !success {
		fmt.Printf("Failed to connect to host %s via any local port - is the server listening?\n", addr.String())
	}

	return nil
}

func udpPortTest(addr *net.UDPAddr) error {
	success := false
	for _, port := range ports {
		local := net.UDPAddr{Port: port}
		c, err := net.DialUDP("udp", &local, addr)
		if err != nil {
			fmt.Printf("FAILED: local UDP port %d failed to access %s\n", port, addr.String())
			continue
		}

		header := []byte(fmt.Sprintf("Message sent from %s\n", local.String()))
		footer := []byte("end of UDP message")
		b := composeByteMessage(header, footer, payloadSize)
		length, err := c.Write(b)
		if err != nil {
			fmt.Printf("error writing to connection: %q\n", err)
		}
		if verbose {
			fmt.Printf("DEBUG: length of packets written: %d\n", length)
		}

		c.Close()
		success = true
	}

	if !success {
		fmt.Printf("Failed to connect to host %s via any local port - is the server listening?\n", addr.String())
	}

	return nil
}

// composes a byte message of len(size) leading with header bytes and trailing with footer bytes
// If size < header or footer, the message is trimmed.
// If if size < header + footer, the message is trimmed, and the trailing bytes of header
// overwrites the leading bytes of footer
func composeByteMessage(header, footer []byte, size int) []byte {
	b := make([]byte, size)
	copy(b[len(b)-len(footer):], footer[:])
	copy(b[:len(header)], header[:])
	return b
}

func parsePorts(p string) ([]int, error) {
	commaSep := strings.Split(p, ",")
	var portsInt []int
	for _, s := range commaSep {
		if !strings.Contains(s, "-") {
			i, err := strconv.Atoi(s)
			if err != nil {
				return nil, errPortFmt
			}
			portsInt = append(portsInt, i)
			continue
		}
		portRange := strings.Split(s, "-")
		if len(portRange) != 2 {
			return nil, errPortFmt
		}
		lower, err := strconv.Atoi(portRange[0])
		if err != nil {
			return nil, errPortFmt
		}
		upper, err := strconv.Atoi(portRange[1])
		if err != nil {
			return nil, errPortFmt
		}
		if lower > upper {
			return nil, errPortFmt
		}
		for i := lower; i <= upper; i++ {
			portsInt = append(portsInt, i)
		}
	}
	return portsInt, nil
}
