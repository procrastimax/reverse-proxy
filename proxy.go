package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	sourceIPAddress      string = "localhost"
	destinationIPAddress string = "localhost"
)

var (
	supportedNetworkTypes = [...]string{"tcp"}
)

// Port represents a simple network port
type Port struct {
	Network string `yaml:"networkType"`
	Port    int    `yaml:"portNumber"`
}

func (p *Port) String() string {
	return fmt.Sprintf("%s:%d", p.Network, p.Port)
}

// Validate checks the correctness of the provided value for the port attributes.
// This functions returns an error if the port could not be validated and nil if everything is ok.
func (p *Port) Validate() error {

	isNetworkTypeOK := false
	// check if the specified network type is currently supported
	// by checking the array with supported network types
	for _, val := range supportedNetworkTypes {
		if strings.ToLower(p.Network) == val {
			isNetworkTypeOK = true
		}
	}

	if isNetworkTypeOK == false {
		return fmt.Errorf("unsupported network type")
	}

	// check the ports > 0 and <= 65535 (aka. 16bit int)
	if p.Port <= 0 && p.Port > 65535 {
		return fmt.Errorf("unvalid port address")
	}
	return nil
}

// ProxyConn represents a proxy connection that connects an ingoing connection to an outgoing one.
type ProxyConn struct {
	listener net.Listener

	dstPort Port
}

// SetupListenerPort creates a listener to accept a connection on the given port.
// This acts as the sourcePort.
func (p *ProxyConn) SetupListenerPort(srcPort Port) error {
	err := srcPort.Validate()
	if err != nil {
		return fmt.Errorf("Validate: %s", err)
	}

	p.listener, err = net.Listen(srcPort.Network, net.JoinHostPort(sourceIPAddress, strconv.Itoa(srcPort.Port)))
	if err != nil {
		return fmt.Errorf("Validate: %s", err)
	}
	log.Printf("created new listener conn: %s\n", srcPort.String())
	return nil
}

// SetupDestinationPort initializes the values used to dial to the destination port.
func (p *ProxyConn) SetupDestinationPort(dstPort Port) error {
	err := dstPort.Validate()
	if err != nil {
		return fmt.Errorf("Validate: %s", err)
	}
	p.dstPort.Port = dstPort.Port
	p.dstPort.Network = dstPort.Network
	return nil
}

// StartListener starts to listen on the port of the specified listener of this proxy.
// All incoming messages on this port are forwarded to the specific destination port on the local machine.
// This function is blocking.
func (p *ProxyConn) StartListener() error {
	if p.listener == nil {
		return fmt.Errorf("listener has not been initialized")
	}

	for {
		lConn, err := p.listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			// close the lConn after every transmission
			defer lConn.Close()
			err := p.newPipe("tcp", 4321, lConn)
			if err != nil {
				log.Printf("newPipe: %s\n", err)
			}
		}()
	}
}

func (p *ProxyConn) newPipe(network string, port int, lConn net.Conn) error {
	if len(p.dstPort.Network) == 0 && p.dstPort.Port == 0 {
		return fmt.Errorf("destination network or port not specified")
	}
	rConn, err := net.Dial(p.dstPort.Network, net.JoinHostPort(destinationIPAddress, strconv.Itoa(p.dstPort.Port)))
	if err != nil {
		return err
	}

	defer rConn.Close()

	finnishChan := make(chan interface{})

	// TODO: self implement copy, to check the chan before reading/ writing
	go func(finnishChan chan interface{}) {
		if _, err = io.Copy(lConn, rConn); err != nil {
			log.Println(err)
		}
		finnishChan <- nil
	}(finnishChan)

	go func(finnishChan chan interface{}) {
		if _, err = io.Copy(rConn, lConn); err != nil {
			log.Println(err)
		}
		finnishChan <- nil
	}(finnishChan)

	// wait here until one of the functions above end with an EOF or an error
	<-finnishChan

	return nil
}

// Close tries to close both the incomming and the outgoing connection
func (p *ProxyConn) Close() error {
	var err error
	// definetely try to close both, if one throws an error, return that error
	err = p.listener.Close()
	return err
}

func (p *ProxyConn) String() string {
	return fmt.Sprintf("%s %s ==> %s", p.dstPort.Network, p.listener.Addr().String(), net.JoinHostPort("localhost", strconv.Itoa(p.dstPort.Port)))
}
