package util

type RingBuffer struct {
	Buffer []interface{}

	Top int
	Length int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		Buffer: make([]interface{}, size),
	}
}

func (this *RingBuffer) Append(v interface{}) {
	this.Buffer[this.Top] = v
	if this.Length < len(this.Buffer) {
		this.Length++
	}
	this.Top++

	if this.Top >= len(this.Buffer) {
		this.Top = 0
	}
}

func (this *RingBuffer) Get(index int) interface{} {
	if index < 0 || index >= this.Length {
		panic("Index out of range")
	}

	if this.Length < len(this.Buffer) {
		return this.Buffer[index]
	}

	index += this.Top
	index %= this.Length
	return this.Buffer[index]
}

func (this *RingBuffer) GetLength() int {
	return this.Length
}
