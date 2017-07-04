package tds

import (
	"os"
	"errors"
	"bufio"
)


type RecordReader interface {
	Read(start, end int) (error, []Record)
	Count() (error, int)
}

type RecordWriter interface {
	Write(from int, data []Record) error
}

type recordReader struct {
	file *os.File
}

func NewRecordReader(file *os.File) RecordReader {
	return &recordReader{file}
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

	this.file.Seek(int64(start) * int64(recordSize), os.SEEK_SET)

	buf := make([]byte, 100 * recordSize)

	var n int
	current := start
	for {
		count := end - current
		if count > 100 {
			count = 100
		}
		n, err = this.file.Read(buf[0:count * recordSize])
		if err != nil {
			return err, nil
		}

		if n < count * recordSize {
			return errors.New("read less data"), nil
		}

		for i := 0; i < count; i++ {
			result = append(result, Record{})
			err = RecordFromBytes(buf[i * recordSize: (i+1) * recordSize], &result[len(result) - 1])
			if err != nil {
				return err, nil
			}
		}
	}

	return nil, result
}

func (this *recordReader) Count() (error, int) {
	stat, err := this.file.Stat()
	if err != nil {
		return err, 0
	}

	if stat.Size() % recordSize != 0 {
		return errors.New("file damaged"), 0
	}

	return nil, int(stat.Size() / recordSize)
}

type recordWriter struct {
	file *os.File
}

func NewRecordWriter(file *os.File) RecordWriter {
	return &recordWriter{file}
}

func (this *recordWriter) Write(from int, data []Record) error {
	this.file.Seek(int64(from) * int64(recordSize), os.SEEK_SET)

	writer := bufio.NewWriter(this.file)

	for i := 0; i < len(data); i++ {
		r := &data[i]
		_, err := writer.Write(r.Bytes())
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
