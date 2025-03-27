package trie

type Trie[T any] interface {
	Add(key string, value T) bool
	Remove(key string) bool
	Get(key string) (T, bool)
	Clear()
}

type trieNode[T any] struct {
	children *[256]*trieNode[T]
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
		}
		return false // New key, no previous value
	}

	current := t.root
	offset := 0

	for {
		// Find common prefix length
		commonLen := 0
		for i := 0; i < len(current.prefix) && offset+i < len(keyBytes); i++ {
			if current.prefix[i] != keyBytes[offset+i] {
				break
			}
			commonLen++
		}

		if commonLen == len(current.prefix) {
			// We've matched the entire current prefix
			if offset+commonLen == len(keyBytes) {
				// Key matches exactly
				if current.hasValue {
					current.value = value
					return true // Previous value existed, replaced it
				}
				// Convert prefix node to leaf
				current.value = value
				current.hasValue = true
				return false // No previous value, new key added
			}

			// Need to continue down the tree
			if current.children == nil {
				// Convert to leaf node
				current.value = value
				current.hasValue = true
				current.prefix = keyBytes[offset:]
				return false // New key, no previous value
			}

			// Move to next child
			nextByte := keyBytes[offset+commonLen]
			if current.children[nextByte] == nil {
				current.children[nextByte] = &trieNode[T]{
					prefix:   keyBytes[offset+commonLen:],
					value:    value,
					hasValue: true,
				}
				return false // New key, no previous value
			}
			current = current.children[nextByte]
			offset += commonLen
			continue
		}

		// Split the node
		newPrefix := current.prefix[:commonLen]
		remainder := current.prefix[commonLen:]

		// Create new child for existing branch
		oldChild := &trieNode[T]{
			prefix:   remainder,
			value:    current.value,
			hasValue: current.hasValue,
			children: current.children,
		}

		// Update current node
		current.prefix = newPrefix
		var zero T
		current.value = zero
		current.hasValue = false
		current.children = new([256]*trieNode[T])

		// Add old child to children
		current.children[remainder[0]] = oldChild

		// If we've used up the key, we're done
		if offset+commonLen == len(keyBytes) {
			current.value = value
			current.hasValue = true
			return false // New key, no previous value (split created a new leaf)
		}

		// Add new leaf node
		newChild := &trieNode[T]{
			prefix:   keyBytes[offset+commonLen:],
			value:    value,
			hasValue: true,
		}
		current.children[keyBytes[offset+commonLen]] = newChild
		return false // New key, no previous value
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
		// Check prefix match
		if len(keyBytes)-offset < len(current.prefix) {
			return zero, false
		}

		for i := 0; i < len(current.prefix); i++ {
			if offset+i >= len(keyBytes) || keyBytes[offset+i] != current.prefix[i] {
				return zero, false
			}
		}

		// If we're at a leaf node
		if current.hasValue {
			if offset+len(current.prefix) == len(keyBytes) {
				return current.value, true
			}
			return zero, false
		}

		// Move to next child
		if offset+len(current.prefix) >= len(keyBytes) {
			return zero, false
		}

		nextByte := keyBytes[offset+len(current.prefix)]
		if current.children[nextByte] == nil {
			return zero, false
		}

		current = current.children[nextByte]
		offset += len(current.prefix)
	}
}

func (t *trie[T]) Remove(key string) bool {
	keyBytes := []byte(key)
	if t.root == nil {
		return false
	}

	current := t.root
	offset := 0

	for {
		if len(keyBytes)-offset < len(current.prefix) {
			return false
		}

		for i := 0; i < len(current.prefix); i++ {
			if offset+i >= len(keyBytes) || keyBytes[offset+i] != current.prefix[i] {
				return false
			}
		}

		if current.hasValue {
			if offset+len(current.prefix) == len(keyBytes) {
				var zero T
				current.value = zero
				current.hasValue = false
				return true
			}
			return false
		}

		if offset+len(current.prefix) >= len(keyBytes) {
			return false
		}

		nextByte := keyBytes[offset+len(current.prefix)]
		if current.children[nextByte] == nil {
			return false
		}

		current = current.children[nextByte]
		offset += len(current.prefix)
	}
}

func (t *trie[T]) Clear() {
	t.root = nil
}
