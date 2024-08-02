package main

import (
	"encoding/base64"
	"github.com/google/uuid"
	"log"
)

type RelayMessage struct {
	PeerID uuid.UUID
	Packet []byte
}

func (h *Hub) SplitRelayMessage(message []byte, client *Client) RelayMessage {
	networkPeerID := string(message[:24]) // 24 bytes for the base64 encoded relay message
	uidBytes, err := base64.StdEncoding.DecodeString(networkPeerID)
	if err != nil {
		log.Println("Error decoding peerID")
		return RelayMessage{}
	}

	uid, err := uuid.FromBytes(uidBytes)
	if err != nil {
		log.Println("Error turning peer ID bytes into UUID")
		return RelayMessage{}
	}

	senderBytes := make([]byte, 24)
	base64.StdEncoding.Encode(senderBytes, client.guid[:])

	packet := senderBytes
	log.Println(message[24:])                // Prepend with the senders uid
	packet = append(packet, message[24:]...) // Append with the message itself
	return RelayMessage{
		uid,
		packet,
	}
}

func (h *Hub) HandleRelayMessage(message RelayMessage) error {
	client := h.clients[message.PeerID.String()]
	log.Println(client.username)
	packet := append([]byte{RES_ID_RELAY_MSG}, message.Packet...)
	client.send <- packet

	return nil
}
