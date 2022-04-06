package ordered

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"sync/atomic"
	"testing"
)

func TestDelete(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder, wantOrder []string
		startingMap, wantMap     map[string]int
		key                      string
	}{
		"nil_delete": {
			key: "notakey",
		},
		"empty_delete": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			key:           "notakey",
		},
		"delete_none": {
			startingOrder: []string{"one"},
			startingMap:   map[string]int{"one": 1},
			wantOrder:     []string{"one"},
			wantMap:       map[string]int{"one": 1},
			key:           "notakey",
		},
		"delete_one": {
			startingOrder: []string{"one"},
			startingMap:   map[string]int{"one": 1},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			key:           "one",
		},
		"delete_five": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			key:           "five",
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			m.Delete(test.key)
			if !reflect.DeepEqual(m.dirty, test.wantMap) {
				t.Errorf("Unexpected map content\nactual: %#v\nwant  : %#v", m.dirty, test.wantMap)
			}
			if !reflect.DeepEqual(m.order, test.wantOrder) {
				t.Errorf("Unexpected order content\nactual: %#v\nwant  : %#v", m.order, test.wantOrder)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder []string
		startingMap   map[string]int
		index         int
		wantKey       string
		wantValue     int
		wantLoaded    bool
	}{
		"nil_index": {
			index:      0,
			wantKey:    "",
			wantValue:  0,
			wantLoaded: false,
		},
		"empty_index": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			index:         0,
			wantKey:       "",
			wantValue:     0,
			wantLoaded:    false,
		},
		"index_zero": {
			startingOrder: []string{"zero"},
			startingMap:   map[string]int{"zero": 0},
			index:         0,
			wantKey:       "zero",
			wantValue:     0,
			wantLoaded:    true,
		},
		"index_three": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			index:         3,
			wantKey:       "three",
			wantValue:     3,
			wantLoaded:    true,
		},
		"index_out_of_range": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			index:         10,
			wantKey:       "",
			wantValue:     0,
			wantLoaded:    false,
		},
		"index_last": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			index:         -1,
			wantKey:       "nine",
			wantValue:     9,
			wantLoaded:    true,
		},
		"index_zero_backwards": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			index:         -10,
			wantKey:       "zero",
			wantValue:     0,
			wantLoaded:    true,
		},
		"index_out_of_range_backwards": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			index:         -11,
			wantKey:       "",
			wantValue:     0,
			wantLoaded:    false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			key, value, loaded := m.Index(test.index)
			if key != test.wantKey {
				t.Errorf("Unexpected key, wanted %q but got %q", test.wantKey, key)
			}
			if value != test.wantValue {
				t.Errorf("Unexpected value, wanted %d but got %d", test.wantValue, value)
			}
			if loaded != test.wantLoaded {
				t.Errorf("Unexpected OK, wanted %t but got %t", test.wantLoaded, loaded)
			}
		})
	}
}

func TestLen(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder []string
		startingMap   map[string]int
		want          int
	}{
		"nil": {
			want: 0,
		},
		"empty": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			want:          0,
		},
		"one": {
			startingOrder: []string{"one"},
			startingMap:   map[string]int{"one": 1},
			want:          1,
		},
		"ten": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			want:          10,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			ln := m.Len()
			if ln != test.want {
				t.Errorf("Unexpected length, wanted %d but got %d", test.want, ln)
			}
		})
	}
}

