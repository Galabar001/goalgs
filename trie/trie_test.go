package trie

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tr := New[int]()
	if tr == nil {
		t.Errorf("New[int]() returned nil, want non-nil")
	}
}

func TestAdd(t *testing.T) {
	tr := New[int]()

	// Test adding first key
	added := tr.Add("hello", 1)
	if added {
		t.Errorf("Add(\"hello\", 1) = %v, want false (new key)", added)
	}

	// Test replacing existing key
	added = tr.Add("hello", 2)
	if !added {
		t.Errorf("Add(\"hello\", 2) on existing key = %v, want true (value replaced)", added)
	}
	value, found := tr.Get("hello")
	if !found || value != 2 {
		t.Errorf("Get(\"hello\") = %d, %v, want 2, true after replacement", value, found)
	}

	// Test adding prefix
	added = tr.Add("help", 3)
	if added {
		t.Errorf("Add(\"help\", 3) = %v, want false (new key)", added)
	}

	// Test adding longer key
	added = tr.Add("hello world", 4)
	if added {
		t.Errorf("Add(\"hello world\", 4) = %v, want false (new key)", added)
	}
}

func TestGet(t *testing.T) {
	tr := New[int]()

	// Test empty trie
	_, found := tr.Get("hello")
	if found {
		t.Errorf("Get(\"hello\") on empty trie found = %v, want false", found)
	}

	// Test after adding
	tr.Add("hello", 1)
	value, found := tr.Get("hello")
	if !found {
		t.Errorf("Get(\"hello\") after Add found = %v, want true", found)
	}
	if value != 1 {
		t.Errorf("Get(\"hello\") = %d, want 1", value)
	}

	// Test non-existent key
	_, found = tr.Get("hell")
	if found {
		t.Errorf("Get(\"hell\") found = %v, want false", found)
	}

	// Test prefix with value and children
	tr.Add("car", 1)
	tr.Add("card", 2)
	value, found = tr.Get("card")
	if !found || value != 2 {
		t.Errorf("Get(\"card\") = %d, %v, want 2, true", value, found)
	}

	// Test the prefix key still works
	value, found = tr.Get("car")
	if !found || value != 1 {
		t.Errorf("Get(\"car\") = %d, %v, want 1, true", value, found)
	}
}

func TestRemove(t *testing.T) {
	tr := New[int]()

	// Test empty trie
	found := tr.Remove("hello")
	if found {
		t.Errorf("Remove(\"hello\") on empty trie = %v, want false", found)
	}

	// Test after adding
	tr.Add("hello", 1)
	found = tr.Remove("hello")
	if !found {
		t.Errorf("Remove(\"hello\") after Add = %v, want true", found)
	}

	// Test removed key
	_, found = tr.Get("hello")
	if found {
		t.Errorf("Get(\"hello\") after Remove found = %v, want false", found)
	}
}

func TestClear(t *testing.T) {
	tr := New[int]()
	tr.Add("hello", 1)
	tr.Clear()
	_, found := tr.Get("hello")
	if found {
		t.Errorf("Get(\"hello\") after Clear found = %v, want false", found)
	}
}

func TestCompact(t *testing.T) {
	tr := New[int]()
	tr.Add("car", 1)
	tr.Add("card", 2)
	tr.Add("cat", 3)
	tr.Compact()

	// Test compacted trie
	value, found := tr.Get("car")
	if !found || value != 1 {
		t.Errorf("Get(\"car\") after Compact = %d, %v, want 1, true", value, found)
	}
	value, found = tr.Get("card")
	if !found || value != 2 {
		t.Errorf("Get(\"card\") after Compact = %d, %v, want 2, true", value, found)
	}
	value, found = tr.Get("cat")
	if !found || value != 3 {
		t.Errorf("Get(\"cat\") after Compact = %d, %v, want 3, true", value, found)
	}
}

// Benchmarks

func BenchmarkAdd(b *testing.B) {
	tr := New[int]()
	keys := []string{"hello", "help", "hello world", "helicopter", "helium"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Add(keys[i%len(keys)], i)
	}
}

