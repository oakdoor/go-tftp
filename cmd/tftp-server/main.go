// Copyright (C) 2025 PA Knowledge Ltd.
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package main

import (
	"log"
	"github.com/oakdoor/go-tftp/tftp"
	"io"
	"os"
	"path/filepath"
	"fmt"
	"flag"
)

func Receiver(outputFolder string, atomicSave bool) func(writeRequest tftp.WriteRequest) {
    return func (writeRequest tftp.WriteRequest) {
        var filename = writeRequest.Name()
        var outputFilePath = filepath.Join(outputFolder, filepath.Clean(filename))

        if atomicSave {
            err := writeToFile(dotpath(outputFolder, filename), writeRequest)
            if err != nil {
                log.Println(err)
                writeRequest.WriteError(tftp.ErrCodeAccessViolation, fmt.Sprintf("Cannot write to file %q", outputFilePath))
                return
            }

            err = os.Rename(dotpath(outputFolder, filename), outputFilePath)
            if err != nil {
                log.Println(err)
                writeRequest.WriteError(tftp.ErrCodeAccessViolation, fmt.Sprintf("Cannot rename file to %q", outputFilePath))
            }
        } else {
            err := writeToFile(outputFilePath, writeRequest)

            if err != nil {
                log.Println(err)
                writeRequest.WriteError(tftp.ErrCodeAccessViolation, fmt.Sprintf("Cannot write to file %q", outputFilePath))
                return
            }
        }
    }
}

func dotpath(outputFolder string, filename string) string {
    return filepath.Join(outputFolder, filepath.Join("." + filepath.Clean(filename)))
}

func writeToFile(fullPath string, writeRequest tftp.WriteRequest) error {
        file, err := os.Create(fullPath)

        if err != nil {
            log.Println(err)
            writeRequest.WriteError(tftp.ErrCodeAccessViolation, fmt.Sprintf("Cannot create file %q", fullPath))
            return err
        }

        defer file.Close()

        _, err = io.Copy(file, writeRequest)

        if err != nil {
            log.Println(err)
            os.Remove(fullPath)
            writeRequest.WriteError(tftp.ErrCodeAccessViolation, fmt.Sprintf("Cannot write to file %q", fullPath))
            return err
        }
        return err
}

func main() {
    var singlePortMode = flag.Bool("single-port", false, "When set the server will not use standard ephemeral ports for the TFTP transaction, making firewall configuration easier.")
    var outputFolder = flag.String("output-folder", "output", "The write location of received files.")
    var port = flag.Int("port", 69, "The UDP port the server will listen on.")
    var atomicSave = flag.Bool("atomic-save", false, "When set the server will write the message contents to the output file atomically.")
    flag.Parse()

    opts:= []tftp.ServerOpt{tftp.ServerSinglePort(*singlePortMode)}

	server, err := tftp.NewServer(fmt.Sprintf(":%d", *port), opts...)
	if err != nil {
		log.Fatal(err)
	}

	fs := tftp.WriteHandlerFunc(Receiver(*outputFolder, *atomicSave))

	server.WriteHandler(fs)

	log.Fatal(server.ListenAndServe())
}
