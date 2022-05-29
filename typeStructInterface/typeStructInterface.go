package typestructinterface

type ChatMessageContent struct {
	TimeStamp          string   `json:"timeStamp,omitempty"`
	Sender             string   `json:"sender,omitempty"`
	MessageRecipientId []string `json:"messageRecipientId,omitempty"`
	ChatRoomId         string   `json:"chatRoomId,omitempty"`
	MessageTextContent string   `json:"messageTextContent,omitempty"`
}
