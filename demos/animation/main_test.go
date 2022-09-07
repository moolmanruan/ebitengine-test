package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestColor(t *testing.T) {
	require.Equal(t, [3]float32{1, 1, 1}, colorFromHex(0xffffff))
	require.Equal(t, [3]float32{0, 0, 0}, colorFromHex(0x000000))
	r := float32(0x80) / 0xff
	g := float32(0x34) / 0xff
	b := float32(0x69) / 0xff
	require.Equal(t, [3]float32{r, g, b}, colorFromHex(0x803469))
}
