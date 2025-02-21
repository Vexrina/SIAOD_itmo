package extendible_hashing

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtendableHash_CreateInsertGetStat(t *testing.T) {
	tcs := []struct {
		name            string
		keys            []string
		values          []string
		startDepth      int
		maxSize         int
		expectedDepth   int
		expectedNumDirs int
	}{
		{
			name:            "basic test",
			keys:            []string{"0", "1"},
			values:          []string{"0", "1"},
			startDepth:      1,
			maxSize:         1,
			expectedDepth:   1,
			expectedNumDirs: 2,
		},
		{
			name:            "split dirs test",
			keys:            []string{"apple", "banana", "grape", "orange", "watermelon"},
			values:          []string{"red", "yellow", "purple", "orange", "green"},
			startDepth:      1,
			maxSize:         2,
			expectedDepth:   3,
			expectedNumDirs: 8,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			eh := NewExtendableHash(tc.startDepth, tc.maxSize, false)
			for idx := range tc.keys {
				eh.Insert(tc.keys[idx], tc.values[idx])
			}
			actualDepth := eh.GetDepth()
			assert.Equal(t, tc.expectedDepth, actualDepth)
			actualNumDirs := eh.GetNumDirs()
			assert.Equal(t, 1<<tc.expectedDepth, actualNumDirs)
		})
	}
}

func TestExtendableHash_GetValues(t *testing.T) {
	eh := NewExtendableHash(3, 3, false)
	eh.Insert("apple", "red")
	eh.Insert("banana", "yellow")
	eh.Insert("grape", "purple")
	eh.Insert("orange", "orange")
	eh.Insert("watermelon", "green")
	eh.Insert("orange", "NOT ORANGE")
	var tcs = []struct {
		key           string
		foundFlag     bool
		expectedValue string
	}{
		{"apple", true, "red"},
		{"grape", true, "purple"},
		{"orange", true, "NOT ORANGE"},
		{"PC", false, ""},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%s  key is found? %v", tc.key, tc.foundFlag), func(t *testing.T) {
			val, found := eh.GetByKey(tc.key)
			assert.Equal(t, tc.foundFlag, found)
			if tc.foundFlag {
				assert.Equal(t, tc.expectedValue, val)
			} else {
				assert.Nil(t, val)
			}
		})
	}

	allKeys := eh.GetAllKeys()
	assert.ElementsMatch(t, allKeys, []string{"apple", "banana", "grape", "orange", "watermelon"})
}

func generateKeys(n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}
	return keys
}

func generateValues(n int) []string {
	values := make([]string, n)
	for i := 0; i < n; i++ {
		values[i] = fmt.Sprintf("value%d", i)
	}
	return values
}

var (
	keys        = generateKeys(10e6) // 10^6 = 1_000_000
	values      = generateValues(10e6)
	globalDepth = 20
	bucketSize  = 256
)

func BenchmarkNewExtendableHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewExtendableHash(globalDepth, bucketSize, false)
	}
}

func BenchmarkExtendableHash_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		eh := NewExtendableHash(globalDepth, bucketSize, false)
		for jndx := range keys {
			eh.Insert(keys[jndx], values[jndx])
		}
	}
}

func BenchmarkExtendableHash_GetByKey_ExistingKey(b *testing.B) {
	eh := NewExtendableHash(globalDepth, bucketSize, false)
	for jndx := range keys {
		eh.Insert(keys[jndx], values[jndx])
	}
	key := "key777"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh.GetByKey(key)
	}
}

func BenchmarkExtendableHash_GetByKey_NotExistingKey(b *testing.B) {
	eh := NewExtendableHash(globalDepth, bucketSize, false)
	for jndx := range keys {
		eh.Insert(keys[jndx], values[jndx])
	}
	key := "key-777"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh.GetByKey(key)
	}
}

func BenchmarkExtendableHash_GetByKey_Together(b *testing.B) {
	eh := NewExtendableHash(globalDepth, bucketSize, false)
	for jndx := range keys {
		eh.Insert(keys[jndx], values[jndx])
	}
	b.ResetTimer()
	for _, key := range keys {
		eh.GetByKey(key)
		eh.GetByKey("-" + key)
	}
}

func BenchmarkExtendableHash_InsertFileALL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		eh := NewExtendableHash(globalDepth, bucketSize, true)
		b.ResetTimer()
		for jndx := range keys {
			eh.Insert(keys[jndx], values[jndx])
		}
	}
}

func BenchmarkExtendableHash_InsertFileOne(b *testing.B) {
	eh := NewExtendableHash(globalDepth, bucketSize, true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh.Insert(keys[50], values[50])
	}
}

func BenchmarkExtendableHash_GetByKey_TogetherFILE(b *testing.B) {
	eh := NewExtendableHash(globalDepth, bucketSize, true)
	for jndx := range keys {
		eh.Insert(keys[jndx], values[jndx])
	}
	b.ResetTimer()
	for _, key := range keys {
		eh.GetByKey(key)
		eh.GetByKey("-" + key)
	}
}
