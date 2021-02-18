package format

import (
	"fmt"
	"storage-go"
)

// custom1
//  header: no
//  contents: CDR1CDR2CDR3 ...CDRn
//  eof: \n
//  name: <node-id-suffix+vpn-id>_<date>+<time>_<total-cdrs>_file<fileseqnum>

// Custom1 holds necessary functions to build Custom1 file format
type Custom1 struct{}

// NewCustom1Format instantiates custom1 file format
func NewCustom1Format() *Custom1 {
	return &Custom1{}
}

// Header writes no header for Custom1 format
func (c1 *Custom1) Header(t *storage.TempFile) error {
	return nil
}

// Write writes data in temp file without any padding bytes for Custom1 format
func (c1 *Custom1) Write(t *storage.TempFile, data []byte) (numBytes int, err error) {
	if numBytes, err = t.File.Write(data); err != nil {
		fmt.Printf("Error in writing byte array to file %s. Finished writing %d bytes", err, numBytes)
		return numBytes, err
	}
	return numBytes, nil
}

// Close closes file with new line character and no header for Custom1 format
func (c1 *Custom1) Close(t *storage.TempFile) (bool, string) {
	if _, err := t.File.Write(storage.NEWLINE_EOF); err != nil {
		fmt.Printf("Error in writing end of file: %s", err)
		return false, ""
	}
	return true, "%n+%v_%D+%T_%N_file%Q"
}
