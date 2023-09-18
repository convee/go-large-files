package pkgs

import (
	"os"
	"path/filepath"
	"time"
)

func GetFolderPath(genDateSubFolder bool, subDirs ...string) string {
	if genDateSubFolder {
		subDirs = append(subDirs, time.Now().Format("20060102"))
	}
	folderPath := filepath.Join(subDirs...)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0777) //0777也可以os.ModePerm
	}
	return folderPath
}

// RemoveFile 移除文件
func RemoveFile(path string) bool {
	if FileExists(path) {
		err := os.Remove(path)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

// FileExists 判断所给路径文件/文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
