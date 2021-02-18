package storage

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var sequence uint32 = 0

const (
	DATE_FORMAT        = "01_02_2006"
	TIME_FORMAT        = "05_04_15"
	PLACEHOLDER_FORMAT = "%(\\d*)([A-Za-z%]{1})" // ex: "%n+%v_%D+%T_%N_file%4Q.u"
)

type parseContext struct {
	nodeIdSuffix string
	vpnId        string
	totalCdrs    uint64
	maxSequence  uint32
	mutex        sync.Mutex
}

func (ctx *parseContext) getFileSeqNum() uint32 {
	atomic.AddUint32(&sequence, 1)
	if sequence > ctx.maxSequence {
		ctx.mutex.Lock()
		if sequence > ctx.maxSequence {
			sequence = 0
		}
		ctx.mutex.Unlock()
	}
	return sequence
}

type parser interface {
	parse(string, *parseContext) string
}

type fileNameParser struct{}

func newFileNameParser() parser {
	return &fileNameParser{}
}

func (p *fileNameParser) parse(result string, ctx *parseContext) string {
	pattern := regexp.MustCompile(PLACEHOLDER_FORMAT)
	if matches := pattern.FindAllStringSubmatch(result, -1); matches != nil {
		for _, match := range matches {
			t := time.Now()
			var evaluated string
			switch match[2] {
			case "n":
				evaluated = ctx.nodeIdSuffix
			case "v":
				evaluated = ctx.vpnId
			case "D":
				evaluated = t.Format(DATE_FORMAT)
			case "T":
				evaluated = t.Format(TIME_FORMAT)
			case "N":
				evaluated = fmt.Sprint(ctx.totalCdrs)
			case "Q":
				evaluated = fmt.Sprint(ctx.getFileSeqNum())
			}
			if evaluated != "" {
				if i, err := strconv.Atoi(match[1]); err == nil {
					evaluated = string(PreFill([]byte(evaluated), i, byte('0')))
				}
				result = strings.ReplaceAll(result, match[0], evaluated)
			}
		}
	}
	return result
}
