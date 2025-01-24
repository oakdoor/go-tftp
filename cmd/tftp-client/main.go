// Copyright (C) 2025 PA Knowledge Ltd.
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package main

import (
	"log"
	"github.com/oakdoor/go-tftp/tftp"
	"os"
	"flag"
	"io"
)

func main() {
    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
    address, options, inputfile := parseCmdLine()

    reader, err := getFileOrStdin(inputfile)
    if err != nil {
    	log.Fatalln(err)
    	os.Exit(1)
    }

    err = sendFile(address, options, reader)

    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }

    os.Exit(0)
}

func parseCmdLine() (string, []tftp.ClientOpt, string) {
    inputfile := flag.String("file", "", "File to send. If not specified, stdin is read instead.")
    windowsize := flag.Int("windowsize", 64, "TFTP windowsize parameter.")
    blocksize := flag.Int("blocksize", 1408, "TFTP blocksize parameter.")
    retransmit := flag.Int("retransmit", 3, "TFTP retransmit parameter.")
    timeout := flag.Int("timeout", 1, "TFTP timeout parameter.")
    singleport := flag.Int("single-port", 0, "The client will use the specified value as the UDP src port for the TFTP transaction, making firewall configuration easier. If not specified or 0, standard TFTP ephemeral ports are used instead.")
    flag.Parse()

    if flag.NArg() != 1 {
        log.Println("USAGE: ")
        log.Println(os.Args[0], "--file test_file [--windowsize [64]] [--blocksize [1408]] [--single-port [0]] [--retransmit [3]] [--timeout [1]] tftp://0.0.0.0/test_file")
        log.Println("or")
        log.Println("echo abc |", os.Args[0], "[--windowsize [64]] [--blocksize [1408]] [--single-port [0]] [--retransmit [3]] [--timeout [1]] tftp://0.0.0.0/test_file")
        log.Println()
        os.Exit(1)
    }

    var address = flag.Args()[0]
    var options = []tftp.ClientOpt{tftp.ClientBlocksize(*blocksize), tftp.ClientWindowsize(*windowsize), tftp.ClientRetransmit(*retransmit), tftp.ClientTimeout(*timeout), tftp.ClientListenPort(*singleport)}
    return address, options, *inputfile
}

func sendFile(address string, options []tftp.ClientOpt, input io.Reader) error {

    client,err := tftp.NewClient(options...)
    if err != nil {
        log.Fatalln(err)
        return err
    }

    return client.Put(address, input, 0)
}

func getFileOrStdin(inputfile string) (io.Reader, error) {
    if inputfile == "" {
        return os.Stdin, nil
    }

    return os.Open(inputfile)
}
