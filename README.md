# TData-recode

这是一个用Go语言编写的Telegram数据处理工具。

## 功能特性

- 通过进程名查找Telegram安装路径
- 自动收集和压缩指定的数据文件
- 自动上传文件到文件托管服务器
- 支持命令行参数配置
- 无窗口静默运行

## 环境配置

1. 安装Go语言环境：
   - 访问 [Go官网](https://golang.org/dl/) 下载并安装Go 1.16或更高版本
   - 设置GOPATH环境变量（可选）

2. 安装必要的工具：
```bash
# 安装goversioninfo工具（用于生成Windows资源文件）
go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest

# 安装依赖包
go get github.com/lxn/win
```

## 编译步骤

1. 克隆或下载项目代码：
```bash
git clone [项目地址] tg_go
cd tg_go
```

2. 生成Windows资源文件：
```bash
# 生成资源文件
goversioninfo -platform-specific=true
```

3. 编译程序：
```bash
# 使用GUI模式编译（无命令行窗口）
go build -ldflags -H=windowsgui

# 或使用普通模式编译（有命令行窗口，用于调试）
go build
```

## 命令行参数

- `-pn`: 进程名称（默认为"Telegram.exe"）
- `-sp`: 源目录路径（可选，如不指定则自动查找）
- `-dp`: 目标目录路径（可选，如不指定则使用程序所在目录）

## 使用示例

1. 使用默认配置：
```bash
tg_go.exe
```

2. 指定源目录：
```bash
tg_go.exe -sp "D:\Telegram Desktop\tdata"
```

3. 指定目标目录：
```bash
tg_go.exe -dp "D:\backup"
```

## 项目结构

- `main.go`: 主程序入口
- `main_windows.go`: Windows特定代码
- `pkg/process`: 进程管理相关功能
- `pkg/file`: 文件操作相关功能
- `pkg/network`: 网络上传相关功能
- `pkg/zip`: 压缩解压相关功能
- `winres.json`: Windows资源配置
- `versioninfo.json`: 程序版本信息

## 上传说明

程序会自动将压缩后的文件上传到文件托管服务器（https://wp.nakano.top）。上传成功后会显示：
- 上传状态信息
- 文件下载地址

## 注意事项

1. 编译环境要求：
   - Windows操作系统（推荐Windows 7或更高版本）
   - Go 1.16+
   - 管理员权限（用于访问进程信息）

2. 运行要求：
   - 程序需要访问Telegram的数据目录
   - 需要网络连接用于文件上传
   - 建议使用管理员权限运行

## 常见问题

1. 如果编译时提示缺少依赖，请运行：
```bash
go mod tidy
```

2. 如果运行时无法找到Telegram进程，可以手动指定源目录：
```bash
tg_go.exe -sp "D:\Telegram Desktop\tdata"
```

3. 如果需要调试程序，可以使用普通模式编译以查看命令行输出。 