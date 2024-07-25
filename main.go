package main

import (
	"crypto/tls"
	"fmt"
	"github.com/simonvetter/modbus"
	"os"
	"time"
)

var mode = "TCP+TLS"
var ip = `127.0.0.1`

func checkAlarm(client *modbus.ModbusClient) {
	// Read the motion sensor status
	sensor, err := client.ReadDiscreteInput(15)
	if err != nil {
		fmt.Printf("failed to read Discret Input: %v\n", err)
	}

	// Read the garage door status
	doorGarage, err := client.ReadDiscreteInput(14)
	if err != nil {
		fmt.Printf("failed to read Discret Input: %v\n", err)
	}

	// Read the outside door status
	doorOutside, err := client.ReadDiscreteInput(13)
	if err != nil {
		fmt.Printf("failed to read Discret Input: %v\n", err)
	}

	// Check the alarm status
	alarm := sensor || !doorGarage || !doorOutside

	if alarm {
		// Set the alarm
		fmt.Printf("Alarm !\n")
		err = client.WriteCoil(4, true)
		if err != nil {
			fmt.Printf("failed to write Coil: %v\n", err)
		}
	} else {
		// Reset the alarm
		fmt.Printf("Ok\n")
		err = client.WriteCoil(4, false)
		if err != nil {
			fmt.Printf("failed to write Coil: %v\n", err)
		}
	}
}

func main() {
	var client *modbus.ModbusClient
	var err error

	// Start the ModBus client / master connection on TCP or TCP+TLS
	if mode == "TCP+TLS" {

		// Load the client certificate with the private key
		clientCert, err := tls.LoadX509KeyPair("HomeIoClientTLS.crt", "key.pem")
		if err != nil {
			fmt.Printf("failed to load client key pair: %v\n", err)
			os.Exit(1)
		}

		// Load the server certificate authority
		serverCertPool, err := modbus.LoadCertPool("CA-OT-Security.crt")
		if err != nil {
			fmt.Printf("failed to load server certificate/CA: %v\n", err)
			os.Exit(1)
		}

		// Create the Modbus TCP+TLS client instance
		client, err = modbus.NewClient(&modbus.ClientConfiguration{
			URL:           "tcp+tls://" + ip + ":5802",
			TLSClientCert: &clientCert,
			TLSRootCAs:    serverCertPool,
		})
		if err != nil {
			fmt.Printf("failed to start modbus TCP+TLS instance: %v\n", err)
			os.Exit(1)
		}

	} else {

		// Create the Modbus TCP client instance
		client, err = modbus.NewClient(&modbus.ClientConfiguration{
			URL: "tcp://" + ip + ":1502",
		})
		if err != nil {
			fmt.Printf("failed to start modbus TCP instance: %v\n", err)
			os.Exit(1)
		}
	}

	// Set the Modbus client unit ID
	err = client.SetUnitId(5)
	if err != nil {
		fmt.Printf("failed to set Unit Id: %v\n", err)
		os.Exit(2)
	}

	// Open the Modbus client connection
	err = client.Open()
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		os.Exit(3)
	}

	// Loop to check the alarm status
	for {
		checkAlarm(client)
		time.Sleep(250 * time.Millisecond)
	}

	/*
		err = client.Close()
		if err != nil {
			fmt.Printf("failed to close connection: %v\n", err)
		}

		os.Exit(0)
	*/
}
