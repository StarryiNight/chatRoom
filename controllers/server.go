package controllers

import (
	"chatRoom/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Server(c *gin.Context) {
	roomid := c.Query("roomid")

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	con := &models.Connection{Send: make(chan []byte, 256), WsConn: wsConn}
	msg := models.Message{nil, roomid, c.GetString("username"), con}

	models.H.Join <- msg

	go msg.Write()
	go msg.Read()

}
