# **github.com/oakdoor/go-tftp**

This is a fork of [pack.ag/tftp](https://github.com/vcabbage/go-tftp), a cross-platform, concurrent TFTP client and server implementation written in Go.

This fork adds CLI client and server applications, and unique quality of life features.

### Standards implemented
- [X] Binary Transfer ([RFC 1350](https://tools.ietf.org/html/rfc1350))
- [X] Netascii Transfer ([RFC 1350](https://tools.ietf.org/html/rfc1350))
- [X] Option Extension ([RFC 2347](https://tools.ietf.org/html/rfc2347))
- [X] Blocksize Option ([RFC 2348](https://tools.ietf.org/html/rfc2348))
- [X] Timeout Interval Option ([RFC 2349](https://tools.ietf.org/html/rfc2349))
- [X] Transfer Size Option ([RFC 2349](https://tools.ietf.org/html/rfc2349))
- [X] Windowsize Option ([RFC 7440](https://tools.ietf.org/html/rfc7440))

### Unique features of this fork
- __Single Port Mode for the TFTP Client__

    TL;DR: It allows TFTP to work through firewalls.

    A standard TFTP server implementation receives requests on port 69 and allocates a new high port (over 1024) dedicated to that request.
    In single port mode, the same port is used for transmit and receive. If the server is started on port 69, all communication will
    be done on port 69.
    
    The primary use case of this feature is to play nicely with firewalls. Most firewalls will prevent the typical case where the server responds
    back on a random port because they have no way of knowing that it is in response to a request that went out on port 69. In single port mode,
    the firewall will see a request go out to a server on port 69 and that server respond back on the same port, which most firewalls will allow.
    
    Of course if the firewall in question is configured to block TFTP connections, this setting won't help you.
    
    Enable single port mode with the `--single-port` flag. This is currently marked experimental as it diverges from the TFTP standard.

## Licenses
This is a fork of [pack.ag/tftp](https://github.com/vcabbage/go-tftp) which is licensed under MIT, retained in [LICENSE](LICENSE).

## Building client and server applications

### Dependencies
Install Go following these [instructions](https://go.dev/doc/install).

### Build commands
```bash
go build -o tftp-client cmd/tftp-client/main.go
```
```bash
go build -o tftp-server cmd/tftp-server/main.go
```

### Usage

```bash
./tftp-client --help
```
```bash
./tftp-server --help
```

### Examples

```bash
./tftp-server --output-folder output/ --port 69
```
```bash
./tftp-client --file test_file --windowsize 64 --blocksize 1408 tftp://0.0.0.0/test_file
```

## The tftp package

### Installing

```bash
go get -u github.com/oakdoor/go-tftp/tftp
```

### API

The API was inspired by Go's well-known net/http API. If you can write a net/http handler or middleware, you should have no problem doing the same with the tftp package.

#### Configuration functions

One area that is noticeably different from net/http is the configuration of clients and servers. This tftp package uses "configuration functions" rather than the direct modification of the
Client/Server struct or a configuration struct passed into the factory functions.

A few explanations of this pattern:
* [Self-referential functions and the design of options](http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html) by Rob Pike
* [Functional options for friendly APIs](https://www.youtube.com/watch?v=24lFtGHWxAQ) by Dave Cheney [video]

If this sounds complicated, don't worry, the public API is quiet simple. The `NewClient` and `NewServer` functions take zero or more configuration functions.

Want all defaults? Don't pass anything.

Want a Client configured for blocksize 9000 and windowsize 16? Pass in `ClientBlocksize(9000)` and `ClientWindowsize(16)`.

```go
// Default Client
tftp.NewClient()

// Client with blocksize 9000, windowsize 16
tftp.NewClient(tftp.ClientBlocksize(9000), tftp.ClientWindowsize(16))

// Configuring with a slice of options
opts := []tftp.ClientOpt{
    tftp.ClientMode(tftp.ModeOctet),
    tftp.ClientBlocksize(9000),
    tftp.ClientWindowsize(16),
    tftp.ClientTimeout(1),
    tftp.ClientTransferSize(true),
    tftp.ClientRetransmit(3),
}

tftp.NewClient(opts...)
```

#### Examples

##### Read file from server, print to stdout

```go
client := tftp.NewClient()
resp, err := client.Get("myftp.local/myfile")
if err != nil {
    log.Fatalln(err)
}

err := io.Copy(os.Stdout, resp)
if err != nil {
    log.Fatalln(err)
}
```

##### Write file to server

```go

file, err := os.Open("myfile")
if err != nil {
    log.Fatalln(err)
}
defer file.Close()

// Get the file info se we can send size (not required)
fileInfo, err := file.Stat()
if err != nil {
    log.Println("error getting file size:", err)
}

client := tftp.NewClient()
err := client.Put("myftp.local/myfile", file, fileInfo.Size())
if err != nil {
    log.Fatalln(err)
}
```


##### Other examples

Full examples including HTTP proxy and database access can be found in [pack.ag/tftp](https://github.com/vcabbage/go-tftp).
