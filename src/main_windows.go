// Copyright (c) 2014 The cef2go authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cef2go

package main

import (
	"cef"
	"encoding/json"

	"net/url"
	"os"

	"syscall"
	"time"
	"unsafe"
	"wingui"
)

type Settings struct {
	SrvURL      string
	Width       int32
	Height      int32
	Title       string
	URL         string
	AppPid      int
	LauncherPid int
}

func main() {
	args := Settings{}
	buf := make([]byte, 1024)
	var err error
	var n int
	wait := make(chan bool, 1)
	go func() {
		n, err = os.Stdin.Read(buf)
		wait <- true
	}()
	select {
	case <-wait:
	case <-time.After(time.Second):
		os.Stdin.Close()
	}
	if err == nil {
		json.Unmarshal(buf[:n], &args)
	}
	hInstance, e := wingui.GetModuleHandle(nil)
	if e != nil {
		wingui.AbortErrNo("GetModuleHandle", e)
	}
	settings := cef.Settings{}
	cef.ExecuteProcess(unsafe.Pointer(hInstance))
	cef.Initialize(settings)
	wndproc := syscall.NewCallback(WndProc)
	hwnd := wingui.CreateWindow(args.Title, wndproc, args.Width, args.Height)
	browserSettings := cef.BrowserSettings{}
	_url0 := args.SrvURL + "/load?url=" + url.QueryEscape(args.URL)
	cef.CreateBrowser(unsafe.Pointer(hwnd), browserSettings, _url0)
	time.AfterFunc(time.Millisecond*500, func() {
		cef.WindowResized(unsafe.Pointer(hwnd))
	})
	cef.RunMessageLoop()
	cef.Shutdown()
	os.Exit(0)
}

func WndProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) (rc uintptr) {
	switch msg {
	case wingui.WM_CREATE:
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	case wingui.WM_SIZE:
		cef.WindowResized(unsafe.Pointer(hwnd))
	case wingui.WM_CLOSE:
		wingui.DestroyWindow(hwnd)
	case wingui.WM_DESTROY:
		cef.QuitMessageLoop()
	default:
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	}
	return
}
