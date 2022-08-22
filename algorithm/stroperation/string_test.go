package stroperation

import "testing"

func TestDelLeast(t *testing.T) {
	input := "aabcddd"
	output := DelLeast(input)
	t.Log(output)
}
