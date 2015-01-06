package main

import (
	"bufio"
	"fmt"
	"github.com/gmallard/stompngo"
	"gopkg.in/alecthomas/kingpin.v1"
	"net"
	"os"
)

var (
	queueName     = kingpin.Flag("queue", "Destination").Default("/queue/client_test").Short('q').String()
	serverAddr    = kingpin.Flag("server", "STOMP server endpoint").Default("amq1.prod.us-east-1.aws.tropo.com").Short('s').String()
	serverPort    = kingpin.Flag("port", "STOMP server port").Default("61613").Short('p').String()
	fileToProcess = kingpin.Flag("file", "File to process").Short('f').String()
	workerCount   = kingpin.Flag("workers", "Number of workers to send/receive").Short('w').Int()
	serverUser    = kingpin.Flag("user", "Username").OverrideDefaultFromEnvar("STOMP_USER").String()
	serverPass    = kingpin.Flag("pass", "Password").OverrideDefaultFromEnvar("STOMP_PASS").String()

	client Client
	done   = make(chan bool)
)

//var done = make(chan bool)

type Client struct {
	Host     string
	Port     string
	User     string
	Password string
	Uuid     string
	Queue    string

	NetConnection   net.Conn
	StompConnection *stompngo.Connection
}

func init() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	// Set default of 4 workers
	if *workerCount == 0 {
		*workerCount = 4
	}

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
		fmt.Println(err)
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
		fmt.Println(err)
		os.Exit(0)
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

//  Start main
//
//
func main() {

	fmt.Println("Starting connection")
	_ = client.Connect()
	defer client.Disconnect()

	dataCh := make(chan string, *workerCount)

	// Start workers
	fmt.Println("Create workers")
	for id := 1; id <= *workerCount; id++ {
		go sender(id, &client, dataCh)
	}

	//  Start reader go routine
	fmt.Println("Read file")
	go fileReader(*fileToProcess, dataCh)

	<-done

	fmt.Println("Done")
}

//  Read from file and put data line by line on channel
func fileReader(path string, dataCh chan<- string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		dataCh <- scanner.Text()
	}
	close(dataCh)
}

//  Read from channel and put on queue
func sender(id int, client *Client, dataCh <-chan string) {

	for message := range dataCh {
		//fmt.Println("Worker", id, "message", len(message))

		headers := stompngo.Headers{"destination", *queueName, "suppress-content-length", "true", "id", client.Uuid, "persistent", "true"}
		err := client.StompConnection.Send(headers, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	done <- true
}
