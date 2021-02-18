package format

import (
	"fmt"
	"storage-go"
)

// custom4:
//  header: no
//  contents: CDR1|CDR2FFFFFF|CDR3FFFFF..|..CDRnFFFF|      ..... where | represents the end of a 2K block
//  eof: no
//  name: <node-id-suffix+vpn-id>_<date>+<time>_<total-cdrs>_file<fileseqnum>.u

// Custom4 holds necessary functions to build Custom4 file format
type Custom4 struct{}

// NewCustom4Format instantiates Custom4 file format
func NewCustom4Format() *Custom4 {
	return &Custom4{}
}

// Header writes no header for Custom4 format
func (c4 *Custom4) Header(t *storage.TempFile) error {
	return nil
}

// Write writes data in temp file in a block of 2K bytes for Custom4 format
func (c4 *Custom4) Write(t *storage.TempFile, data []byte) (numBytes int, err error) {
	paddedBytes := storage.PostFill(data, 2*1024, byte('f'))
	if numBytes, err = t.File.Write(paddedBytes); err != nil {
		fmt.Printf("Error in writing byte array to file %s. Finished writing %d bytes", err, numBytes)
		return numBytes, err
	}
	return numBytes, nil
}

// Close closes file with EOF and no header for Custom4 format
func (c4 *Custom4) Close(t *storage.TempFile) (bool, string) {
	return true, "%n+%v_%D+%T_%N_file%Q.u"
}
