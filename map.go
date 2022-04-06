// ordered contains ordered data structures
package ordered

import (
	"fmt"
	"reflect"
	"sync"
)

// Map is an ordered map data structure that is safe for concurrent use by
// multiple goroutines without additional locking or coordination.
//
// The zero Map is empty and ready for use. A Map must not be copied after first
// use.
type Map[K comparable, V any] struct {
	order []K
	dirty map[K]V
	mu    sync.RWMutex
}

// Delete deletes the vlaue for a key
func (m *Map[K, V]) Delete(key K) {
	m.LoadAndDelete(key)
}

// Index loads the key and value of the key at index n. The loaded result
// reports whether the index was in range. Negative value of n index from the
// end of the Map.
func (m *Map[K, V]) Index(n int) (key K, value V, loaded bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if n >= 0 {
		if len(m.order) <= n {
			return
		}
		key = m.order[n]
	} else {
		if n < -len(m.order) {
			return
		}
		key = m.order[len(m.order)+n]
	}

	value, loaded = m.dirty[key]
	return
}

// Len returns the number of keys in Map
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.order)
}

// Load returns the value stored in the map for a key, or nil if no value is
// present. The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.load(key)
}

func (m *Map[K, V]) load(key K) (value V, ok bool) {
	value, ok = m.dirty[key]
	return
}

// LoadAndDelete deletes the value for a key, returning the previous value if
// any. The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 0; i < len(m.order); i++ {
		if m.order[i] != key {
			continue
		}

		m.order = append(m.order[:i], m.order[i+1:]...)
		break
	}

	value, loaded = m.dirty[key]
	delete(m.dirty, key)
	return
}

// LoadAndDeleteFirst delered the last key, returning the kay and its previous
// value if any. The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDeleteFirst() (key K, value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.order) < 1 {
		return
	}

	key = m.order[0]

	value, loaded = m.dirty[key]
	delete(m.dirty, key)
	return
}

// LoadAndDeleteLast delered the last key, returning the kay and its previous
// value if any. The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDeleteLast() (key K, value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.order) < 1 {
		return
	}

	key = m.order[len(m.order)-1]

	value, loaded = m.dirty[key]
	delete(m.dirty, key)
	return
}

// LoadOrStore returns the existing value for the key if present. Otherwise, it
// stores and returns the given value, adding it to the end. The loaded result
// is true if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	actual, loaded = m.load(key)
	if !loaded {
		m.store(key, value)
		actual = value
	}
	return
}

// Range calls f sequentially for each key and value present in the map. If f
// returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: nevery index will be visited in order, but if the value for any key
// stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
func (m *Map[K, V]) Range(f func(index int, key K, value V) bool) {
	m.mu.RLock()

	if len(m.order) < 1 {
		m.mu.RUnlock()
		return
	}

	for index, key := range m.order {
		if index == 0 {
			m.mu.RUnlock()
		}

		if f(index, key, m.dirty[key]) {
			continue
		}

		return
	}
}

// Store sets the value for a key adding it to the end if it was not in the map.
func (m *Map[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store(key, value)
}

func (m *Map[K, V]) store(key K, value V) {
	if _, ok := m.dirty[key]; !ok {
		m.order = append(m.order, key)
	}

	if m.dirty == nil {
		m.dirty = make(map[K]V)
	}
	m.dirty[key] = value
}

// StoreFirst sets the value for a key adding it to the beginning if it was not
// in the map.
func (m *Map[K, V]) StoreFirst(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.dirty[key]; !ok {
		m.order = append([]K{key}, m.order...)
	}

	if m.dirty == nil {
		m.dirty = make(map[K]V)
	}
	m.dirty[key] = value
}

// String formats the map for printing
func (m *Map[K, V]) String() string {
	return typeName(m) + m.string()
}

func (m *Map[K, V]) string() (s string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s = "["
	var space string
	for i := range m.order {
		s += space + fmt.Sprint(m.order[i]) + ":" + fmt.Sprint(m.dirty[m.order[i]])
		space = " "
	}
	s += "]"
	return
}

// Swap swaps the position of the keys at indicies i and j.
func (m *Map[K, V]) Swap(i, j int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if i < 0 || i >= len(m.order) || j < 0 || j >= len(m.order) {
		return
	}

	m.order[i], m.order[j] = m.order[j], m.order[i]
}

// Ordered represents all orderable types.
//
// Deprecated: This will be removed when the constraints package is added to
// the standard library.
type Ordered interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64 |
		~string
}

// SortMap is a Map which fully impliments sort.Interface. Sort is not
// sortable if the key type is float and NaN is used as a key.
type SortMap[K Ordered, V any] struct {
	Map[K, V]
}

// Less returns true if the key at index i is less than the key at index j.
func (m *SortMap[K, V]) Less(i, j int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if i < 0 || i >= len(m.order) || j < 0 || j >= len(m.order) {
		return false
	}

	return m.order[i] < m.order[j]
}

// String formats the map for printing
func (m *SortMap[K, V]) String() string {
	return typeName(m) + m.string()
}

func typeName(t any) string {
	elem := reflect.TypeOf(t).Elem()
	return elem.PkgPath() + "." + elem.Name()
}
