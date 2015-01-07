# Stomp File Sender

Small utility which reads a file line by line and sends to a queue via Stomp.  This was mainly an excercise to use channels for the first time.
## Usage

```
usage: stomp-file-sender [<flags>]

Flags:
  --help              Show help.
  -q, --queue="/queue/client_test"
                      Destination
  -s, --server="amq1.prod.somewhere.com"
                      STOMP server endpoint
  -p, --port="61613"  STOMP server port
  -f, --file=FILE     File to process
  -w, --workers=WORKERS ( Default 4 )
                      Number of workers to send/receive
  --user=USER         Stomp username
  --pass=PASS         Stomp password
  -d, --debug         Debug mode
  --formatter="text"  Formatter (text or json)
  --version           Show application version.
```

### Building

    go get github.com/nitrous-io/goop
    goop install
    go build stomp-file-sender.go


### Using


    ./stomp-file-sender --user foo --pass bar --file example.log --server amq1.prod.somewhere.com



## Change log
* 0.2.0

* Support JSON or TEXT logging

* 0.1.2

 * Progress percentage indicator as whole number

```log
INFO[0000] Starting connection
DEBU[0000] Initilizing workers
DEBU[0000] All workers initialized                       Worker Count=5
INFO[0000] Starting to process input file                Lines=597283
DEBU[0002] Progress data                                 Completed=0% LinesProcessed=1000 LinesRemaining=596283 LinesTotal=597283
DEBU[0005] Progress data                                 Completed=0% LinesProcessed=2000 LinesRemaining=595283 LinesTotal=597283
DEBU[0008] Progress data                                 Completed=0% LinesProcessed=3000 LinesRemaining=594283 LinesTotal=597283
DEBU[0011] Progress data                                 Completed=0% LinesProcessed=4000 LinesRemaining=593283 LinesTotal=597283
DEBU[0013] Progress data                                 Completed=0% LinesProcessed=5000 LinesRemaining=592283 LinesTotal=597283
```

* 0.1.1

 * Progress percentage indicator

* 0.1.0

 * Support for long lines
 * Improved logging

