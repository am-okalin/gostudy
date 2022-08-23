package strings

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	r := strings.NewReader("0123456789")
	b := make([]byte, 5, 10)
	n, err := r.Read(b)
	t.Log(n, err, string(b))
}

func TestSeeker(t *testing.T) {
	r := strings.NewReader("0123456789")
	tests := []struct {
		off     int64
		seek    int
		n       int
		want    string
		wantpos int64
		readerr error
		seekerr string
	}{
		{seek: io.SeekStart, off: 0, n: 20, want: "0123456789"},
		{seek: io.SeekStart, off: 1, n: 1, want: "1"},
		{seek: io.SeekCurrent, off: 1, wantpos: 3, n: 2, want: "34"},
		{seek: io.SeekStart, off: -1, seekerr: "strings.Reader.Seek: negative position"},
		{seek: io.SeekStart, off: 1 << 33, wantpos: 1 << 33, readerr: io.EOF},
		{seek: io.SeekCurrent, off: 1, wantpos: 1<<33 + 1, readerr: io.EOF},
		{seek: io.SeekStart, n: 5, want: "01234"},
		{seek: io.SeekCurrent, n: 5, want: "56789"},
		{seek: io.SeekEnd, off: -1, n: 1, wantpos: 9, want: "9"},
		{seek: io.SeekEnd, off: 1, n: 1, wantpos: 11, readerr: io.EOF},
	}

	for i, tt := range tests {
		index := r.Size() - int64(r.Len())
		pos, err := r.Seek(tt.off, tt.seek)
		t.Log("当前计数为：", index, "->", pos)
		if err == nil && tt.seekerr != "" {
			t.Errorf("%d. want seek error %q", i, tt.seekerr)
			continue
		}
		if err != nil && err.Error() != tt.seekerr {
			t.Errorf("%d. seek error = %q; want %q", i, err.Error(), tt.seekerr)
			continue
		}
		if tt.wantpos != 0 && tt.wantpos != pos {
			t.Errorf("%d. pos = %d, want %d", i, pos, tt.wantpos)
		}
		buf := make([]byte, tt.n)
		n, err := r.Read(buf)
		if err != tt.readerr {
			t.Errorf("%d. read = %v; want %v", i, err, tt.readerr)
			continue
		}
		got := string(buf[:n])
		if got != tt.want {
			t.Errorf("%d. got %q; want %q", i, got, tt.want)
		}
	}
}

func TestWriterTo(t *testing.T) {
	const str = "0123456789"
	for i := 0; i <= len(str); i++ {
		s := str[i:]
		r := strings.NewReader(s)
		var b bytes.Buffer
		n, err := r.WriteTo(&b)
		if expect := int64(len(s)); n != expect {
			t.Errorf("got %v; want %v", n, expect)
		}
		if err != nil {
			t.Errorf("for length %d: got error = %v; want nil", len(s), err)
		}
		if b.String() != s {
			t.Errorf("got string %q; want %q", b.String(), s)
		}
		if r.Len() != 0 {
			t.Errorf("reader contains %v bytes; want 0", r.Len())
		}
	}
}
