package tools

import (
	"fmt"
	"math/rand"
	"time"
)

var _randInit bool = false

func randInit() {
	if !_randInit {
		rand.Seed(time.Now().UnixNano())
		_randInit = true
	}
}

// RandNum 随机数
func RandNum() int64 {
	randInit()
	return rand.Int63()
}

func RandNumString() string {
	return fmt.Sprintf("%d", RandNum())
}

// RandNumN 随机数
func RandNumN(n int64) int64 {
	randInit()
	return rand.Int63n(n)
}
