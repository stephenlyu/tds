package util

import (
	"time"
	"os"
	"archive/zip"
	"io"
	"path/filepath"
	"math"
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

func Round(val float64, places int) float64 {
	var t float64
	f := math.Pow10(places)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}