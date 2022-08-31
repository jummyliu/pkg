package file

import (
	"os"
	"path/filepath"
)

// Write 写数据到指定的文件（指定权限），如：0644
//	如果指定文件的父目录不存在，则会创建
func Write(filename string, data []byte, perm os.FileMode) error {
	if err := DirCheckAndCreate(filename); err != nil {
		return err
	}
	os.WriteFile(filename, data, perm)
	return nil
}

// Merge 合并文件
//	TODO:
func Merge() {}

// GetFile 返回指定文件 *os.File，默认以 0644 打开文件
//	如果指定文件的父目录不存在，则会创建
//	文件不存在，则创建
//	文件存在，则追加内容
func GetFile(filename string) (*os.File, error) {
	return GetFileWithFlag(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE)
}

// GetFileWithFlag 以指定操作权限打开文件，并返回 *os.File，默认以 0644 打开文件
//	如果指定文件的父目录不存在，则会创建
func GetFileWithFlag(filename string, flag int) (*os.File, error) {
	if err := DirCheckAndCreate(filename); err != nil {
		return nil, err
	}
	return os.OpenFile(filename, flag, 0644)
}

// DirCheckAndCreate 检查给定路径的父目录是否存在，并创建
func DirCheckAndCreate(filename string) error {
	// check parent directory and create
	dir := filepath.Dir(filename)
	if file, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// dir is not exists
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	} else if !file.IsDir() {
		return err
	}
	return nil
}
