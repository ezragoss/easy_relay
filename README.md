# Websocket Relay Server

This is a websocket relay server with rooms for use in basic multiplayer game data relaying.

The structure of this relay server is specifically meant to work with Godot, but is agnostic of whatever data is being relayed so should work with any client that interfaces with it correctly.

## High Level

The lifetime of a connection looks like

Connect -> Populate Metadata -> Host/Join Match -> Relay -> Disconnect

## File Tour

### match.go

This is the logic for the 