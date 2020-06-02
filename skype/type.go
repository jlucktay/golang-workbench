package main

type SkypeExport struct {
	UserID        string          `json:"userId"`
	ExportDate    string          `json:"exportDate"`
	Conversations []Conversations `json:"conversations"`
}

type Conversations struct {
	ID               string                 `json:"id"`
	DisplayName      string                 `json:"displayName"`
	Version          float64                `json:"version"`
	Properties       Properties             `json:"properties"`
	ThreadProperties map[string]interface{} `json:"threadProperties"`
	MessageList      []MessageList          `json:"MessageList"`
}

type Properties struct {
	Conversationblocked bool        `json:"conversationblocked"`
	Lastimreceivedtime  interface{} `json:"lastimreceivedtime"`
	Consumptionhorizon  string      `json:"consumptionhorizon"`
	Conversationstatus  interface{} `json:"conversationstatus"`
}

type MessageList struct {
	AMSReferences       []interface{}          `json:"amsreferences"`
	Content             string                 `json:"content"`
	ConversationID      string                 `json:"conversationid"`
	DisplayName         string                 `json:"displayName"`
	From                string                 `json:"from"`
	Id                  string                 `json:"id"`
	MessageType         string                 `json:"messagetype"`
	OriginalArrivalTime string                 `json:"originalarrivaltime"`
	Properties          map[string]interface{} `json:"properties"`
	Version             float64                `json:"version"`
}
