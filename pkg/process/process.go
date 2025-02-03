package process

import (
    "fmt"
    "syscall"
    "unsafe"
)

var (
    kernel32            = syscall.NewLazyDLL("kernel32.dll")
    psapi              = syscall.NewLazyDLL("psapi.dll")
    enumProcesses      = psapi.NewProc("EnumProcesses")
    openProcess        = kernel32.NewProc("OpenProcess")
    queryFullProcessImageName = kernel32.NewProc("QueryFullProcessImageNameW")
    closeHandle        = kernel32.NewProc("CloseHandle")
)

// ProcessInfo 存储进程信息
type ProcessInfo struct {
    ID   uint32
    Name string
    Path string
}

// GetProcessByName 根据进程名获取进程信息
func GetProcessByName(processName string) (*ProcessInfo, error) {
    const maxProcesses = 1024
    processes := make([]uint32, maxProcesses)
    var needed uint32

    r, _, err := enumProcesses.Call(
        uintptr(unsafe.Pointer(&processes[0])),
        uintptr(len(processes)*4),
        uintptr(unsafe.Pointer(&needed)),
    )

    if r == 0 {
        return nil, fmt.Errorf("EnumProcesses failed: %v", err)
    }

    numProcesses := needed / 4

    for i := uint32(0); i < numProcesses; i++ {
        if processPath := getProcessPath(processes[i]); processPath != "" {
            if getProcessNameFromPath(processPath) == processName {
                return &ProcessInfo{
                    ID:   processes[i],
                    Name: processName,
                    Path: processPath,
                }, nil
            }
        }
    }

    return nil, fmt.Errorf("process %s not found", processName)
}

// getProcessPath 获取进程完整路径
func getProcessPath(processID uint32) string {
    handle, _, _ := openProcess.Call(
        0x1000, // PROCESS_QUERY_LIMITED_INFORMATION
        0,
        uintptr(processID),
    )
    
    if handle == 0 {
        return ""
    }
    defer closeHandle.Call(handle)

    var size uint32 = 260 // MAX_PATH
    buffer := make([]uint16, size)
    r, _, _ := queryFullProcessImageName.Call(
        handle,
        0,
        uintptr(unsafe.Pointer(&buffer[0])),
        uintptr(unsafe.Pointer(&size)),
    )

    if r == 0 {
        return ""
    }

    return syscall.UTF16ToString(buffer[:size])
}

// getProcessNameFromPath 从路径中提取进程名
func getProcessNameFromPath(path string) string {
    for i := len(path) - 1; i >= 0; i-- {
        if path[i] == '\\' {
            return path[i+1:]
        }
    }
    return path
} 