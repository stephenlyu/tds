package util

import (
	"time"
	"os"
	"archive/zip"
	"io"
	"path/filepath"
	"math"
	"runtime"
	"bytes"
	"compress/zlib"
)

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func Tick() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

func NanoTick() uint64 {
	return uint64(time.Now().UnixNano())
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

func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}

//进行zlib压缩
func ZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func ZlibUnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	io.Copy(&out, r)
	return out.Bytes(), nil
}
