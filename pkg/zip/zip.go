package zip

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

// CompressDirectory 将目录压缩成zip文件
func CompressDirectory(sourceDir, targetFile string) error {
    // 创建zip文件
    zipfile, err := os.Create(targetFile)
    if err != nil {
        return fmt.Errorf("创建zip文件失败: %v", err)
    }
    defer zipfile.Close()

    archive := zip.NewWriter(zipfile)
    defer archive.Close()

    // 遍历源目录
    err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // 获取文件头信息
        header, err := zip.FileInfoHeader(info)
        if err != nil {
            return fmt.Errorf("创建文件头信息失败: %v", err)
        }

        // 将绝对路径转换为相对路径
        relPath, err := filepath.Rel(sourceDir, path)
        if err != nil {
            return fmt.Errorf("获取相对路径失败: %v", err)
        }

        // 使用跨平台的路径分隔符
        header.Name = strings.ReplaceAll(relPath, "\\", "/")

        if info.IsDir() {
            header.Name += "/"
        } else {
            header.Method = zip.Deflate
        }

        // 创建压缩文件
        writer, err := archive.CreateHeader(header)
        if err != nil {
            return fmt.Errorf("创建压缩文件失败: %v", err)
        }

        if info.IsDir() {
            return nil
        }

        // 打开源文件
        file, err := os.Open(path)
        if err != nil {
            return fmt.Errorf("打开源文件失败: %v", err)
        }
        defer file.Close()

        // 复制文件内容到压缩文件
        _, err = io.Copy(writer, file)
        return err
    })

    return err
}

// DecompressFile 解压zip文件到指定目录
func DecompressFile(zipFile, targetDir string) error {
    reader, err := zip.OpenReader(zipFile)
    if err != nil {
        return fmt.Errorf("打开zip文件失败: %v", err)
    }
    defer reader.Close()

    // 创建目标目录
    if err := os.MkdirAll(targetDir, 0755); err != nil {
        return fmt.Errorf("创建目标目录失败: %v", err)
    }

    // 遍历压缩文件
    for _, file := range reader.File {
        path := filepath.Join(targetDir, file.Name)

        // 检查路径是否超出目标目录
        if !strings.HasPrefix(path, filepath.Clean(targetDir)+string(os.PathSeparator)) {
            return fmt.Errorf("非法的文件路径: %s", file.Name)
        }

        if file.FileInfo().IsDir() {
            os.MkdirAll(path, file.Mode())
            continue
        }

        // 创建目标文件所在的目录
        if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
            return fmt.Errorf("创建目录失败: %v", err)
        }

        // 创建目标文件
        targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
        if err != nil {
            return fmt.Errorf("创建目标文件失败: %v", err)
        }

        // 打开压缩文件
        sourceFile, err := file.Open()
        if err != nil {
            targetFile.Close()
            return fmt.Errorf("打开压缩文件失败: %v", err)
        }

        // 复制文件内容
        _, err = io.Copy(targetFile, sourceFile)
        targetFile.Close()
        sourceFile.Close()

        if err != nil {
            return fmt.Errorf("解压文件失败: %v", err)
        }
    }

    return nil
} 