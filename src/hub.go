package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

type RawMessage []byte
type Hub struct {
	// registered clients
	clients map[string]*Client
	// existing matches: Match GUID -> Match pointer
	matches map[string]*Match
	// A mapping of matches by client
	matchByClient map[*Client]*Match

	// channels
	broadcast chan struct {
		RawMessage
		*Client
	}
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	Action string `json:"action"`
	Data   json.RawMessage
}

const (
	SERVER_COMMAND = 0
	RELAY_MESSAGE  = 1
)

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[string]*Client),
		matches:       make(map[string]*Match),
		matchByClient: make(map[*Client]*Match),

		broadcast: make(chan struct {
			RawMessage
			*Client
		}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// FindMatchByName returns a match whose name matches the query name, returns nil if nothing found
func (h *Hub) FindMatchByName(name string) *Match {
	for guid := range h.matches {
		if h.matches[guid].meta.Name == name {
			return h.matches[guid]
		}
	}

	return nil
}

func (h *Hub) HandleRegistration(client *Client) {
	h.clients[client.guid.String()] = client
	log.Printf("Registering user with GUID %v, username %v", client.guid, client.username)

	notify := []byte{RES_ID_CONFIRMATION}
	notify = append(notify, CONF_CONNECTED)
	notify = append(notify, []byte(base64.StdEncoding.EncodeToString(client.guid[:]))...)
	client.send <- notify
}

func (h *Hub) HandleUnregistration(client *Client) {
	if _, c := h.clients[client.guid.String()]; c {
		delete(h.clients, client.guid.String())
		close(client.send)
	}
	log.Printf("Unregistering user with GUID %v, username %v", client.guid, client.username)
}

func ExtractAction(message []byte) (string, error) {
	var data map[string]json.RawMessage
	if err := json.Unmarshal(message, &data); err != nil {
		return "", err
	}

	var action string
	if err := json.Unmarshal(data["action"], &action); err != nil {
		return "", err
	}

	return action, nil
}

const (
	ACTION_MATCH = "match"
)

const (
	MatchCreate string = "match_create"
	MatchList          = "match_list"
	MatchJoin          = "match_join"
	MatchMove          = "match_move"
)

// HandleMatchMessage ingests a raw JSON message representing a MatchMessage and validates it, either creating the match or sending the message to the match channels
//func (h *Hub) HandleMatchMessage(message []byte, client *Client) error {
//	var matchMessage MatchMessage
//	if err := json.Unmarshal(message, &matchMessage); err != nil {
//		return err
//	}
//
//	switch matchMessage.Action {
//	case MatchCreate:
//		name := matchMessage.Meta.Name
//		if match, err := client.HostMatch(name); err != nil {
//			return err
//		} else {
//			h.matches[match.meta.Guid] = match
//		}
//	case MatchList:
//		matchListing := make([]MatchData, 0)
//		for guid := range h.matches {
//			match := h.matches[guid]
//			matchData := match.meta
//			matchListing = append(matchListing, matchData)
//		}
//		if packet, err := json.Marshal(matchListing); err != nil {
//			log.Println("Could not marshall match listing")
//			return err
//		} else {
//			// Send the listing as data to just the client
//			client.send <- packet
//		}
//	case MatchJoin:
//		for guid := range h.matches {
//			if guid == matchMessage.Meta.Guid || h.matches[guid].meta.Name == matchMessage.Meta.Name {
//				match := h.matches[guid]
//				match.clients[client.guid.String()] = client
//				h.matchByClient[client] = match
//			}
//		}
//	default:
//		return errors.New(fmt.Sprintf("Unrecognized action '%v'", matchMessage.Action))
//	}
//
//	log.Printf("Current match listing:")
//	for match := range h.matches {
//		log.Printf("\t%v", h.matches[match].meta.Name)
//		for player := range h.matches[match].clients {
//			log.Printf("\t\t%v %v", h.matches[match].clients[player].username, h.matches[match].clients[player].guid)
//		}
//	}
//
//	return nil
//}

func (h *Hub) HandleMessage(message []byte, client *Client) error {
	log.Println("Handling message...")
	var classifyingPrefix = message[0]
	log.Println("Analyzing prefix...")
	switch classifyingPrefix {
	case SERVER_COMMAND:
		log.Println("Handling server message")
		// Interpret the remainder of the packet as JSON
		return h.HandleServerCommand(client, message[1:])
	case RELAY_MESSAGE:
		log.Println("Handling relay message")
		// Structure the relay message struct
		relayMessage := h.SplitRelayMessage(message[1:], client)
		return h.HandleRelayMessage(relayMessage)
	default:
		log.Println("Classifying byte not recognized")
	}

	return nil
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.HandleRegistration(client)
		case client := <-h.unregister:
			h.HandleUnregistration(client)
		case packet := <-h.broadcast:
			message := packet.RawMessage
			client := packet.Client

			log.Println(message)

			if err := h.HandleMessage(message, client); err != nil {
				log.Println(err)
			}
		}
	}
}
