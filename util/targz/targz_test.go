package targz

import "testing"

func TestCompress(t *testing.T) {
	Compress("data", "temp/data.tar.gz")
}

func TestDeCompress(t *testing.T) {
	DeCompress("temp/data.tar.gz", "temp")
}