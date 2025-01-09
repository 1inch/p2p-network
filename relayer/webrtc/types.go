package webrtc

import (
	pb "github.com/1inch/p2p-network/proto"
)

// IncommingMessage incomming webrtc message.
type IncommingMessage struct {
	Request    *pb.ResolverRequest
	PublicKeys [][]byte
}

// OutcommingMessage outcomming webrtc message.
type OutcommingMessage struct {
	Response  *pb.ResolverResponse
	PublicKey []byte
}
