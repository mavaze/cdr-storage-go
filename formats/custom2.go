package format

import (
	"fmt"
	"storage-go"
)

// custom2:
//  header: 24 byte
//  contents: LEN1CDR1LEN2CDR2LEN3CDR3...LENnCDRn
//  eof: no
//  name: <node-id-suffix+vpn-id>_<date>+<time>_<total-cdrs>_file<fileseqnum>.u

// Custom2 holds necessary functions to build Custom2 file format
type Custom2 struct{}

// NewCustom2Format instantiates custom2 file format
func NewCustom2Format() *Custom2 {
	return &Custom2{}
}

// Header writes 24 byte header for Custom2 format
func (c2 *Custom2) Header(t *storage.TempFile) error {
	return nil
}

// Write writes data in temp file without any padding bytes for Custom2 format
func (c2 *Custom2) Write(t *storage.TempFile, data []byte) (numBytes int, err error) {
	if numBytes, err = t.File.Write(data); err != nil {
		fmt.Printf("Error in writing byte array to file %s. Finished writing %d bytes", err, numBytes)
		return numBytes, err
	}
	return numBytes, nil
}

// Close closes file with no EOF but update in header for Custom2 format
func (c2 *Custom2) Close(t *storage.TempFile) (bool, string) {
	return true, "%n+%v_%D+%T_%N_file%Q.u"
}
