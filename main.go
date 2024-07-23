package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	mode := "TCP+TLS"
	var client *modbus.ModbusClient
	var err error

	if mode == "TCP+TLS" {
		//clientCert, err = GenerateSelfSignedCertificate("Controller")
		clientCert, err := tls.LoadX509KeyPair("HomeIoClientTLS.crt", "key.pem")
		if err != nil {
			fmt.Printf("failed to load client key pair: %v\n", err)
			os.Exit(1)
		}
		serverCertPool, err := modbus.LoadCertPool("CA-OT-Security.crt")
		if err != nil {
			fmt.Printf("failed to load server certificate/CA: %v\n", err)
			os.Exit(1)
		}

		client, err = modbus.NewClient(&modbus.ClientConfiguration{
			URL:           "tcp+tls://127.0.0.1:5802",
			TLSClientCert: &clientCert,
			TLSRootCAs:    serverCertPool,
		})
	} else {
		client, err = modbus.NewClient(&modbus.ClientConfiguration{
			URL: "tcp://127.0.0.1:1502",
		})
	}

	client.SetUnitId(5)

	if err != nil {
		fmt.Printf("failed to create modbus client: %v\n", err)
		os.Exit(1)
	}

	err = client.Open()
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		os.Exit(2)
	}

	for {
		time.Sleep(250 * time.Millisecond)
		sensor, err := client.ReadDiscreteInput(15)
		if err != nil {
			fmt.Printf("failed to read Discret Input: %v\n", err)
		}

		door_garage, err := client.ReadDiscreteInput(14)
		if err != nil {
			fmt.Printf("failed to read Discret Input: %v\n", err)
		}

		door_outside, err := client.ReadDiscreteInput(13)
		if err != nil {
			fmt.Printf("failed to read Discret Input: %v\n", err)
		}

		alarm := sensor || !door_garage || !door_outside
		if alarm {
			fmt.Printf("Alarm !\n")
			err = client.WriteCoil(4, true)
			if err != nil {
				fmt.Printf("failed to write Coil: %v\n", err)
			}
		} else {
			fmt.Printf("Ok\n")
			err = client.WriteCoil(4, false)
			if err != nil {
				fmt.Printf("failed to write Coil: %v\n", err)
			}
		}
	}

	err = client.Close()
	if err != nil {
		fmt.Printf("failed to close connection: %v\n", err)
	}

	os.Exit(0)
}