func TestLess(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder []int
		startingMap   map[int]string
		i, j          int
		want          bool
	}{
		"nill_less": {
			i:    0,
			j:    1,
			want: false,
		},
		"empty_less": {
			startingOrder: []int{},
			startingMap:   map[int]string{},
			i:             0,
			j:             1,
			want:          false,
		},
		"less": {
			startingOrder: []int{1, 2, 3},
			startingMap:   map[int]string{1: "one", 2: "two", 3: "three"},
			i:             0,
			j:             1,
			want:          true,
		},
		"more": {
			startingOrder: []int{1, 2, 3},
			startingMap:   map[int]string{1: "one", 2: "two", 3: "three"},
			i:             1,
			j:             0,
			want:          false,
		},
		"equal": {
			startingOrder: []int{1, 2, 3},
			startingMap:   map[int]string{1: "one", 2: "two", 3: "three"},
			i:             1,
			j:             1,
			want:          false,
		},
		"i_low": {
			startingOrder: []int{1, 2, 3},
			startingMap:   map[int]string{1: "one", 2: "two", 3: "three"},
			i:             -1,
			j:             1,
			want:          false,
		},
		"i_high": {
			startingOrder: []int{1, 2, 3},
			startingMap:   map[int]string{1: "one", 2: "two", 3: "three"},
			i:             3,
			j:             1,
			want:          false,
		},
		"j_low": {
			startingOrder: []int{1, 2, 3},
			startingMap:   map[int]string{1: "one", 2: "two", 3: "three"},
			i:             1,
			j:             -1,
			want:          false,
		},
		"j_high": {
			startingOrder: []int{1, 2, 3},
			startingMap:   map[int]string{1: "one", 2: "two", 3: "three"},
			i:             1,
			j:             3,
			want:          false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[int, string]{
				Map[int, string]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			less := m.Less(test.i, test.j)
			if less != test.want {
				t.Errorf("Unexpected less, wanted %t but got %t", test.want, less)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder []string
		startingMap   map[string]int
		key           string
		wantValue     int
		wantOK        bool
	}{
		"nil_load": {
			key:       "notakey",
			wantValue: 0,
			wantOK:    false,
		},
		"empty_load": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			key:           "notakey",
			wantValue:     0,
			wantOK:        false,
		},
		"load_none": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			key:           "notakey",
			wantValue:     0,
			wantOK:        false,
		},
		"load_seven": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			key:           "seven",
			wantValue:     7,
			wantOK:        true,
		},
		"load_zero": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			key:           "zero",
			wantValue:     0,
			wantOK:        true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			value, ok := m.Load(test.key)
			if value != test.wantValue {
				t.Errorf("Unexpected value, wanted %d but got %d", test.wantValue, value)
			}
			if ok != test.wantOK {
				t.Errorf("Unexpected OK, wanted %t but got %t", test.wantOK, ok)
			}
		})
	}
}

func TestLoadAndDelete(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder, wantOrder []string
		startingMap, wantMap     map[string]int
		key                      string
		wantValue                int
		wantLoaded               bool
	}{
		"nil_loadAndDelete": {
			key:        "notakey",
			wantValue:  0,
			wantLoaded: false,
		},
		"empty_loadAndDelete": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			key:           "notakey",
			wantValue:     0,
			wantLoaded:    false,
		},
		"loadAndDelete_none": {
			startingOrder: []string{"one"},
			startingMap:   map[string]int{"one": 1},
			wantOrder:     []string{"one"},
			wantMap:       map[string]int{"one": 1},
			key:           "notakey",
			wantValue:     0,
			wantLoaded:    false,
		},
		"loadAndDelete_one": {
			startingOrder: []string{"one"},
			startingMap:   map[string]int{"one": 1},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			key:           "one",
			wantValue:     1,
			wantLoaded:    true,
		},
		"loadAndDelete_five": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			key:           "five",
			wantValue:     5,
			wantLoaded:    true,
		},
		"loadAndDelete_zero": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			key:           "zero",
			wantValue:     0,
			wantLoaded:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			value, loaded := m.LoadAndDelete(test.key)
			if !reflect.DeepEqual(m.dirty, test.wantMap) {
				t.Errorf("Unexpected map content\nactual: %#v\nwant  : %#v", m.dirty, test.wantMap)
			}
			if !reflect.DeepEqual(m.order, test.wantOrder) {
				t.Errorf("Unexpected order content\nactual: %#v\nwant  : %#v", m.order, test.wantOrder)
			}
			if value != test.wantValue {
				t.Errorf("Unexpected value, wanted %d but got %d", test.wantValue, value)
			}
			if loaded != test.wantLoaded {
				t.Errorf("Unexpected OK, wanted %t but got %t", test.wantLoaded, loaded)
			}
		})
	}
}

