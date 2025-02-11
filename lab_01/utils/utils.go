package utils

type KeyValue struct {
	Key   string
	Value any
}

func NewKeyValue(key string, value any) *KeyValue {
	return &KeyValue{Key: key, Value: value}
}
