/** Server commands are messages sent to the server by a client for the purpose of altering the state on the server

Example:
- Player wants to create a match, they send a server command to create a match and be made the host
- Player wants to fill out or change their online metadata
- Player wants to join a match
- Player wants a list of matches to join

*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"log"
)

const (
	SET_PLAYER_METADATA = "set_metadata"
	HOST_MATCH          = "host_match"
	JOIN_MATCH          = "join_match"
	LEAVE_MATCH         = "leave_match"
	LIST_MATCHES        = "list_matches"
	SET_MATCH_METADATA  = "set_match_metadata"
)

type ServerCommand struct {
	Action string `json:"action"`
	Inputs json.RawMessage
}

type SetPlayerMetadata struct {
	Username string `json:"username"`
}

type HostMatch struct {
	Name string `json:"name"`
}

type SetMatchMetadata struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type JoinMatch struct {
	UUID string `json:"uuid"`
}

type LeaveMatch struct {
	UUID string `json:"uuid"`
}

type ClientDescription struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

type MatchDescription struct {
	Name string `json:"name"`
	Guid string `json:"guid"`
}

// ListMatches
// Todo: Do we want to filter anything on the server side vs the client side?
type ListMatches struct{}

func (h *Hub) HandleServerCommand(client *Client, jsonData []byte) error {
	log.Println("Handling Server Command...")
	//var action ServerCommand
	var obj map[string]json.RawMessage

	if err := json.Unmarshal(jsonData, &obj); err != nil {
		return err
	}

	var action string
	if err := json.Unmarshal(obj["action"], &action); err != nil {
		return err
	}

	log.Println("Unmarshalled Command...")

	switch action {
	case SET_PLAYER_METADATA:
		return h.HandleSetPlayerMetadata(client, jsonData)
	case HOST_MATCH:
		var name string
		json.Unmarshal(obj["name"], &name)
		return h.HandleHostMatch(client, HostMatch{Name: name})
	case JOIN_MATCH:
		var uid string
		json.Unmarshal(obj["uuid"], &uid)
		return h.HandleJoinMatch(client, JoinMatch{UUID: uid})
	case LIST_MATCHES:
		return h.HandleListMatches(client)
	case LEAVE_MATCH:
		var uid string
		json.Unmarshal(obj["uuid"], &uid)
		return h.HandleLeaveMatch(client, LeaveMatch{UUID: uid})
	case SET_MATCH_METADATA:
		return h.HandleSetMatchMetadata(client, jsonData)
	default:
		return nil
	}
}

func (h *Hub) HandleSetPlayerMetadata(client *Client, jsonData []byte) error {
	return nil
}

func (h *Hub) HandleHostMatch(client *Client, message HostMatch) error {
	log.Println("Host match requested...")

	guid, err := uuid.NewUUID()
	if err != nil {
		log.Println("Could not create UUID for new match")
		return err
	}

	log.Printf("HANDLING GUID %s", guid.String())

	var name = message.Name
	if name == "" {
		name = Generate(2, "_")
	}

	clients := make(map[string]*Client)
	clients[client.guid.String()] = client

	meta := MatchData{
		Guid: guid,
		Name: name,
	}

	match := &Match{
		meta:       meta,
		host:       client,
		clients:    clients,
		maxClients: 4, // TODO: We can parameterize this
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		end:        make(chan bool),
	}

	h.matches[guid.String()] = match
	h.matchByClient[client] = match

	go match.run()

	return nil
}

func (h *Hub) HandleJoinMatch(client *Client, match JoinMatch) error {
	log.Println("Join match requested...")

	log.Printf("Base64 string to decode %s", match.UUID)

	decodedIDBytes, err := base64.StdEncoding.DecodeString(match.UUID)
	if err != nil {
		log.Println("Error decoding UUID")
		return err
	}

	uid, err := uuid.FromBytes(decodedIDBytes)
	if err != nil {
		log.Println("Error parsing decoded bytes")
		return err
	}

	log.Printf("HANDLING GUID %s", uid.String())
	matchObj := h.matches[uid.String()]

	if matchObj == nil {
		msg := "Match does not exist"
		log.Println(msg)
		response := []byte{CONF_FAILED_JOIN}
		response = append(response, []byte(msg)...)
		client.send <- response
		return nil
	}

	if len(matchObj.clients) == matchObj.maxClients {
		msg := "Max clients reached"
		log.Println(msg)
		response := []byte{CONF_FAILED_JOIN}
		response = append(response, []byte(msg)...)
		client.send <- response
		return nil
	}

	for _, existingClient := range matchObj.clients {
		if packet, err := json.Marshal(ClientDescription{Username: existingClient.username, UUID: base64.StdEncoding.EncodeToString(existingClient.guid[:])}); err != nil {
			log.Println("Could not marshall the client description")
			return err
		} else {
			notify := []byte{RES_ID_PEER_CONNECTED}
			notify = append(notify, packet...)
			client.send <- notify
		}
	}

	matchObj.clients[client.guid.String()] = client
	h.matchByClient[client] = matchObj

	msg := "Match successfully joined"
	response := []byte{CONF_JOIN_MATCH}
	response = append(response, []byte(msg)...)
	client.send <- response

	if packet, err := json.Marshal(ClientDescription{Username: client.username, UUID: base64.StdEncoding.EncodeToString(client.guid[:])}); err != nil {
		log.Println("Could not marshall the client description")
		return err
	} else {
		notify := []byte{RES_ID_PEER_CONNECTED}
		notify = append(notify, packet...)
		matchObj.broadcast <- notify
	}
	return nil
}

func (h *Hub) HandleLeaveMatch(client *Client, match LeaveMatch) error {
	return nil
}

func (h *Hub) HandleListMatches(client *Client) error {
	log.Println("List matches requested...")
	matchListing := make([]MatchDescription, 0)
	for guid := range h.matches {
		match := h.matches[guid]
		matchData := match.meta
		matchListing = append(matchListing, MatchDescription{Name: matchData.Name, Guid: base64.StdEncoding.EncodeToString(matchData.Guid[:])})
	}
	if packet, err := json.Marshal(matchListing); err != nil {
		log.Println("Could not marshall match listing")
		return err
	} else {
		log.Println("Sending response back to client")
		// Send the listing as data to just the client with the proper identifier byte prefix
		response := []byte{RES_ID_COMMAND_RES}
		response = append(response, packet...)
		client.send <- response
	}

	return nil
}

func (h *Hub) HandleSetMatchMetadata(client *Client, jsonData []byte) error {
	return nil
}