func TestLoadAndDeleteFirst(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder, wantOrder []string
		startingMap, wantMap     map[string]int
		wantKey                  string
		wantValue                int
		wantLoaded               bool
	}{
		"nil_load": {
			wantKey:    "",
			wantValue:  0,
			wantLoaded: false,
		},
		"empty_load": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			wantKey:       "",
			wantValue:     0,
			wantLoaded:    false,
		},
		"loadAndDelete_zero": {
			startingOrder: []string{"zero"},
			startingMap:   map[string]int{"zero": 0},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			wantKey:       "zero",
			wantValue:     0,
			wantLoaded:    true,
		},
		"loadAndDelete_one": {
			startingOrder: []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantKey:       "one",
			wantValue:     1,
			wantLoaded:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			key, value, ok := m.LoadAndDeleteFirst()
			if key != test.wantKey {
				t.Errorf("Unexpected key, wanted %q but got %q", test.wantKey, key)
			}
			if value != test.wantValue {
				t.Errorf("Unexpected value, wanted %d but got %d", test.wantValue, value)
			}
			if ok != test.wantLoaded {
				t.Errorf("Unexpected OK, wanted %t but got %t", test.wantLoaded, ok)
			}
		})
	}
}

func TestLoadAndDeleteLast(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder, wantOrder []string
		startingMap, wantMap     map[string]int
		wantKey                  string
		wantValue                int
		wantLoaded               bool
	}{
		"nil_load": {
			wantKey:    "",
			wantValue:  0,
			wantLoaded: false,
		},
		"empty_load": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			wantKey:       "",
			wantValue:     0,
			wantLoaded:    false,
		},
		"loadAndDelete_zero": {
			startingOrder: []string{"zero"},
			startingMap:   map[string]int{"zero": 0},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			wantKey:       "zero",
			wantValue:     0,
			wantLoaded:    true,
		},
		"loadAndDelete_nine": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8},
			wantKey:       "nine",
			wantValue:     9,
			wantLoaded:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			key, value, ok := m.LoadAndDeleteLast()
			if key != test.wantKey {
				t.Errorf("Unexpected key, wanted %q but got %q", test.wantKey, key)
			}
			if value != test.wantValue {
				t.Errorf("Unexpected value, wanted %d but got %d", test.wantValue, value)
			}
			if ok != test.wantLoaded {
				t.Errorf("Unexpected OK, wanted %t but got %t", test.wantLoaded, ok)
			}
		})
	}
}

func TestLoadOrStore(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder, wantOrder []string
		startingMap, wantMap     map[string]int
		key                      string
		value                    int
		wantValue                int
		wantLoaded               bool
	}{
		"nil_Store": {
			wantOrder:  []string{"one"},
			wantMap:    map[string]int{"one": 1},
			key:        "one",
			value:      1,
			wantValue:  1,
			wantLoaded: false,
		},
		"empty_Store": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			wantOrder:     []string{"one"},
			wantMap:       map[string]int{"one": 1},
			key:           "one",
			value:         1,
			wantValue:     1,
			wantLoaded:    false,
		},
		"store_ten": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10},
			key:           "ten",
			value:         10,
			wantValue:     10,
			wantLoaded:    false,
		},
		"load_nine": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			key:           "nine",
			value:         10,
			wantValue:     9,
			wantLoaded:    true,
		},
		"load_zero": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			key:           "zero",
			value:         10,
			wantValue:     0,
			wantLoaded:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			value, ok := m.LoadOrStore(test.key, test.value)
			if !reflect.DeepEqual(m.dirty, test.wantMap) {
				t.Errorf("Unexpected map content\nactual: %#v\nwant  : %#v", m.dirty, test.wantMap)
			}
			if !reflect.DeepEqual(m.order, test.wantOrder) {
				t.Errorf("Unexpected order content\nactual: %#v\nwant  : %#v", m.order, test.wantOrder)
			}
			if value != test.wantValue {
				t.Errorf("Unexpected value, wanted %d but got %d", test.wantValue, value)
			}
			if ok != test.wantLoaded {
				t.Errorf("Unexpected OK, wanted %t but got %t", test.wantLoaded, ok)
			}
		})
	}
}

