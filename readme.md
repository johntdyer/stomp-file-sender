# Stomp File Sender

Small utility which reads a file line by line and sends to a queue via Stomp.  This was mainly an excercise to use channels for the first time.
## Usage

```
usage: stomp-file-sender [<flags>]

Flags:
  --help              Show help.
  -q, --queue="/queue/client_test"
                      Destination
  -s, --server="amq1.prod.us-east-1.aws.tropo.com"
                      STOMP server endpoint
  -p, --port="61613"  STOMP server port
  -f, --file=FILE     File to process
  -w, --workers=WORKERS
                      Number of workers to send/receive
  --user=USER         Username
  --pass=PASS         Password
  --version           Show application version.
```

### Building

    go get github.com/nitrous-io/goop
    goop install
    go build stomp-file-sender.go


### Using


    ./stomp-file-sender --user foo --pass bar --file example.log --server amq1.prod.somewhere.com


