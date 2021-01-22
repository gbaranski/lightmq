package packets

const (
	// TypeConnect - Client request to connect to Server
	//
	// Direction: Client to Server
	TypeConnect = iota + 1

	// TypeConnACK - Connect acknowledgment
	//
	// Direction: Server to Client
	TypeConnACK

	// TypePublish - Publish message
	//
	// Direction: Server to Client or Client to Server
	TypePublish

	// TypePubACK - Publish acknowledgment
	//
	// Direction: Server to Client or Client to Server
	TypePubACK

	// TypePubREC - Publish received (assured delivery part 1)
	//
	// Direction: Server to Client or Client to Server
	TypePubREC

	// TypePubREL - Publish release (assured delivery part 2)
	//
	// Direction: Server to Client or Client to Server
	TypePubREL

	// TypePubCOMP - Publish complete (assured delivery part 3)
	//
	// Direction: Server to Client or Client to Server
	TypePubCOMP

	// TypeSubscribe - Client subscribe request
	//
	// Direction: Client to Server
	TypeSubscribe

	// TypeSubACK - Subscribe acknowledgment
	//
	// Direction: Server to Client
	TypeSubACK

	// TypeUnsubscribe - Unsubscribe request
	//
	// Direction: Client to Server
	TypeUnsubscribe

	// TypeUnsubACK - Unsubscribe acknowledgment
	//
	// Direction: Server to Client
	TypeUnsubACK

	// TypePingREQ - Ping request
	//
	// Direction: Client to Server
	TypePingREQ

	// TypePingRESP - Ping response
	//
	// Direction: Server to Client
	TypePingRESP

	// TypeDisconnect - Client is disconnecting
	//
	// Direction: Client to Server
	TypeDisconnect
)
