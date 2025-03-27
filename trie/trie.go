package trie

type Trie[T any] interface {
	Add(key string, value T) bool
	Remove(key string) bool
	Get(key string) (T, bool)
	Clear()
	Compact()
}

type trieNode[T any] struct {
	children map[byte]*trieNode[T]
	prefix   []byte
	value    T
	hasValue bool
}

type trie[T any] struct {
	root *trieNode[T]
}

func New[T any]() Trie[T] {
	return &trie[T]{}
}

func (t *trie[T]) Add(key string, value T) bool {
	keyBytes := []byte(key)
	if t.root == nil {
		t.root = &trieNode[T]{
			prefix:   keyBytes,
			value:    value,
			hasValue: true,
			children: make(map[byte]*trieNode[T]),
		}
		return false
	}

	current := t.root
	offset := 0

	for {
		commonLen := 0
		for i := 0; i < len(current.prefix) && offset+i < len(keyBytes); i++ {
			if current.prefix[i] != keyBytes[offset+i] {
				break
			}
			commonLen++
		}

		if commonLen == len(current.prefix) {
			if offset+commonLen == len(keyBytes) {
				if current.hasValue {
					current.value = value
					return true
				}
				current.value = value
				current.hasValue = true
				return false
			}

			nextByte := keyBytes[offset+commonLen]
			if current.children == nil {
				current.children = make(map[byte]*trieNode[T])
			}
			if current.children[nextByte] == nil {
				current.children[nextByte] = &trieNode[T]{
					prefix:   keyBytes[offset+commonLen:],
					value:    value,
					hasValue: true,
					children: make(map[byte]*trieNode[T]),
				}
				return false
			}
			current = current.children[nextByte]
			offset += commonLen
			continue
		}

		// Split the node
		newPrefix := current.prefix[:commonLen]
		remainder := current.prefix[commonLen:]
		oldChild := &trieNode[T]{
			prefix:   remainder,
			value:    current.value,
			hasValue: current.hasValue,
			children: current.children,
		}

		current.prefix = newPrefix
		current.children = make(map[byte]*trieNode[T])
		current.children[remainder[0]] = oldChild

		if offset+commonLen == len(keyBytes) {
			current.value = value
			current.hasValue = true
			return false
		}

		newChild := &trieNode[T]{
			prefix:   keyBytes[offset+commonLen:],
			value:    value,
			hasValue: true,
			children: make(map[byte]*trieNode[T]),
		}
		current.children[keyBytes[offset+commonLen]] = newChild
		return false
	}
}

func (t *trie[T]) Get(key string) (T, bool) {
	keyBytes := []byte(key)
	var zero T
	if t.root == nil {
		return zero, false
	}

	current := t.root
	offset := 0

	for {
		if len(keyBytes)-offset < len(current.prefix) {
			return zero, false
		}

		for i := 0; i < len(current.prefix); i++ {
			if keyBytes[offset+i] != current.prefix[i] {
				return zero, false
			}
		}
		offset += len(current.prefix)

		if offset == len(keyBytes) {
			if current.hasValue {
				return current.value, true
			}
			return zero, false
		}

		nextByte := keyBytes[offset]
		current = current.children[nextByte]
		if current == nil {
			return zero, false
		}
	}
}

func (t *trie[T]) Remove(key string) bool {
	keyBytes := []byte(key)
	if t.root == nil {
		return false
	}

	return removeHelper(t, &t.root, keyBytes, 0)
}

func removeHelper[T any](t *trie[T], nodePtr **trieNode[T], keyBytes []byte, offset int) bool {
	current := *nodePtr
	if current == nil {
		return false
	}

	if len(keyBytes)-offset < len(current.prefix) {
		return false
	}

	for i := 0; i < len(current.prefix); i++ {
		if keyBytes[offset+i] != current.prefix[i] {
			return false
		}
	}

	if offset+len(current.prefix) == len(keyBytes) {
		if current.hasValue {
			current.hasValue = false
			var zero T
			current.value = zero
			return true
		}
		return false
	}

	nextByte := keyBytes[offset+len(current.prefix)]
	if current.children == nil || current.children[nextByte] == nil {
		return false
	}

	nextChild := current.children[nextByte]
	if removeHelper(t, &nextChild, keyBytes, offset+len(current.prefix)) {
		if len(nextChild.children) == 0 && !nextChild.hasValue {
			delete(current.children, nextByte)
		} else {
			current.children[nextByte] = nextChild
		}
		return true
	}
	return false
}

func (t *trie[T]) Clear() {
	t.root = nil
}

func (t *trie[T]) Compact() {
	if t.root == nil {
		return
	}
	t.root = compactNode(t.root)
}

func compactNode[T any](node *trieNode[T]) *trieNode[T] {
	if node == nil {
		return nil
	}

	// Ensure children is a map
	if node.children == nil {
		node.children = make(map[byte]*trieNode[T])
	}

	// Count non-nil children
	childCount := 0
	var onlyChild *trieNode[T]
	for _, child := range node.children {
		if child != nil {
			childCount++
			onlyChild = child
		}
	}

	// Merge with single child if no value at current node
	if childCount == 1 && !node.hasValue {
		mergedPrefix := append(node.prefix, onlyChild.prefix...)
		return &trieNode[T]{
			prefix:   mergedPrefix,
			value:    onlyChild.value,
			hasValue: onlyChild.hasValue,
			children: onlyChild.children,
		}
	}

	// Recursively compact children
	for b, child := range node.children {
		if child != nil {
			node.children[b] = compactNode(child)
		}
	}

	return node
}
