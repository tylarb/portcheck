/*
Released under MIT license, copyright 2019 Tyler Ramer
*/

package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "portcheck-server",
	Short: "portcheck-server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if udp && tcp {
			return fmt.Errorf("Specify exactly one of tcp (-t) or udp (-u)")
		}
		if !tcp && !udp {
			return fmt.Errorf("Exactly one of tcp (-t) or udp (-u) flag must be specified")
		}
		if listenPort < 1 || listenPort > 65535 {
			return fmt.Errorf("Please specify a listning port between 1-65535")
		}

		return startServer()
	},
}

var listenPort int
var tcp bool
var udp bool
var listenType string
var bufferSize int
var verbose bool

func init() {
	rootCmd.PersistentFlags().IntVarP(&listenPort, "port", "p", 0, "port for server to listen on")
	rootCmd.PersistentFlags().BoolVarP(&tcp, "tcp", "t", false, "check tcp ports")
	rootCmd.PersistentFlags().BoolVarP(&udp, "udp", "u", false, "check udp ports")
	rootCmd.PersistentFlags().IntVarP(&bufferSize, "buffersize", "s", 1500, "size of read buffer, recommened to set to MTU")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Debug level logs")
}

func startServer() error {
	if tcp {
		return startTCPServer()
	}
	return startUDPServer()
}

func startTCPServer() error {

	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(listenPort))
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %q\n", err)
		}
		go func(c net.Conn) {
			addr := c.RemoteAddr()
			b := make([]byte, bufferSize)
			length, err := c.Read(b)
			if err != nil {
				fmt.Printf("error on packet read: %q", err)
			}

			fmt.Printf("TCP packet received from %s size %d bytes\n", addr.String(), length)
			if verbose {
				fmt.Printf("message: %s\n", string(b))
			}

			c.Close()
		}(conn)
	}
}

func startUDPServer() error {
	l, err := net.ListenPacket("udp", "0.0.0.0:"+strconv.Itoa(listenPort))
	if err != nil {
		return err
	}
	defer l.Close()
	doneChan := make(chan error, 1)
	go func() {
		for {
			b := make([]byte, bufferSize)
			length, addr, err := l.ReadFrom(b)
			if err != nil {
				fmt.Printf("error on packet read: %q", err)
			}
			fmt.Printf("UDP packet received from %s size %d bytes\n", addr.String(), length)
			if verbose {
				fmt.Printf("message: %s\n", string(b))
			}
		}
	}()
	select {
	case err = <-doneChan:
	}
	return err
}

func main() {
	rootCmd.Execute()
}
