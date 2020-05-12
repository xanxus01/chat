package main

import (
	"awesomeProject/socketBox"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)



var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {

			return true
		},
	}
	onlineCon []*websocket.Conn
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wbsCon *websocket.Conn
		err    error
		data   []byte
	)

	// 完成http应答，在httpheader中放下如下参数
	if wbsCon, err = upgrader.Upgrade(w, r, nil); err != nil {
		fmt.Println(err)
		return // 获取连接失败直接返回
	}

	onlineCon = append(onlineCon, wbsCon)

	for {
		// 只能发送Text, Binary 类型的数据,下划线意思是忽略这个变量.
		if _, data, err = wbsCon.ReadMessage(); err != nil {
			fmt.Println("读取消息错误")
			goto ERR // 跳转到关闭连接
		}
		fmt.Println(data)
		for _, v := range onlineCon {
			v.WriteMessage(websocket.TextMessage, data)
		}
		if err = wbsCon.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println("发送消息错误")
			goto ERR // 发送消息失败，关闭连接
		}
	}

ERR:
	// 关闭连接
	fmt.Println("关闭连接")
	wbsCon.Close()
}

func main() {
	// 当有请求访问ws时，执行此回调方法
	http.HandleFunc("/ws", socketBox.WsHandler)
	// 监听127.0.0.1:7777
	//go func() {
	//	for  {
	//		fmt.Println(socketBox.GetArr())
	//		time.Sleep(1*time.Second)
	//	}
	//}()
	err := http.ListenAndServe("0.0.0.0:7777", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}