func BenchmarkGet(b *testing.B) {
	tr := New[int]()
	keys := []string{"hello", "help", "hello world", "helicopter", "helium"}
	for i, key := range keys {
		tr.Add(key, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Get(keys[i%len(keys)])
	}
}

func BenchmarkRemove(b *testing.B) {
	tr := New[int]()
	keys := []string{"hello", "help", "hello world", "helicopter", "helium"}
	for i, key := range keys {
		tr.Add(key, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Remove(keys[i%len(keys)])
		tr.Add(keys[i%len(keys)], i) // Reset for next iteration
	}
}

const numEntries = 1_000_000

func generateRandomKeys(count, maxLen int) []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	keys := make(map[string]struct{}, count)
	for len(keys) < count {
		length := 1 + rng.Intn(maxLen)
		var sb strings.Builder
		for j := 0; j < length; j++ {
			sb.WriteByte('a' + byte(rng.Intn(26)))
		}
		keys[sb.String()] = struct{}{}
	}
	result := make([]string, 0, count)
	for k := range keys {
		result = append(result, k)
	}
	return result
}

// Non-Compacted Benchmarks

func BenchmarkTrieAddRandom(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr := New[int]()
		for j, key := range keys {
			tr.Add(key, j)
		}
	}
}

func BenchmarkTrieAddRandomCompacted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr := New[int]()
		for j, key := range keys {
			tr.Add(key, j)
		}
		tr.Compact()
	}
}

func BenchmarkMapAddRandom(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]int)
		for j, key := range keys {
			m[key] = j
		}
	}
}

func BenchmarkTrieAddSorted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr := New[int]()
		for j, key := range keys {
			tr.Add(key, j)
		}
	}
}

func BenchmarkTrieAddSortedCompacted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr := New[int]()
		for j, key := range keys {
			tr.Add(key, j)
		}
		tr.Compact()
	}
}

func BenchmarkMapAddSorted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]int)
		for j, key := range keys {
			m[key] = j
		}
	}
}

func BenchmarkTrieGetRandom(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	tr := New[int]()
	for i, key := range keys {
		tr.Add(key, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Get(keys[i%numEntries])
	}
}

func BenchmarkTrieGetRandomCompacted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	tr := New[int]()
	for i, key := range keys {
		tr.Add(key, i)
	}
	tr.Compact()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Get(keys[i%numEntries])
	}
}

func BenchmarkMapGetRandom(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	m := make(map[string]int)
	for i, key := range keys {
		m[key] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[keys[i%numEntries]]
	}
}

func BenchmarkTrieGetSorted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	tr := New[int]()
	for i, key := range keys {
		tr.Add(key, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Get(keys[i%numEntries])
	}
}

func BenchmarkTrieGetSortedCompacted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	tr := New[int]()
	for i, key := range keys {
		tr.Add(key, i)
	}
	tr.Compact()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Get(keys[i%numEntries])
	}
}

func BenchmarkMapGetSorted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	m := make(map[string]int)
	for i, key := range keys {
		m[key] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[keys[i%numEntries]]
	}
}

func BenchmarkTrieRemoveRandom(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	tr := New[int]()
	for i, key := range keys {
		tr.Add(key, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Remove(keys[i%numEntries])
	}
}

func BenchmarkTrieRemoveRandomCompacted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	tr := New[int]()
	for i, key := range keys {
		tr.Add(key, i)
	}
	tr.Compact()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Remove(keys[i%numEntries])
	}
}

func BenchmarkMapRemoveRandom(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	m := make(map[string]int)
	for i, key := range keys {
		m[key] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(m, keys[i%numEntries])
	}
}

func BenchmarkTrieRemoveSorted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	tr := New[int]()
	for i, key := range keys {
		tr.Add(key, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Remove(keys[i%numEntries])
	}
}

func BenchmarkTrieRemoveSortedCompacted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	tr := New[int]()
	for i, key := range keys {
		tr.Add(key, i)
	}
	tr.Compact()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Remove(keys[i%numEntries])
	}
}

func BenchmarkMapRemoveSorted(b *testing.B) {
	keys := generateRandomKeys(numEntries, 128)
	sort.Strings(keys)
	m := make(map[string]int)
	for i, key := range keys {
		m[key] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(m, keys[i%numEntries])
	}
}
