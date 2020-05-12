package socketBox

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type serverBox struct {
	onlineConnections map[string]*Connection
}

var (
	wbsCon *websocket.Conn
	err    error
	data   []byte
	onlineConnections  = make(map[string]*Connection)
)

var upgrader = websocket.Upgrader{
	// 读取存储空间大小
	ReadBufferSize: 1024,
	// 写入存储空间大小
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		r.ParseForm()
		if r.Form["token"][0] == "" {
			return false
		}
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	if wbsCon, err = upgrader.Upgrade(w, r, nil); err != nil {
		return // 获取连接失败直接返回
	}
	myCon, _ := NewConnection(wbsCon,r.Form["token"][0])

	//heartbeat
	go func() {
		for  {
			fmt.Println(1)
			if err = myCon.WriteMsg([]byte("ping"));err!=nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	//listenning
	for {
		if data , err = myCon.ReadMsg();err != nil{
			fmt.Println("r err")
			goto ERR
		}
		if err = myCon.WriteMsg(data);err !=nil{
			fmt.Println("w err")
			goto ERR
		}
		//if _,data , err = myCon.conn.ReadMessage();err != nil{
		//	fmt.Println("r err")
		//	goto ERR
		//}
		//if err = myCon.conn.WriteMessage(websocket.TextMessage,data);err !=nil{
		//	fmt.Println("w err")
		//	goto ERR
		//}
	}

	ERR:
			myCon.Close()
}

func GetArr() map[string]*Connection {
	return onlineConnections
}

