# cef2win
�������Ǳ���õ�,����ֱ�ӱ����õ�cef����,����ʵ����ҳ����ı��ػ�.ʹ��html css js��д���س���.
# ���÷���   
����ĳ����ڲ�����cef��ʱ����Ŀ¼ΪcefĿ¼������Ϊcef.exe��Ȼ���ȡ���ý��̵ı�׼���룬Ȼ������׼����д��json�ַ����Ķ������ֽ����鼴�ɡ�   
json�ṹ���£����ִ�Сд��   
golang�ṹ��ʵ����   
```golang
type Settings struct {   
	SrvURL      string   
	Width       int32   
	Height      int32   
	Title       string   
	URL         string   
	AppPid      int   
	LauncherPid int   
}   
```
����˵����   
SrvURL�������ʼ��֮��򿪵���ҳurl��   
URL����ʼ��app��url   
        �����ʼ��֮�󣬻��SrvURL?url=URL����������Ŀ����ʹ��SrvURL���URL�Ľ����������Ѻõ���ʾ��   
        ��Ȼ����ֱ������SrvURLΪ��ڵ�ַ��URL���ա�   
Width�������ʼ��֮��Ŀ�ȣ���λ���ء�   
Height�������ʼ��֮��ĸ߶ȣ���λ���ء�   
Title���������   
AppPid������ʹ�ã���0���ɡ�   
LauncherPid������ʹ�ã���0���ɡ�   

go����ʾ����

```golang
        cefPath:="cef/cef.exe"
	settings := Settings{
		Title:       title,
		URL:         url0,
		SrvURL:      srvURL,
		Width:       int32(width),
		Height:      int32(height),
		AppPid:      cmdApp.Process.Pid,
		LauncherPid: os.Getpid(),
	}
	cmdCEF = exec.Command(cefPath)
	cmdCEF.Dir = filepath.Dir(cefPath)
	out, err := cmdCEF.StdinPipe()
	if err != nil {
		fmt.Printf("ERR:%s", err)
		os.Exit(1)
	}
	err = cmdCEF.Start()
	if err != nil {
		fmt.Printf("ERR:%s", err)
		os.Exit(1)
	}
	b, _ := json.Marshal(settings)
	_, err = out.Write(b)
	if err != nil {
		fmt.Printf("ERR:%s", err)
		os.Exit(1)
	}
```