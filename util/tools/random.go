package tools

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// UUID 生成UUID
// is是否去除-符号 默认去除
func UUID(is ...bool) string {
	if len(is) > 0 && is[0] {
		return fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil))
	}
	return strings.Replace(fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil)), "-", "", -1)
}

type random struct {
	rnd  *rand.Rand
	Lock sync.Mutex
}

func NewRandom() *random {
	return &random{rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// String 生成指定长度的随机字符串
func (r *random) String(n int, allowedChars ...[]rune) string {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	var letters []rune
	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}
	b := make([]rune, n)
	for i := range b {
		r.rnd.Seed(time.Now().UnixNano())
		b[i] = letters[r.rnd.Intn(len(letters))]
	}
	return string(b)
}

// Code 生成指定长度的随机数字
func (r *random) Code(length int) string {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	var container string
	for i := 0; i < length; i++ {
		r.rnd.Seed(time.Now().UnixNano())
		container += fmt.Sprintf("%01v", r.rnd.Int31n(10))
	}
	return container
}

// Num 生成指定范围内随机值
// [min,max)
func (r *random) Num(min, max int) int {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.rnd.Seed(time.Now().UnixNano())
	randNum := r.rnd.Intn(max-min) + min
	return randNum
}
