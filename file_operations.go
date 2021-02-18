package storage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync/atomic"
	"time"
)

// TempFile points to file to be stored in tmpfs storage
type TempFile struct {
	path      string
	File      *os.File
	format    FileFormat
	count     uint64
	createdAt time.Time
	isClosed  bool
}

// PersistentFile points to file to be stored in local-storage
type PersistentFile struct {
	path         string
	parser       parser
	parseContext *parseContext
}

func (p *PersistentFile) move(oldFile *os.File, fileName string) (err error) {
	newFileName := p.path + "/" + fileName
	if err = os.Rename(oldFile.Name(), newFileName); err == nil {
		fmt.Printf("Moved temp file: %s to persistent storage as %s\n", oldFile.Name(), fileName)
	}
	return err
}

func (p *PersistentFile) persist(tmpFile *TempFile) error {
	if newFilePattern, _ := tmpFile.close(); newFilePattern != "" {
		p.parseContext.totalCdrs = tmpFile.count
		newFileName := p.parser.parse(newFilePattern, p.parseContext)
		return p.move(tmpFile.File, newFileName)
	}
	return errors.New("Failed to evaluate filename pattern")
}

func (t *TempFile) create() (err error) {
	if t.File, err = ioutil.TempFile(t.path, "tmpfs-*.cdr"); err == nil {
		fmt.Println("Created temp file:", t.File.Name())
		t.count = 0
		t.createdAt = time.Now()
		t.isClosed = false
		if err = t.format.Header(t); err != nil {
			return fmt.Errorf("Failed to insert header to file %+v", err)
		}
	}
	return fmt.Errorf("Failed to create temp file %+v", err)
}

func (t *TempFile) write(data []byte) error {
	if n, err := t.format.Write(t, data); err != nil {
		fmt.Printf("Error in writing byte array to file %s. Finished writing %d bytes", err, n)
		return err
	}
	atomic.AddUint64(&t.count, 1)
	return nil
}

func (t *TempFile) close() (newFilePattern string, err error) {
	if t.File == nil {
		return "", nil
	}
	fmt.Println("Closing temp file", t.File.Name())
	var isClosed bool
	if isClosed, newFilePattern = t.format.Close(t); isClosed {
		if err := t.File.Close(); err == nil {
			t.isClosed = true
			return newFilePattern, nil
		}
		return "", fmt.Errorf("Failed to close temp file: %s because %+v", t.File.Name(), err)
	}
	return "", fmt.Errorf("Failed to post process temp file: %s", t.File.Name())
}
