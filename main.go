package main

import (
	"log"
)

const (
	configPath = "config.yaml"
)

func main() {

	proxyConfig, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %s\n", err)
	}

	// construct needed ports based on config

	proxyConnList := make([]ProxyConn, 0)

	for _, conn := range proxyConfig.Connections {
		proxyConn := ProxyConn{}

		err := proxyConn.SetupListenerPort(conn.SourcePort)
		if err != nil {
			log.Printf("SetupListener: %s\n", err)
			// don't add unitialized proxyConns to list
			continue
		}

		err = proxyConn.SetupDestinationPort(conn.DstPort)
		if err != nil {
			log.Printf("SetupDestinationPort: %s\n", err)
			// don't add unitialized proxyConns to list
			continue
		}

		// add new proxy port to list
		proxyConnList = append(proxyConnList, proxyConn)
	}

	log.Println("Setup proxy connections...")
	log.Println("Starting to listen on proxy connections...")
	loggerChan := make(chan string)

	// starting to listen on all proxyPorts
	for _, proxyConn := range proxyConnList {
		// each proxyPort creates a new goroutine
		go func(proxyConn ProxyConn) {
			err := proxyConn.StartListener()
			if err != nil {
				loggerChan <- err.Error()
				proxyConn.Close()
			}
		}(proxyConn)
	}

	for {
		// endlessly log errors to stdout
		log.Printf("ERROR: %s\n", <-loggerChan)
	}
}
