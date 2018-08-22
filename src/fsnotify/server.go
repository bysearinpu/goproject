package main;

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"net"
	"fmt"
	"os"
	"os/signal"
)

const (
	confFilePath = "./conf/conf.json";
)

//我们这里只是演示，配置项只设置一个
type Conf struct {
	Port int `json:port`;
}

func main() {
	//读取文件内容
	data, err := ioutil.ReadFile(confFilePath);
	if err != nil {
		log.Fatal(err);
	}
	var c Conf;
	//解析配置文件
	err = json.Unmarshal(data, &c);
	if err != nil {
		log.Fatal(err);
	}
	//根据配置项来监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port));
	if err != nil {
		log.Fatal(err);
	}
	log.Println("server start");
	go func() {
		ch := make(chan os.Signal);
		//获取程序退出信号
		signal.Notify(ch, os.Interrupt, os.Kill);
		<-ch;
		log.Println("server exit");
		os.Exit(1);
	}();
	for {
		conn, err := lis.Accept();
		if err != nil {
			continue;
		}
		go func(conn net.Conn) {
			defer conn.Close();
			conn.Write([]byte("hello\n"));
		}(conn);
	}
}