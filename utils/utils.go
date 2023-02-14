package utils

import (
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// UUID 封装获取 uuid 的函数
func UUID() string {
	return uuid.NewV4().String()
}

// RandomStr 随机字符串，包含大小写字母、数字、一般字符
//	字符集：abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~!@#$%^&*-+
func RandomStr(length int) string {
	// 字符集
	source := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~!@#$%^&*-+"
	return RandomStrWithSource(length, source)
}

// RandomStrWithSource 随机字符串，自行提供字符集
func RandomStrWithSource(length int, source string) string {
	sourceLen := len(source)
	if sourceLen == 0 {
		return ""
	}
	str := make([]byte, 0, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		str = append(str, source[r.Intn(sourceLen)])
	}
	return string(str)
}

// GetExecutablePath 获取执行文件路径，在 main 包下面调用，能获取到准确地路径
//	通过判断是否在临时目录，区分 go run 和 go build
//	go run:
// 		runtime.Caller(1) 获取调用者的文件路径
//
//	go build:
// 		os.Executable() 获取文件路径
func GetExecutablePath() string {
	file, err := os.Executable()
	if err != nil {
		return ""
	}
	// 通过判断是否在临时目录，区分 go run 和 go build
	if !strings.HasPrefix(file, os.TempDir()) {
		return filepath.Dir(file)
	}
	// 返回上一级调用者的地址
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Dir(file)
}

// FindIndex 返回元素在 slice 中所在的位置
//	如果未命中，则返回 slice 长度
func FindIndex[T comparable](arr []T, target T) int {
	for i, item := range arr {
		if target == item {
			return i
		}
	}
	return len(arr)
}

// CompareStringSlice 比较两个 []string
//	如果完全一致，则返回true，否则返回false
func CompareStringSlice(str1, str2 []string) bool {
	if len(str1) != len(str2) {
		return false
	}
	for i, str := range str1 {
		if str != str2[i] {
			return false
		}
	}
	return true
}

// CompareStringSlicePrefix 比较两个 []string，判断 str 中每一项是否包含 substr 中的每一项
func CompareStringSlicePrefix(str, substr []string) bool {
	if len(str) != len(substr) {
		return false
	}
	for i, str := range str {
		if !strings.HasPrefix(str, substr[i]) {
			return false
		}
	}
	return true
}

// FillSlice 使用指定元素填充 slice 到指定长度
func FillSlice[T any](slice []T, length int, fill T) []T {
	if len(slice) == length {
		return slice
	}
	newSlice := make([]T, 0, length)
	newSlice = append(newSlice, slice...)
	for i := len(slice); i < length; i++ {
		newSlice = append(newSlice, fill)
	}
	return newSlice
}

// Reverse 在原slice上进行修改，进行倒序
func Reverse[T any](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}
	l := len(arr)
	for i := 0; i < l/2; i++ {
		arr[i], arr[l-i-1] = arr[l-i-1], arr[i]
	}
	return arr
}