func TestRange(t *testing.T) {
	type row struct {
		key   string
		value int
	}
	for name, test := range map[string]struct {
		startingOrder []string
		startingMap   map[string]int
		endOn         int
		mutateOn      int
		wantRows      []row
	}{
		"nil_range": {
			wantRows: []row{},
		},
		"empty_range": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			wantRows:      []row{},
		},
		"range": {
			startingOrder: []string{"one", "two", "three"},
			startingMap:   map[string]int{"one": 1, "two": 2, "three": 3},
			wantRows:      []row{{"one", 1}, {"two", 2}, {"three", 3}},
		},
		"range_end_early": {
			startingOrder: []string{"one", "two", "three"},
			startingMap:   map[string]int{"one": 1, "two": 2, "three": 3},
			endOn:         1,
			wantRows:      []row{{"one", 1}, {"two", 2}},
		},
		"range_but_mutate": {
			startingOrder: []string{"one", "two", "three"},
			startingMap:   map[string]int{"one": 1, "two": 2, "three": 3},
			wantRows:      []row{{"one", 1}, {"two", 2}, {"three", 0}},
			mutateOn:      1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			s := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			var calls int64
			gotRows := make([]row, 0, len(test.wantRows))
			s.Range(func(index int, key string, value int) bool {
				atomic.AddInt64(&calls, 1)
				if index > cap(gotRows) {
					t.Errorf("Index %d over expected rows %d", index, cap(gotRows))
				}

				gotRows = append(gotRows, row{key, value})

				if test.mutateOn > 0 && test.mutateOn == index {
					s.order = []string{}
					s.dirty = map[string]int{}
				}

				if test.endOn > 0 && test.endOn == index {
					return false
				}
				return true
			})
			if !reflect.DeepEqual(gotRows, test.wantRows) {
				t.Errorf("Got unexpected rows\nactual: %#v\nwant  : %#v", gotRows, test.wantRows)
			}
		})
	}
}

func TestSort(t *testing.T) {
	s := SortMap[float64, string]{}
	for i := 0; i < 1000; i++ {
		n := rand.Float64()
		s.store(rand.Float64(), strconv.FormatFloat(n, 'E', -1, 64))
	}
	if sort.IsSorted(&s) {
		t.Error("you're extrodinarily unlucky")
	}
	t.Log("unsorted:", &s)
	sort.Sort(&s)
	t.Log("sorted  :", &s)
	if !sort.IsSorted(&s) {
		t.Error("It should be sorted")
	}
}

func TestStore(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder, wantOrder []string
		startingMap, wantMap     map[string]int
		key                      string
		value                    int
	}{
		"nil_insert": {
			wantOrder: []string{"one"},
			wantMap:   map[string]int{"one": 1},
			key:       "one",
			value:     1,
		},
		"empty_insert": {
			startingMap: map[string]int{},
			wantOrder:   []string{"one"},
			wantMap:     map[string]int{"one": 1},
			key:         "one",
			value:       1,
		},
		"nonempty_insert": {
			startingMap:   map[string]int{"one": 1},
			startingOrder: []string{"one"},
			wantOrder:     []string{"one", "two"},
			wantMap:       map[string]int{"one": 1, "two": 2},
			key:           "two",
			value:         2,
		},
		"duplicate_insert": {
			startingMap:   map[string]int{"one": 1, "two": 2},
			startingOrder: []string{"one", "two"},
			wantOrder:     []string{"one", "two"},
			wantMap:       map[string]int{"one": 1, "two": 2},
			key:           "one",
			value:         1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			m.Store(test.key, test.value)
			if !reflect.DeepEqual(m.dirty, test.wantMap) {
				t.Errorf("Unexpected map content\nactual: %#v\nwant  : %#v", m.dirty, test.wantMap)
			}
			if !reflect.DeepEqual(m.order, test.wantOrder) {
				t.Errorf("Unexpected order content\nactual: %#v\nwant  : %#v", m.order, test.wantOrder)
			}
		})
	}
}

func TestStoreFirst(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder, wantOrder []string
		startingMap, wantMap     map[string]int
		key                      string
		value                    int
	}{
		"nil_insert": {
			wantOrder: []string{"one"},
			wantMap:   map[string]int{"one": 1},
			key:       "one",
			value:     1,
		},
		"empty_insert": {
			startingMap: map[string]int{},
			wantOrder:   []string{"one"},
			wantMap:     map[string]int{"one": 1},
			key:         "one",
			value:       1,
		},
		"nonempty_insert": {
			startingMap:   map[string]int{"one": 1},
			startingOrder: []string{"one"},
			wantOrder:     []string{"two", "one"},
			wantMap:       map[string]int{"two": 2, "one": 1},
			key:           "two",
			value:         2,
		},
		"duplicate_insert": {
			startingMap:   map[string]int{"one": 1, "two": 2},
			startingOrder: []string{"one", "two"},
			wantOrder:     []string{"one", "two"},
			wantMap:       map[string]int{"one": 1, "two": 2},
			key:           "two",
			value:         2,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			m.StoreFirst(test.key, test.value)
			if !reflect.DeepEqual(m.dirty, test.wantMap) {
				t.Errorf("Unexpected map content\nactual: %#v\nwant  : %#v", m.dirty, test.wantMap)
			}
			if !reflect.DeepEqual(m.order, test.wantOrder) {
				t.Errorf("Unexpected order content\nactual: %#v\nwant  : %#v", m.order, test.wantOrder)
			}
		})
	}
}

