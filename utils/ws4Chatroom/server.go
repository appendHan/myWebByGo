package ws4Chatroom

type hub struct {
	//客户端列表
	clients map[*client]bool
	//广播
	broadcast chan []byte
	//注册
	register chan *client
	//注销
	unRegister chan *client
}

func NewAndRun() *hub {
	h := &hub{
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unRegister: make(chan *client),
		clients:    make(map[*client]bool),
	}
	go h.run()
	return h
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			//处理注册
			h.clients[c] = true
		case c := <-h.unRegister:
			//处理注销
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case message := <-h.broadcast:
			//处理广播
			for c := range h.clients {
				select {
				case c.send <- message:
				default:
					close(c.send)
					delete(h.clients, c)
				}

			}
		}
	}
}
