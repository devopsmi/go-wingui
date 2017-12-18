// Copyright (c) 2014 The cef2go authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cef2go

package main

import (
	"cef"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
	"wingui"

	ini "github.com/go-ini/ini"
)

var (
	Logger   *log.Logger = log.New(os.Stdout, "[main] ", log.Lshortfile)
	cfg, _               = ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, "app.ini")
	cmd      *exec.Cmd
	once     *sync.Once = &sync.Once{}
	localURL            = "127.0.0.1:51212"
)

func main() {
	hInstance, e := wingui.GetModuleHandle(nil)
	if e != nil {
		wingui.AbortErrNo("GetModuleHandle", e)
	}

	cef.ExecuteProcess(unsafe.Pointer(hInstance))

	settings := cef.Settings{}

	cef.Initialize(settings)

	title, _ := cfg.Section("").GetKey("app_title")

	wndproc := syscall.NewCallback(WndProc)
	Logger.Println("CreateWindow")
	hwnd := wingui.CreateWindow(title.String(), wndproc)

	browserSettings := cef.BrowserSettings{}

	url0, _ := cfg.Section("").GetKey("start_url")
	_url := url0.String()
	_url = "http://" + localURL + "/load?url=" + url.QueryEscape(_url)
	cef.CreateBrowser(unsafe.Pointer(hwnd), browserSettings, _url)

	// It should be enough to call WindowResized after 10ms,
	// though to be sure let's extend it to 500ms.
	time.AfterFunc(time.Millisecond*500, func() {
		cef.WindowResized(unsafe.Pointer(hwnd))
	})
	go execStart()
	cef.RunMessageLoop()
	cef.Shutdown()
	(*cmd).Process.Kill()
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
func execStart() {
	once.Do(func() {
		go run_http_server()
		argsString, _ := cfg.Section("").GetKey("start_exec")
		if argsString.String() == "" {
			return
		}
		args := strings.Split(argsString.String(), " ")
		//fmt.Printf("\nARGS:%v", args)
		if len(args) > 1 {
			cmd = exec.Command(args[0], args[1:]...)
		} else {
			cmd = exec.Command(args[0])
		}
		err := cmd.Start()
		if err != nil {
			fmt.Printf("ERR:%s", err)
			os.Exit(100)
		}
	})
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if strings.HasPrefix(req.URL.Path, "/load") {
		f, _ := os.OpenFile("start.html", os.O_RDONLY, 0666)
		c, _ := ioutil.ReadAll(f)
		fmt.Fprintf(w, string(c))
	} else if strings.HasPrefix(req.URL.Path, "/ping") {
		_url := req.URL.Query().Get("url")
		err := HTTPGet(_url, 5000)
		if err == nil {
			fmt.Fprintf(w, "1")
		}
		fmt.Println(err)
	}
}

func run_http_server() {
	http.HandleFunc("/", handler)
	listener, err := net.Listen("tcp", localURL)
	if err != nil {
		panic(err)
	}
	localURL += fmt.Sprint("%d", listener.Addr().(*net.TCPAddr).Port)
	fmt.Printf("\nUsing port:%d\n", listener.Addr().(*net.TCPAddr).Port)
	http.Serve(listener, nil)
}

func HTTPGet(URL string, timeout int) (err error) {
	tr := &http.Transport{}
	var resp *http.Response
	var client *http.Client
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		tr.CloseIdleConnections()
	}()
	client = &http.Client{Timeout: time.Millisecond * time.Duration(timeout), Transport: tr}
	resp, err = client.Get(URL)
	if err != nil {
		return
	}
	return
}
