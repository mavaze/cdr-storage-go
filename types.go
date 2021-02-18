package storage

type callDataRecord = []byte

var NEWLINE_EOF = []byte{'\n'}

// FileFormat governs CDR file format
type FileFormat interface {
	newFileFormatBuilder
}

type newFileFormatBuilder interface {
	Header(*TempFile) error
	Write(*TempFile, []byte) (int, error)
	Close(*TempFile) (bool, string)
}

func PreFill(data []byte, size int, char byte) []byte {
	var l int
	data = Trim(data, size)
	if l = len(data); l == size {
		return data
	}
	tmp := make([]byte, size)
	for i := range tmp[:size-l] {
		tmp[i] = char
	}
	copy(tmp[size-l:], data)
	return tmp
}

func PostFill(data []byte, size int, char byte) []byte {
	var l int
	data = Trim(data, size)
	if l = len(data); l == size {
		return data
	}
	tmp := make([]byte, size)
	copy(tmp, data)
	for i := range tmp[l:] {
		tmp[i+l] = char
	}
	return tmp
}

func Trim(data []byte, size int) []byte {
	l := len(data)
	if l == size {
		return data
	}
	if l > size {
		return data[l-size:]
	}
	return data
}
