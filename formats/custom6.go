package format

import (
	"fmt"
	"storage-go"
)

// custom6:
//  header: no
//  contents: CDR1|CDR2FFFFFF|CDR3FFFFF..|..CDRnFFFF|      .... where | represents the end of a 8K block
//  eof: no
//  name: <node-id-suffix+vpn-id>_<date>+<time>_<total-cdrs>_file<fileseqnum>.u

// Custom6 holds necessary functions to build Custom6 file format
type Custom6 struct{}

// NewCustom6Format instantiates Custom6 file format
func NewCustom6Format() *Custom6 {
	return &Custom6{}
}

// Header writes no header for Custom6 format
func (c1 *Custom6) Header(t *storage.TempFile) error {
	return nil
}

// write writes data in temp file in a block of 8K bytes for Custom6 format
func (c1 *Custom6) Write(t *storage.TempFile, data []byte) (numBytes int, err error) {
	paddedBytes := storage.PostFill(data, 8*1024, byte('f'))
	if numBytes, err = t.File.Write(paddedBytes); err != nil {
		fmt.Printf("Error in writing byte array to file %s. Finished writing %d bytes", err, numBytes)
		return numBytes, err
	}
	return numBytes, nil
}

// Close closes file with new line character and no header for Custom6 format
func (c1 *Custom6) Close(t *storage.TempFile) (bool, string) {
	return true, "%n+%v_%D+%T_%N_file%Q.u"
}
