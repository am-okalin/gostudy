package hash

import (
	"bytes"
	"crypto/sha256"
	"encoding"
	"encoding/hex"
	"hash"
	"testing"
)

func fromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

type marshalTest struct {
	name   string
	new    func() hash.Hash
	golden []byte
}

func Test1(t *testing.T) {
	tt := marshalTest{"sha256", sha256.New, fromHex("736861032bed68b99987cae48183b2b049d393d0050868e4e8ba3730e9112b08765929b7c0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f80000000000000000000000000000f9")}

	buf := make([]byte, 256)
	for i := range buf {
		buf[i] = byte(i)
	}

	h := tt.new()
	h.Write(buf[:256])
	sum := h.Sum(nil)

	h2 := tt.new()
	h3 := tt.new()
	const split = 249
	for i := 0; i < split; i++ {
		h2.Write(buf[i : i+1])
	}
	h2m, ok := h2.(encoding.BinaryMarshaler)
	if !ok {
		t.Fatalf("Hash does not implement MarshalBinary")
	}
	enc, err := h2m.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: %v", err)
	}
	if !bytes.Equal(enc, tt.golden) {
		t.Errorf("MarshalBinary = %x, want %x", enc, tt.golden)
	}
	h3u, ok := h3.(encoding.BinaryUnmarshaler)
	if !ok {
		t.Fatalf("Hash does not implement UnmarshalBinary")
	}
	if err := h3u.UnmarshalBinary(enc); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}
	h2.Write(buf[split:])
	h3.Write(buf[split:])
	sum2 := h2.Sum(nil)
	sum3 := h3.Sum(nil)
	if !bytes.Equal(sum2, sum) {
		t.Fatalf("Sum after MarshalBinary = %x, want %x", sum2, sum)
	}
	if !bytes.Equal(sum3, sum) {
		t.Fatalf("Sum after UnmarshalBinary = %x, want %x", sum3, sum)
	}
}

func Test2(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("abcde123"))
	b := h.Sum(nil)
	s := hex.EncodeToString(b)
	t.Log(s)
}
