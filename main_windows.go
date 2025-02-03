//go:build windows
// +build windows

package main

import (
    "github.com/lxn/win"
)

func init() {
    // 隐藏控制台窗口
    console := win.GetConsoleWindow()
    if console != 0 {
        win.ShowWindow(console, win.SW_HIDE)
    }
} 