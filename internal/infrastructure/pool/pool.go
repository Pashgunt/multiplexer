package pool

import (
	"sync"
)

var ConfigMapPool = &sync.Pool{
	New: func() interface{} { return make(map[string]interface{}) },
}

func PutConfigMapPool(configMap map[string]interface{}) {
	for key := range configMap {
		delete(configMap, key)
	}

	ConfigMapPool.Put(configMap)
}

var BinaryDataPool = &sync.Pool{
	New: func() interface{} { return make([]byte, 1024) },
}

func PutBinaryDataPool(binaryData []byte) {
	binaryData = binaryData[:0]
	BinaryDataPool.Put(binaryData)
}
