package controllers

import (
	"chatRoom/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//设置参数
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//服务器
func Server(c *gin.Context) {
	//获取进入的房间
	roomid := c.Query("roomid")
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//把信息传入新进入聊天室的通道 进行广播
	con := &models.Connection{Send: make(chan []byte, 256), WsConn: wsConn}
	msg := models.Message{nil, roomid, c.GetString("username"), con}

	models.H.Join <- msg

	go msg.Write()
	go msg.Read()

}
