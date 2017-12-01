package datasource

import (
	"os"
	"errors"
	"bufio"

	. "github.com/stephenlyu/tds/entity"
)

var (
	ERR_FILE_DAMAGED = errors.New("file damaged")
)

type RecordMarshaller interface {
	ToBytes(record *Record) ([]byte, error)
	FromBytes(bytes []byte, record *Record) error
}

type RecordReader interface {
	Read(start, end int) (error, []Record)
	Count() (error, int)
}

type RecordWriter interface {
	Write(from int, data []Record) error
	WriteRaw(from int, data []byte) error
}

type recordReader struct {
	file *os.File
	recordSize int
	marshaller RecordMarshaller
}

func NewRecordReader(file *os.File, recordSize int, marshaller RecordMarshaller) RecordReader {
	return &recordReader{file, recordSize, marshaller}
}

func (this *recordReader) Read(start, end int) (error, []Record) {
	var err error
	if end == -1 {
		err, end = this.Count()
		if err != nil {
			return err, nil
		}
	}

	result := make([]Record, end - start)

	this.file.Seek(int64(start) * int64(this.recordSize), os.SEEK_SET)

	buf := make([]byte, 100 * this.recordSize)

	var n int
	current := start
	for current < end {
		count := end - current
		if count > 100 {
			count = 100
		}
		n, err = this.file.Read(buf[0:count * this.recordSize])
		if err != nil {
			return err, nil
		}

		if n < count * this.recordSize {
			return errors.New("read less data"), nil
		}

		for i := 0; i < count; i++ {
			result[current-start+i] = Record{}
			err = this.marshaller.FromBytes(buf[i * this.recordSize: (i+1) * this.recordSize], &result[current-start+i])
			if err != nil {
				return err, nil
			}
		}
		current += count
	}

	return nil, result
}

func (this *recordReader) Count() (error, int) {
	stat, err := this.file.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0
		}
		return err, 0
	}

	if stat.Size() % int64(this.recordSize) != 0 {
		return ERR_FILE_DAMAGED, 0
	}

	return nil, int(stat.Size() / int64(this.recordSize))
}

type recordWriter struct {
	file *os.File
	recordSize int
	marshaller RecordMarshaller
}

func NewRecordWriter(file *os.File, recordSize int, marshaller RecordMarshaller) RecordWriter {
	return &recordWriter{file, recordSize, marshaller}
}

func (this *recordWriter) Write(from int, data []Record) error {
	err := this.file.Truncate(int64(from) * int64(this.recordSize))
	if err != nil {
		return err
	}

	_, err = this.file.Seek(int64(from) * int64(this.recordSize), os.SEEK_SET)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(this.file)

	for i := 0; i < len(data); i++ {
		r := &data[i]
		bytes, err := this.marshaller.ToBytes(r)
		if err != nil {
			return err
		}
		_, err = writer.Write(bytes)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func (this *recordWriter) WriteRaw(from int, data []byte) error {
	_, err := this.file.Seek(int64(from) * int64(this.recordSize), os.SEEK_SET)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(this.file)

	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return writer.Flush()
}