func TestString(t *testing.T) {
	for name, test := range map[string]struct {
		object fmt.Stringer
		want   string
	}{
		"Map": {
			want: "github.com/brackendawson/ordered.Map[string,int][seven:7 nine:9 one:1]",
			object: &Map[string, int]{
				order: []string{"seven", "nine", "one"},
				dirty: map[string]int{"one": 1, "seven": 7, "nine": 9},
			},
		},
		"SortMap": {
			want: "github.com/brackendawson/ordered.SortMap[float64,bool][6.7:true 1:false -0.2:true]",
			object: &SortMap[float64, bool]{
				Map: Map[float64, bool]{
					order: []float64{6.7, 1, -0.2},
					dirty: map[float64]bool{6.7: true, 1: false, -0.2: true},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := test.object.String()
			if test.want != got {
				t.Errorf("Not equal:\n\twant: %s\n\tgot : %s", test.want, got)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	for name, test := range map[string]struct {
		startingOrder, wantOrder []string
		startingMap, wantMap     map[string]int
		i, j                     int
	}{
		"nil": {
			i: 1,
			j: 2,
		},
		"empty": {
			startingOrder: []string{},
			startingMap:   map[string]int{},
			wantOrder:     []string{},
			wantMap:       map[string]int{},
			i:             1,
			j:             2,
		},
		"one_two": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "two", "one", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             1,
			j:             2,
		},
		"zero_one": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"one", "zero", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             0,
			j:             1,
		},
		"eight_nine": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "nine", "eight"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             8,
			j:             9,
		},
		"one_zero": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"one", "zero", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             1,
			j:             0,
		},
		"nine_eight": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "nine", "eight"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             9,
			j:             8,
		},
		"one_one": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             1,
			j:             1,
		},
		"one_ten": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             1,
			j:             10,
		},
		"ten_one": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             10,
			j:             1,
		},
		"one_minus": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             1,
			j:             -1,
		},
		"minus_one": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             -1,
			j:             1,
		},
		"one_hundred": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             1,
			j:             100,
		},
		"two_negative_tree": {
			startingOrder: []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			startingMap:   map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			wantOrder:     []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
			wantMap:       map[string]int{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9},
			i:             2,
			j:             -3,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := SortMap[string, int]{
				Map[string, int]{
					order: test.startingOrder,
					dirty: test.startingMap,
				},
			}
			m.Swap(test.i, test.j)
			if !reflect.DeepEqual(m.dirty, test.wantMap) {
				t.Errorf("Unexpected map content\nactual: %#v\nwant  : %#v", m.dirty, test.wantMap)
			}
			if !reflect.DeepEqual(m.order, test.wantOrder) {
				t.Errorf("Unexpected order content\nactual: %#v\nwant  : %#v", m.order, test.wantOrder)
			}
		})
	}
}

// newLoadAndDelete is not an improvement
func (m *Map[K, V]) newLoadAndDelete(key K) (value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 0; i < len(m.order); i++ {
		if m.order[i] != key {
			continue
		}

		end := m.order[i+1:]
		m.order = m.order[:len(m.order)-1]
		copy(m.order[i:], end)
		break
	}

	value, loaded = m.dirty[key]
	delete(m.dirty, key)
	return
}

func BenchmarkLoadAndDelete(b *testing.B) {
	root := Map[int, string]{}
	for i := 0; i < 1000; i++ {
		root.Store(i, strconv.Itoa(i))
	}
	b.Run("old", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// b.StopTimer()
			m := Map[int, string]{}
			copy(m.order, root.order)
			// b.StartTimer()
			for i := 400; i < 600; i++ {
				m.LoadAndDelete(i)
			}
		}
	})
	b.Run("new", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// b.StopTimer()
			m := Map[int, string]{}
			copy(m.order, root.order)
			// b.StartTimer()
			for i := 400; i < 600; i++ {
				m.newLoadAndDelete(i)
			}
		}
	})
}
