package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildCompositeId(t *testing.T) {
	require.Equal(t, "first:second:third", BuildCompositeId("first", "second", "third"))
}

func TestSplitCompositeId(t *testing.T) {
	result := SplitCompositeId("first:second:third")
	require.Equal(t, "first", result[0])
	require.Equal(t, "second", result[1])
	require.Equal(t, "third", result[2])
}
