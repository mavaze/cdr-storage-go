package format

import (
	"fmt"
	"storage-go"
)

// custom5:
//  header: no
//  contents: CDR1CDR2CDR3 ...CDRn
//  eof: no
//  name: <node-id-suffix+vpn-id>_<date>+<time>_<total-cdrs>_file<fixed-length-seqnum>.u

// Custom5 holds necessary functions to build Custom5 file format
type Custom5 struct{}

// NewCustom5Format instantiates Custom5 file format
func NewCustom5Format() *Custom5 {
	return &Custom5{}
}

// Header writes no header for Custom5 format
func (c1 *Custom5) Header(t *storage.TempFile) error {
	return nil
}

// Write writes data in temp file without any padding bytes for Custom5 format
func (c1 *Custom5) Write(t *storage.TempFile, data []byte) (numBytes int, err error) {
	if numBytes, err = t.File.Write(data); err != nil {
		fmt.Printf("Error in writing byte array to file %s. Finished writing %d bytes", err, numBytes)
		return numBytes, err
	}
	return numBytes, nil
}

// Close closes file with new line character and no header for Custom5 format
func (c1 *Custom5) Close(t *storage.TempFile) (bool, string) {
	return true, "%n+%v_%D+%T_%N_file%4Q.u"
}
