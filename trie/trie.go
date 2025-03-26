package trie

type UTF8Slice []byte

type Trie[T any] interface {
	Add(key UTF8Slice, value T) bool
	Remove(key UTF8Slice) bool
	Get(key UTF8Slice) (T, bool)
	Clear()
}

type trieNode struct {
}

type trie[T any] struct {
	root *trieNode
}

func New[T any]() Trie[T] {
	return &trie[T]{}
}

func (t *trie[T]) Add(key UTF8Slice, value T) bool {
	return false
}

func (t *trie[T]) Remove(key UTF8Slice) bool {
	return false
}

func (t *trie[T]) Get(key UTF8Slice) (T, bool) {
	var value T
	return value, false
}

func (t *trie[T]) Clear() {
	t.root = nil
}
