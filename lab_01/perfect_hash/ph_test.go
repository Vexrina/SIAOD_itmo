package perfect_hash

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vexrina/siaod_itmo/lab_01/utils"
)

func hashFuncBruteforce(size, maxCounter int, ignoring []string) {
	seen := make(map[int]string)
	for i := 0; i < 1000000; i++ {
		key := fmt.Sprintf("%d", i)

		if slices.Contains(ignoring, key) {
			continue
		}

		hash := hashFunc(key, size)

		if prevKey, exists := seen[hash]; exists {
			fmt.Printf("Collision found! '%s' and '%s' both hash to %d\n", prevKey, key, hash)
			fmt.Println("increment size")
			size++
			ignoring = append(ignoring, prevKey, key)
			if len(ignoring) >= maxCounter {
				fmt.Println(ignoring)
				fmt.Println(size)
				return
			} else {
				hashFuncBruteforce(size, maxCounter, ignoring)
				return
			}
		} else {
			seen[hash] = key
		}
	}
}

//func TestTest(t *testing.T) {
//	hashFuncBruteforce(25, 4, []string{})
//}

func TestPerfectHash_HashFunc(t *testing.T) {
	//hashFuncBruteforce()
	tcs := []struct {
		key         string
		size        int
		expectedRes int
	}{
		{"0", 16, 15},
		{"1", 16, 12},
		{"2", 16, 5},
		{"3", 16, 2},
		{"4", 16, 3},
		{"5", 16, 0},
		{"6", 16, 9},
		{"7", 16, 6},
		{"8", 16, 7},
		{"9", 16, 4},
		{"10", 16, 4}, // collision found (size 16 -> keysSize = 4)
		{"0", 4, 3},
		{"1", 4, 0},
		{"2", 4, 1},
		{"3", 4, 2},
		{"4", 4, 3}, // collision found (size 4 -> keysSize = 2)
		{"some very diffrent key - 0", 9, 1},
		{"some very diffrent key - 9", 9, 1},
	}
	for _, tc := range tcs {
		t.Run(
			fmt.Sprintf("Key - %s, size - %d", tc.key, tc.size),
			func(t *testing.T) {
				actual := hashFunc(tc.key, tc.size)
				assert.Equal(t, tc.expectedRes, actual)
			},
		)
	}
}

func TestPerfectHash_NewPerfectHash(t *testing.T) {
	testcases := []struct {
		name            string
		keys            []string
		values          []any
		expectedIndexes []int
		expectedSize    int
	}{
		{
			name:            "basic test",
			keys:            []string{"1", "2", "3"},
			values:          []any{1, "2", 3.0},
			expectedIndexes: []int{1, 2, 4},
			expectedSize:    9,
		},
		{
			name:            "test with one collision",
			keys:            []string{"9", "10", "4", "3"},
			values:          []any{1, "2", 3.0, 2},
			expectedIndexes: []int{6, 13, 15, 16},
			expectedSize:    17,
		},
		{
			name:            "test with two collision",
			keys:            []string{"0", "4", "10", "15", "20"},
			values:          []any{1, "2", 3.0, 2, 5},
			expectedIndexes: []int{3, 9, 18, 21, 25},
			expectedSize:    27,
		},
	}

	for _, tc := range testcases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				ph := NewPerfectHash(tc.keys, tc.values)
				assert.Equal(t, tc.expectedSize, ph.size)
				//fmt.Println(ph.GetAllIndexes())
				for _, idx := range tc.expectedIndexes {
					assert.NotEqual(t, "", ph.table[idx])
				}
			},
		)
	}
}

func TestPerfectHash_Lookup(t *testing.T) {
	ph := NewPerfectHash([]string{"0", "4", "10", "15", "20"}, []any{1, "2", 3.0, 2, true})

	tcs := []struct {
		key         string
		expectedRes bool
	}{
		{"0", true},
		{"4", true},
		{"10", true},
		{"20", true},
		{"15", true},
		{"16", false},
		{"17", false},
		{"2153", false},
	}
	for _, tc := range tcs {
		t.Run(
			fmt.Sprintf("key - %s; exist - %v", tc.key, tc.expectedRes),
			func(t *testing.T) {
				actual := ph.Lookup(tc.key)
				assert.Equal(t, tc.expectedRes, actual)
			},
		)
	}
}

