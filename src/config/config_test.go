package config

import (
	"github.com/bry00/m/buffers"
	"testing"
)

func TestBlockSizeLimitMB(t *testing.T) {
	c := NewDefaultConfig()
	expected := buffers.DefaultBlockSizeLimit / buffers.MB
	got := c.DataBuffer.BlockSizeLimitMB
	if got != expected {
		t.Errorf("DataBuffer.BlockSizeLimitMB ==> %d; want %d", got, expected)
	}
}

func TestTotalSizeLimitMB(t *testing.T) {
	c := NewDefaultConfig()
	expected := buffers.DefaultTotalSizeLimit / buffers.MB
	got := c.DataBuffer.TotalSizeLimitMB
	if got != expected {
		t.Errorf("DataBuffer.TotalSizeLimitMB ==> %d; want %d", got, expected)
	}
}
