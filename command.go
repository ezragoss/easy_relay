/** Server commands are messages sent to the server by a client for the purpose of altering the state on the server

Example:
- Player wants to create a match, they send a server command to create a match and be made the host
- Player wants to fill out or change their online metadata
- Player wants to join a match
- Player wants a list of matches to join

*/

package main

import (
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
	Name string `json:"name"`
}

type JoinMatch struct {
	UUID string `json:"uuid"`
}

type LeaveMatch struct {
	UUID string `json:"uuid"`
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
		return h.HandleJoinMatch(client, jsonData)
	case LIST_MATCHES:
		return h.HandleListMatches(client)
	case LEAVE_MATCH:
		return h.HandleLeaveMatch(client, jsonData)
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

	var name = message.Name
	if name == "" {
		name = Generate(2, "_")
	}

	clients := make(map[string]*Client)
	clients[client.guid.String()] = client

	meta := MatchData{
		Guid: guid.String(),
		Name: name,
	}

	match := &Match{
		meta:       meta,
		host:       client,
		clients:    make(map[string]*Client),
		maxClients: 4, // TODO: We can parameterize this
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *string),
	}

	h.matches[guid.String()] = match
	h.matchByClient[client] = match

	return nil
}

func (h *Hub) HandleJoinMatch(client *Client, jsonData []byte) error {
	//
	//if err:= json.Unmarshal(jsonData, )

	return nil
}

func (h *Hub) HandleLeaveMatch(client *Client, jsonData []byte) error {
	return nil
}

func (h *Hub) HandleListMatches(client *Client) error {
	log.Println("List matches requested...")
	matchListing := make([]MatchData, 0)
	for guid := range h.matches {
		match := h.matches[guid]
		matchData := match.meta
		matchListing = append(matchListing, matchData)
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