func TestPerfectHash_GetValueByKey(t *testing.T) {
	ph := NewPerfectHash([]string{"0", "4", "10", "15", "20"}, []any{1, "2", 3.0, 2, true})

	tcs := []struct {
		key         string
		expectedRes any
		expectedErr bool
	}{
		{"0", 1, false},
		{"4", "2", false},
		{"10", 3.0, false},
		{"20", true, false},
		{"15", 2, false},
		{"515", nil, true},
	}
	for _, tc := range tcs {
		t.Run(
			fmt.Sprintf("key - %s; exist - %v", tc.key, tc.expectedErr),
			func(t *testing.T) {
				actual, err := ph.GetValueByKey(tc.key)
				if !tc.expectedErr {
					require.NoError(t, err)
				}
				assert.Equal(t, tc.expectedRes, actual)
			},
		)
	}
}

func TestPerfectHash_GetAllSomething(t *testing.T) {
	ph := NewPerfectHash(
		[]string{"0", "4", "10", "15", "20"},
		[]any{1, "2", 3.0, 2, true},
	)
	tcs := []struct {
		funcName    string
		executable  func() any
		expectedRes any
	}{
		{
			funcName: "Keys",
			executable: func() any {
				return ph.GetAllKeys()
			},
			expectedRes: []string{"0", "4", "10", "15", "20"},
		},
		{
			funcName: "Values",
			executable: func() any {
				return ph.GetAllValues()
			},
			expectedRes: []any{1, "2", 3.0, 2, true},
		},
		{
			funcName: "KeyValuesValues",
			executable: func() any {
				return ph.GetAllKeysValues()
			},
			expectedRes: []*utils.KeyValue{
				{"0", 1},
				{"4", "2"},
				{"10", 3.0},
				{"15", 2},
				{"20", true},
			},
		},
		{
			funcName: "Indexes",
			executable: func() any {
				return ph.GetAllIndexes()
			},
			expectedRes: []any{3, 9, 18, 21, 25},
		},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("GetAll%s test", tc.funcName), func(t *testing.T) {
			actual := tc.executable()
			assert.ElementsMatch(t, tc.expectedRes, actual)
		})
	}
}

func TestPerfectHash_PutNewKeyValue(t *testing.T) {
	oldKeys := []string{"0", "4", "10", "15", "20"}
	oldValues := []any{1, "2", 3.0, 2, true}
	ph := NewPerfectHash(oldKeys, oldValues)
	newPh := ph.PutNewKeyValue("abracadbra", "blablabla")
	expectedIndexes := []int{0, 13, 15, 26, 33, 35}
	expectedSize := 37
	assert.Equal(t, expectedSize, newPh.size)
	for _, idx := range expectedIndexes {
		assert.NotEqual(t, "", newPh.table[idx])
	}
}

func generateKeys(n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}
	return keys
}

func generateValues(n int) []any {
	values := make([]any, n)
	for i := 0; i < n; i++ {
		values[i] = i
	}
	return values
}

var (
	keys   = generateKeys(1000)
	values = generateValues(1000)
	ph     = NewPerfectHash(keys, values)
)

// go test -benchmem -run=^$ -bench ^BenchmarkHashFunc
func BenchmarkHashFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashFunc("benchmark_key", 1024)
	}
}

func BenchmarkNewPerfectHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPerfectHash(keys, values)
	}
}

func BenchmarkLookupExisting(b *testing.B) {
	key := "key500" // Существующий ключ
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ph.Lookup(key)
	}
}

func BenchmarkLookupNonExisting(b *testing.B) {
	key := "non_existing_key" // Несуществующий ключ
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ph.Lookup(key)
	}
}

func BenchmarkGetValueByKey(b *testing.B) {
	key := "key500"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ph.GetValueByKey(key)
	}
}

func BenchmarkPutNewKeyValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ph.PutNewKeyValue(fmt.Sprintf("newKey%d", i), i)
	}
}
