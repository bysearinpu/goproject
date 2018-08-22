package main;

import (
	"github.com/fsnotify/fsnotify"
	"fmt"
	"path/filepath"
	"os"
	"flag"
	"os/exec"
	"bytes"
	"errors"
	"regexp"
	"strconv"
)

type Watch struct {
	watch *fsnotify.Watcher;
}

//监控目录
func (w *Watch) watchDir(dir string,app string) {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path);
			if err != nil {
				return err;
			}
			err = w.watch.Add(path);
			if err != nil {
				return err;
			}
			fmt.Println("监控 : ", path);
		}
		return nil;
	});
	changeChan:=make(chan int)
	var changeFlag int = 0
	go func() {
		for {

			select {
			case ev := <-w.watch.Events:
				{
					m ,_:=regexp.MatchString("\\.exe-?" , ev.Name)
					m2 ,_:=regexp.MatchString("___jb_(old|tmp)___" , ev.Name)

					if ev.Op&fsnotify.Create == fsnotify.Create && !m && !m2{
						fmt.Println("创建文件 : ", ev.Name);
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(ev.Name);
						if err == nil && fi.IsDir() {
							w.watch.Add(ev.Name);
							fmt.Println("添加监控 : ", ev.Name);
						}
						changeFlag=1
					}
					if ev.Op&fsnotify.Write == fsnotify.Write && !m  && !m2{
						fmt.Println("写入文件 : ", ev.Name);
						changeFlag=1
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove && !m  && !m2{
						fmt.Println("删除文件 : ", ev.Name);
						//如果删除文件是目录，则移除监控
						fi, err := os.Stat(ev.Name);
						if err == nil && fi.IsDir() {
							w.watch.Remove(ev.Name);
							fmt.Println("删除监控 : ", ev.Name);
						}
						changeFlag=1
					}
				/*	if ev.Op&fsnotify.Rename == fsnotify.Rename && !m  && !m2{
						fmt.Println("重命名文件 : ", ev.Name);
						//如果重命名文件是目录，则移除监控
						//注意这里无法使用os.Stat来判断是否是目录了
						//因为重命名后，go已经无法找到原文件来获取信息了
						//所以这里就简单粗爆的直接remove好了
						w.watch.Remove(ev.Name);
						changeFlag=1
					}*/
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod && !m  && !m2{
						fmt.Println("修改权限 : ", ev.Name);
						changeFlag=1
					}
				}
			case err := <-w.watch.Errors:
				{
					fmt.Println("error : ", err);
					return;
				}
			}
			if(changeFlag==1){
				changeChan<-1
				changeFlag=0
			}
		}
	}();

	for{
		if(1==<-changeChan){
		res1,_:=findAndEnddApp(app)
		if(res1==0){
			fmt.Println("无此应用")
		}
		err:=startApp(app)
		if(err!=nil){
			fmt.Println(err)
		}
		fmt.Println("编译启动成功")
	   }
	}
}

func getPid2(processName string) (int, error) {
	//通过wmic process get name,processid | findstr server.exe获取进程ID
	buf := bytes.Buffer{};
	cmd := exec.Command("wmic", "process", "get", "name,processid");
	cmd.Stdout = &buf;
	cmd.Run();
	cmd2 := exec.Command("findstr", processName);
	cmd2.Stdin = &buf;
	data, _ := cmd2.CombinedOutput();
	if len(data) == 0 {
		return -1, errors.New("not find");
	}
	info := string(data);
	//这里通过正则把进程id提取出来
	reg := regexp.MustCompile(`[0-9]+`);
	pid := reg.FindString(info);
	return strconv.Atoi(pid);
}

func startApp(app string ) error{
	_,err :=exec.Command("go", "build","-o",app+".exe", app+".go").Output();
	if(err!=nil){
		return err
	}
	err2 :=exec.Command("./"+app+".exe").Start()
	if(err2!=nil){
		return err2
	}
	return nil
}

func findAndEnddApp(app string) (int,error){
	pid,_:= getPid2(app+".exe")
	process, err := os.FindProcess(pid);
	if err == nil {
		//让进程退出
		process.Kill();
	    return 1,nil
	}
	return 0,err
}

func main() {
    path :=flag.String("path","./","path need to fsnotify")
    app := flag.String("app","","app filename required")
    flag.Parse()
    if(*app==""){
    	panic("app filename required")
	}
	watch, _ := fsnotify.NewWatcher()
	w := Watch{
		watch: watch,
	}
	w.watchDir(*path,*app);
	select {};
}
