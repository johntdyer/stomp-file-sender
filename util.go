package main

import (
	"bytes"
	logrus "github.com/Sirupsen/logrus"
	"github.com/gmallard/stompngo"
	"io"
	"net"
	"os"
)

func lineCounter(path string) (int, error) {
	inFile, err := os.Open(path)
	if err != nil {
		logrus.Fatal(err)
	}

	defer inFile.Close()
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

// Setups connection options
func (client *Client) setOpts() {

	client.Host = *serverAddr
	client.Port = *serverPort
	client.Uuid = stompngo.Uuid()
	client.Queue = *queueName

	if *serverUser != "" {
		client.User = *serverUser
	}

	if *serverPass != "" {
		client.Password = *serverPass
	}
}

// Creates net connection
func (client *Client) netConnection() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", net.JoinHostPort(client.Host, client.Port))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	client.NetConnection = conn
	return
}

func (client *Client) stompConnection() *stompngo.Connection {
	headers := stompngo.Headers{
		"accept-version", "1.1",
		"host", client.Host,
		"login", client.User,
		"passcode", client.Password,
	}

	conn, err := stompngo.Connect(client.NetConnection, headers)
	if err != nil {
		logrus.Fatal(err)

	}

	client.StompConnection = conn
	return conn
}

func (client *Client) Connect() (conn *stompngo.Connection) {
	client.setOpts()
	client.netConnection()

	conn = client.stompConnection()

	return
}

func (client *Client) Disconnect() {
	client.StompConnection.Disconnect(stompngo.Headers{})
	client.NetConnection.Close()
}
