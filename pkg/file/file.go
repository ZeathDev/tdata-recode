package file

import (
    "fmt"
    "io"
    "math/rand"
    "os"
    "path/filepath"
    "time"
)

// CopyFile 复制文件
func CopyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("无法打开源文件: %v", err)
    }
    defer sourceFile.Close()

    // 创建目标文件所在的目录
    if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
        return fmt.Errorf("创建目标目录失败: %v", err)
    }

    destFile, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("无法创建目标文件: %v", err)
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, sourceFile)
    if err != nil {
        return fmt.Errorf("复制文件失败: %v", err)
    }

    return nil
}

// GetFiles 获取指定目录下的所有文件
func GetFiles(path string) ([]string, error) {
    var files []string
    err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}

// CreateDirectory 创建目录
func CreateDirectory(path string) error {
    return os.MkdirAll(path, 0755)
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string {
    rand.Seed(time.Now().UnixNano())
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    result := make([]byte, length)
    for i := range result {
        result[i] = chars[rand.Intn(len(chars))]
    }
    return string(result)
}

// GetExecutablePath 获取可执行文件路径
func GetExecutablePath() (string, error) {
    ex, err := os.Executable()
    if err != nil {
        return "", err
    }
    return filepath.Dir(ex), nil
}

// IsDirectory 判断是否为目录
func IsDirectory(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    return info.IsDir()
}

// DeleteDirectory 删除目录
func DeleteDirectory(path string) error {
    return os.RemoveAll(path)
}

// DeleteFile 删除文件
func DeleteFile(path string) error {
    return os.Remove(path)
} 