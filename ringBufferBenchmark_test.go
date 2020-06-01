package ringbuffer

import (
	"strings"
	"testing"
)

func BenchmarkNewRingBuffer(b *testing.B) {
	rb := NewRingBuffer(1024)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rb.Write(data)
		_, _ = rb.Read(buf)
	}
}
