package main

import "github.com/google/uuid"

type MatchMessage struct {
	Action string     `json:"action"`
	Meta   *MatchData `json:"meta,omitempty"`
}

type MatchData struct {
	Guid uuid.UUID `json:"guid,omitempty"`
	Name string    `json:"name,omitempty"`
	//State   string `json:"state,omitempty"`
	//Private bool   `json:"private,omitempty"`
	//Key     string `json:"key,omitempty"`
}

// MatchStates
const (
	NotReady = "not-ready"
	READY    = "ready"
	ACTIVE   = "active"
	ENDED    = "ended"
)

type Match struct {
	host       *Client
	clients    map[string]*Client // guid -> client
	maxClients int

	meta MatchData

	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	end        chan bool
}

func (m *Match) run() {
	for {
		select {
		case broadcast := <-m.broadcast:
			for _, client := range m.clients {
				client.send <- broadcast
			}
		case end := <-m.end:
			if end {
				return
			}
		}
	}
}

//func (c *Client) HostMatch(name string) (*Match, error) {
//	guid, err := uuid.NewUUID()
//	if err != nil {
//		log.Println("Could not create UUID for new match")
//		return nil, err
//	}
//
//	if name == "" {
//		name = guid.String()
//	}
//
//	clients := make(map[string]*Client)
//	clients[c.guid.String()] = c
//
//	meta := MatchData{
//		Guid:  guid.String(),
//		Name:  name,
//		State: NotReady,
//	}
//
//	return &Match{
//		meta:          meta,
//		clients:       clients,
//		maxClients:    4, // TODO: We can parameterize this
//		register:      make(chan *Client),
//		unregister:    make(chan *Client),
//		metaBroadcast: make(chan *MatchMessage),
//		relay:         make(chan *string),
//	}, nil
//}

//func (m *Match) GetName() string {
//	return m.meta.Name
//}
//
//const (
//	PlayerJoined = "player_joined"
//	PlayerLeft   = "player_left"
//)
//
//func (m *Match) AttemptJoin(client *Client) error {
//	if false {
//		return errors.New("player already in match")
//	}
//
//	if m.meta.State != NotReady {
//		return errors.New("match has already begun")
//	}
//
//	if len(m.clients) == m.maxClients {
//		return errors.New("match is full")
//	}
//
//	m.register <- client
//	m.metaBroadcast <- &MatchMessage{
//		Action: PlayerJoined,
//		Meta:   &m.meta,
//	}
//
//	return nil
//}
//
//func (m *Match) OnPlayerLeave(client *Client) error {
//	m.unregister <- client
//	m.metaBroadcast <- &MatchMessage{
//		Action: PlayerLeft,
//		Meta:   &m.meta,
//	}
//
//	return nil
//}
