/** Server wide constants

 */

package main

/*
* These are used as prefixes for messages to the client
 */
const (
	RES_ID_RELAY_MSG         = byte(0)
	RES_ID_COMMAND_RES       = byte(1)
	RES_ID_CONFIRMATION      = byte(2)
	RES_ID_PEER_CONNECTED    = byte(3)
	RES_ID_PEER_DISCONNECTED = byte(4)
)

/*
* These are used to identify a type of confirmation response to a client
 */
const (
	CONF_JOIN_MATCH   = byte(0)
	CONF_FAILED_JOIN  = byte(1)
	CONF_HOSTED_MATCH = byte(2)
	CONF_FAILED_HOST  = byte(3)
	CONF_CONNECTED    = byte(4)
)

/*
* These are used to identify the prefix bit on a packet to know whether its a server command or relay message
 */
const (
	CMD_PREFIX   = byte(0)
	RELAY_PREFIX = byte(1)
)

/*
* These are used to identify the peer ID in relay packet data
 */
const (
	TARGET_PEER_BROADCAST = int32(0)
	TARGET_PEER_HOST      = int32(1)
)
