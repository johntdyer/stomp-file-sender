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
  --version           Show application version.
```

### Building

    go get github.com/nitrous-io/goop
    goop install
    go build stomp-file-sender.go


### Using


    ./stomp-file-sender --user foo --pass bar --file example.log --server amq1.prod.somewhere.com



## Change log

* 0.1.1

 * Progress percentage indicator

* 0.1.0

 * Support for long lines
 * Improved logging

