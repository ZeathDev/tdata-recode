package network

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
)

// APIResponse API返回结构
type APIResponse struct {
    Code    int    `json:"code"`
    Msg     string `json:"msg"`
    Hash    string `json:"hash"`
    Name    string `json:"name"`
    Size    int    `json:"size"`
    Type    string `json:"type"`
    DownURL string `json:"downurl"`
    ViewURL string `json:"viewurl,omitempty"`
}

// UploadConfig 上传配置
type UploadConfig struct {
    FilePath string
    Show     bool   // 是否首页显示
    IsPwd    bool   // 是否设置密码
    Pwd      string // 下载密码
}

// DefaultUploadConfig 返回默认的上传配置
func DefaultUploadConfig(filePath string) UploadConfig {
    return UploadConfig{
        FilePath: filePath,
        Show:     true,
        IsPwd:    false,
        Pwd:      "",
    }
}

// UploadFile 上传文件到服务器
func UploadFile(config UploadConfig) error {
    // 打开文件
    file, err := os.Open(config.FilePath)
    if err != nil {
        return fmt.Errorf("打开文件失败: %v", err)
    }
    defer file.Close()

    // 创建multipart writer
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    // 创建文件表单字段
    part, err := writer.CreateFormFile("file", filepath.Base(config.FilePath))
    if err != nil {
        return fmt.Errorf("创建表单字段失败: %v", err)
    }

    // 复制文件内容
    if _, err = io.Copy(part, file); err != nil {
        return fmt.Errorf("复制文件内容失败: %v", err)
    }

    // 添加其他表单字段
    writer.WriteField("format", "json")
    if config.Show {
        writer.WriteField("show", "1")
    } else {
        writer.WriteField("show", "0")
    }
    if config.IsPwd {
        writer.WriteField("ispwd", "1")
        writer.WriteField("pwd", config.Pwd)
    } else {
        writer.WriteField("ispwd", "0")
    }

    // 关闭writer
    if err = writer.Close(); err != nil {
        return fmt.Errorf("关闭writer失败: %v", err)
    }

    // 创建请求
    req, err := http.NewRequest("POST", "https://wp.nakano.top/api.php", body)
    if err != nil {
        return fmt.Errorf("创建请求失败: %v", err)
    }

    // 设置请求头
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("发送请求失败: %v", err)
    }
    defer resp.Body.Close()

    // 解析响应
    var apiResp APIResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return fmt.Errorf("解析响应失败: %v", err)
    }

    // 检查上传状态
    if apiResp.Code != 0 {
        return fmt.Errorf("上传失败: %s", apiResp.Msg)
    }

    fmt.Printf("文件上传成功！\n下载地址: %s\n", apiResp.DownURL)
    return nil
} 