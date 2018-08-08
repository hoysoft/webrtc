package datachannel

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

/*
ChannelOpen represents a DATA_CHANNEL_OPEN Message

 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Message Type |  Channel Type |            Priority           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Reliability Parameter                      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         Label Length          |       Protocol Length         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
|                             Label                             |
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
|                            Protocol                           |
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/
type ChannelOpen struct {
	ChannelType          byte
	Priority             uint16
	ReliabilityParameter uint32

	Label    []byte
	Protocol []byte
}

const (
	channelOpenHeaderLength = 12

	offsetChannelType    = 1
	offsetPriority       = 2
	offsetReliability    = 4
	offsetLabelLength    = 8
	offsetProtocolLength = 10
)

// Marshal returns raw bytes for the given message
func (c *ChannelOpen) Marshal() ([]byte, error) {
	labelLength := len(c.Label)
	protocolLength := len(c.Protocol)

	totalLen := channelOpenHeaderLength + labelLength + protocolLength
	raw := make([]byte, totalLen)

	raw[offsetMessageType] = uint8(DataChannelOpen)
	raw[offsetChannelType] = c.ChannelType
	binary.BigEndian.PutUint16(raw[offsetPriority:], c.Priority)
	binary.BigEndian.PutUint32(raw[offsetReliability:], c.ReliabilityParameter)
	binary.BigEndian.PutUint16(raw[offsetLabelLength:], uint16(labelLength))
	binary.BigEndian.PutUint16(raw[offsetProtocolLength:], uint16(protocolLength))
	endLabel := channelOpenHeaderLength + labelLength
	copy(raw[channelOpenHeaderLength:endLabel], c.Label)
	copy(raw[endLabel:endLabel+protocolLength], c.Protocol)

	return raw, nil
}

// Unmarshal populates the struct with the given raw data
func (c *ChannelOpen) Unmarshal(raw []byte) error {
	if len(raw) < channelOpenHeaderLength {
		return errors.Errorf("Length of input is not long enough to satisfy header %d", len(raw))
	}
	c.ChannelType = raw[offsetChannelType]
	c.Priority = binary.BigEndian.Uint16(raw[offsetPriority:])
	c.ReliabilityParameter = binary.BigEndian.Uint32(raw[offsetReliability:])

	labelLength := binary.BigEndian.Uint16(raw[offsetLabelLength:])
	protocolLength := binary.BigEndian.Uint16(raw[offsetProtocolLength:])

	if len(raw) != int(channelOpenHeaderLength+labelLength+protocolLength) {
		return errors.Errorf("Label + Protocol length don't match full packet length")
	}

	c.Label = raw[channelOpenHeaderLength : channelOpenHeaderLength+labelLength]
	c.Protocol = raw[channelOpenHeaderLength+labelLength : channelOpenHeaderLength+labelLength+protocolLength]
	return nil
}
