package ws

type Message struct {
	Event 		string 		`json:"event"`
	Author  	string 		`json:"author"`
	AuthorID	string		`json:"author_id"`
	Room  		string 		`json:"room"`
	Content 	string 		`json:"data"`
}

type messageHandler func(c *Client, msg Message)

type router struct {
	routes map[string]messageHandler
}

var Router = &router{
	routes: make(map[string]messageHandler),
}

func (r *router) Register(event string, handler messageHandler) {
	r.routes[event] = handler
}

func (r *router) Handle(c *Client, msg Message) {
	if fn, ok := r.routes[msg.Event]; ok {
		fn(c, msg)
	} else {
		return
	}
}

