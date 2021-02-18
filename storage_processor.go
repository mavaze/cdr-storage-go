package storage

import (
	"fmt"
	"sync"
	"time"
)

type item = []byte

const (
	ROTATION_EVALUATION_INTERVAL_MS = 5 * 1000 * time.Millisecond

	TEMP_BASE_PATH  = "/opt/docker/filesystem/cdr/archive"
	LOCAL_BASE_PATH = "/opt/docker/filesystem/cdr/active"
)

// Processor orchstrates storage facility
type Processor struct {
	receiver       chan []item
	sender         chan []item
	done           chan bool
	mutex          sync.Mutex
	tmpFile        *TempFile
	persistentFile *PersistentFile
	fileRotator    RotationPolicy
}

var instance *Processor
var once sync.Once

// GetStorageProcessor returns singleton storage processor
func GetStorageProcessor(format FileFormat) *Processor {
	once.Do(func() {
		fileRotator := &FileRotator{}
		fileRotator.addPolicy(&CountBasedRotator{maxCount: 100})
		fileRotator.addPolicy(&VolumeBasedRotator{maxSize: 10 * 1024})
		fileRotator.addPolicy(&TimeBasedRotator{intervalInMillis: 30 * 1000})
		instance = &Processor{
			receiver:    make(chan []item),
			sender:      make(chan []item),
			done:        make(chan bool),
			fileRotator: fileRotator,
			tmpFile:     &TempFile{path: TEMP_BASE_PATH, format: format},
			persistentFile: &PersistentFile{
				path:         LOCAL_BASE_PATH,
				parser:       newFileNameParser(),
				parseContext: &parseContext{nodeIdSuffix: "ga", vpnId: "sgw4", maxSequence: 9999},
			},
		}
		go instance.listener()
	})
	return instance
}

// Save forwards packets to channel for further processing
func (p *Processor) Save(items []item) {
	p.receiver <- items
}

// Close on graceful shutdown closes channel and stops timer
func (p *Processor) Close() {
	p.done <- true
}

func (p *Processor) listener() {
	defer func() {
		fmt.Println("Closing channel")
		close(p.receiver)
	}()

	monitorTimer := time.NewTicker(ROTATION_EVALUATION_INTERVAL_MS)

	keepRunning := true
	for {
		if !keepRunning {
			break
		}
		select {
		case <-p.done:
			monitorTimer.Stop()
			keepRunning = false
			return
		case items, ok := <-p.receiver:
			fmt.Println("in receiver")
			if !ok {
				fmt.Println("channel closed")
				break
			}
			// create temp file if not created yet by monitor routine
			if p.tmpFile.File == nil {
				p.tmpFile.create()
			}
			for _, item := range items {
				if err := p.tmpFile.write(item); err != nil {
					p.sender <- items
				}
			}
		case <-monitorTimer.C:
			fmt.Println("Evaluating all rotation policies at " + time.Now().Format("2006-01-02 15:04:05"))
			if toBeRotated := p.fileRotator.rotate(p.tmpFile); toBeRotated {
				p.persistentFile.persist(p.tmpFile)
				p.tmpFile.create()
			}
		}
	}
}
