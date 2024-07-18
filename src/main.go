package main

import (
	"go-websocket-broadcast/src/core"
	"io"
	"net/http"
	"time"
	"go-websocket-broadcast/src/controllers"
	"go-websocket-broadcast/src/server"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	//配置初始化
	core.Config.Init()
	//日志初始化
	core.Config.Logger.Init()
	//客户端管理服务启动
	go server.Manager.Start()
	//关闭超时连接客户端
	go server.Manager.CloseTask()

	log.Infof("Server started %s ...", core.Config.Listen)
	http.HandleFunc("/ws", wsPage)
	// 配置TLS
	//err := http.ListenAndServeTLS(core.Config.Listen, core.Config.ServerCrt, core.Config.ServerKey,nil)
	//if err != nil {
	//	log.Fatal("ListenAndServeTLS error:", err)
	//}
	var pc controllers.PushController
	http.HandleFunc("/message/push", pc.Push)
	http.HandleFunc("/message/update_status", pc.UpdateReadStatus)
	log.Fatal(http.ListenAndServe(core.Config.Listen, nil))
}

//广播客户端连接handle
func wsPage(res http.ResponseWriter, req *http.Request) {
	//解析一个连接
	conn, error := upgrader.Upgrade(res, req, nil)
	if error != nil {
		io.WriteString(res, "这是一个websocket.")
		return
	}

	uid:= uuid.NewV4()
	sha1 := uid.String()

	//初始化一个客户端对象
	client := &server.Client{ID: sha1, Socket: conn, Send: make(chan []byte), RegisterTime: time.Now().Unix()}
	//注册一个对象到channel
	server.Manager.Register <- client

	go client.Read()
	go client.Write()
}
