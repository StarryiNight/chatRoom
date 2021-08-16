package models

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Msg      []byte
	Roomid   string
	Username string
	Conn     *Connection
}

type Connection struct {
	WsConn *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	Rooms     map[string]map[*Connection]bool
	Broadcast chan Message
	Join      chan Message
	Quit      chan Message
}

var H = Hub{
	Rooms:     make(map[string]map[*Connection]bool),
	Broadcast: make(chan Message),
	Join:      make(chan Message),
	Quit:      make(chan Message),
}

func (h *Hub) Run() {
	for true {
		select {
		//广播消息
		case m := <-h.Broadcast:
			conns := h.Rooms[m.Roomid]
			for con := range conns {
				if con == m.Conn {
					continue
				}
				select {
				case con.Send <- m.Msg:
				default:
					close(con.Send)
					delete(conns, con)
					if len(conns) == 0 {
						delete(h.Rooms, m.Roomid)
					}

				}
			}
			//进入聊天室
		case m := <-h.Join:
			conns := h.Rooms[m.Roomid]
			if conns == nil {
				conns = make(map[*Connection]bool)
				h.Rooms[m.Roomid] = conns
			}
			h.Rooms[m.Roomid][m.Conn] = true

			for con := range conns {
				str := "欢迎" + m.Username + "加入" + m.Roomid + "聊天室"
				msg := []byte(str)
				select {
				case con.Send <- msg:
				}

			}
			//退出聊天室
		case m := <-h.Quit:
			conns := h.Rooms[m.Roomid]
			if conns != nil {
				if _, ok := conns[m.Conn]; ok {
					delete(conns, m.Conn)
					close(m.Conn.Send)
					for con := range conns {
						str := m.Username + "离开了" + m.Roomid + "聊天室"
						msg := []byte(str)
						select {
						case con.Send <- msg:
						}
						if len(conns) == 0 {
							delete(h.Rooms, m.Roomid)
						}
					}
				}

			}
		}

	}
}

func (m Message) Read() {
	c := m.Conn

	defer func() {
		H.Quit <- m
		c.WsConn.Close()
	}()

	c.WsConn.SetReadLimit(512)
	//读写超时设置为60s
	c.WsConn.SetReadDeadline(time.Now().Add(60 * time.Second))
	//设置pong处理方式
	c.WsConn.SetPongHandler(func(string) error {
		c.WsConn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := c.WsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		//传入广播通道
		msg := Message{data, m.Roomid, m.Roomid, c}
		H.Broadcast <- msg
	}
}

func (m Message) Write() {
	c := m.Conn
	ticker := time.NewTicker(50 * time.Second)

	defer func() {
		ticker.Stop()
		c.WsConn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	c.WsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.WsConn.WriteMessage(mt, payload)
}
