package util

import (
	"time"
	"os"
	"archive/zip"
	"io"
	"path/filepath"
)

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func Tick() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

func UnzipFile(fileName string, outputDir string) error {
	os.MkdirAll(outputDir, 0777)

	r, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	unzipOneFile := func (f *zip.File, destFile string) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		writer, err := os.Create(destFile)
		if err != nil {
			return err
		}
		defer writer.Close()

		_, err = io.CopyN(writer, rc, int64(f.UncompressedSize64))
		return err
	}

	for _, f := range r.File {
		destFile := filepath.Join(outputDir, f.Name)
		err = unzipOneFile(f, destFile)
		if err != nil {
			return err
		}
	}
	return nil
}
