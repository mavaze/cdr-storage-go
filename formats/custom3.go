package format

import (
	"fmt"
	"storage-go"
)

// custom3:
//  header: no
//  contents: CDR1CDR2CDR3 ...CDRn
//  eof: no
//  name: <node-id-suffix+vpn-id>_<date>+<time>_<total-cdrs>_file<fileseqnum>.u

// Custom3 holds necessary functions to build Custom3 file format
type Custom3 struct{}

// NewCustom3Format instantiates custom1 file format
func NewCustom3Format() *Custom3 {
	return &Custom3{}
}

// Header writes no header for Custom3 format
func (c1 *Custom3) Header(t *storage.TempFile) error {
	return nil
}

// Write writes data in temp file without any padding bytes for Custom3 format
func (c1 *Custom3) Write(t *storage.TempFile, data []byte) (numBytes int, err error) {
	if numBytes, err = t.File.Write(data); err != nil {
		fmt.Printf("Error in writing byte array to file %s. Finished writing %d bytes", err, numBytes)
		return numBytes, err
	}
	return numBytes, nil
}

// Close closes file with no EOF and no header for Custom3 format
func (c1 *Custom3) Close(t *storage.TempFile) (bool, string) {
	return true, "%n+%v_%D+%T_%N_file%Q.u"
}
