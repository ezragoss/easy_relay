package main

type RelayMessage struct {
	PeerID int32
	Packet []byte
}

func (h *Hub) SplitRelayMessage(message []byte) RelayMessage {
	var peerId int32 // 32-bit integer
	peerId |= int32(message[1])
	peerId |= int32(message[2])
	peerId |= int32(message[3])
	peerId |= int32(message[4])
	var packet = message[5:] // The packet to relay to the given
	return RelayMessage{
		peerId,
		packet,
	}
}

func (h *Hub) HandleRelayMessage(message RelayMessage) error {

	return nil
}
