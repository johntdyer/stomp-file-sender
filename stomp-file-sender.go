package main

import (
	"fmt"
	logrus "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/gmallard/stompngo"
	"gopkg.in/alecthomas/kingpin.v1"
	"net"
	"os"
	"strconv"
	"sync"
)

var (
	queueName     = kingpin.Flag("queue", "Destination").Default("/queue/client_test").Short('q').String()
	serverAddr    = kingpin.Flag("server", "STOMP server endpoint").Default("amq1.prod.us-east-1.aws.tropo.com").Short('s').String()
	serverPort    = kingpin.Flag("port", "STOMP server port").Default("61613").Short('p').String()
	fileToProcess = kingpin.Flag("file", "File to process").Short('f').String()
	workerCount   = kingpin.Flag("workers", "Number of workers to send/receive").Short('w').Int()
	serverUser    = kingpin.Flag("user", "Username").OverrideDefaultFromEnvar("STOMP_USER").String()
	serverPass    = kingpin.Flag("pass", "Password").OverrideDefaultFromEnvar("STOMP_PASS").String()
	debug         = kingpin.Flag("debug", "Enable debug mode.").Short('d').Bool()

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
	kingpin.Version("0.1.1")
	kingpin.Parse()
	logrus.SetOutput(os.Stderr)
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	// Set default of 4 workers
	if *workerCount == 0 {
		*workerCount = 4
	}

}

//  Start main
//
//
func main() {
	producer_wg := &sync.WaitGroup{}
	consumer_wg := &sync.WaitGroup{}

	// Add producer wait group
	producer_wg.Add(1)

	// Add worker wait groups
	consumer_wg.Add(*workerCount)

	logrus.Info("Starting connection")
	_ = client.Connect()
	defer client.Disconnect()

	dataCh := make(chan string)

	// Start workers
	logrus.Debug("Initilizing workers")
	for id := 1; id <= *workerCount; id++ {
		go worker(id, &client, dataCh, consumer_wg)
	}
	logrus.WithFields(logrus.Fields{"Worker Count": *workerCount}).Debug("All workers initialized")

	//  Start reader go routine
	go fileReader(*fileToProcess, dataCh, producer_wg)

	// End wait groups
	producer_wg.Wait()
	consumer_wg.Wait()

	logrus.Info("Done")
}

//  Read from file and put data line by line on channel
func fileReader(path string, dataCh chan<- string, producer_wg *sync.WaitGroup) {
	linesProcessedCounter := 0

	defer producer_wg.Done()

	inFile, err := os.Open(path)
	if err != nil {
		logrus.Fatal(err)
	}
	defer inFile.Close()

	totalLineCount, _ := lineCounter(path)

	logrus.WithFields(logrus.Fields{"Lines": totalLineCount}).Info("Starting to process input file")
	scanner := NewScanner(inFile)

	for scanner.Scan() {
		linesProcessedCounter++
		if linesProcessedCounter%1000 == 0 && *debug {

			logrus.WithFields(logrus.Fields{
				"Completed": fmt.Sprintf("%s", humanize.Ftoa((float64(linesProcessedCounter)/float64(totalLineCount))*100.0)),
				"Finished":  linesProcessedCounter,
				"Total":     totalLineCount,
				"Remaining": (totalLineCount - linesProcessedCounter),
			}).Debug("Progress data")
		}
		dataCh <- scanner.Text()
	}

	logrus.WithFields(logrus.Fields{"Lines": linesProcessedCounter}).Debug("Finished reading " + path)

	close(dataCh)
}

//  Read from channel and put on queue
func worker(id int, client *Client, dataCh <-chan string, w *sync.WaitGroup) {
	defer w.Done()
	for message := range dataCh {
		headers := stompngo.Headers{"destination", *queueName, "suppress-content-length", "true", "id", client.Uuid, "persistent", "true"}
		err := client.StompConnection.Send(headers, message)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	logrus.Debug("Worker # [" + strconv.Itoa(id) + "] done")
}
