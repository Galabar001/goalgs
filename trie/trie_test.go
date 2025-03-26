package trie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tr := New[int]()
	require.NotNil(t, tr)
}

func TestAdd(_ *testing.T) {

}

func TestRemove(t *testing.T) {
	tr := New[int]()
	found := tr.Remove([]byte("the_key"))
	require.False(t, found)
}

func TestGet(t *testing.T) {
	tr := New[int]()
	_, found := tr.Get([]byte("the_key"))
	require.False(t, found)
}

func TestClear(_ *testing.T) {
	tr := New[int]()
	tr.Clear()
}
