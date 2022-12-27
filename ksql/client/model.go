package client

type Response struct {
	Type          string        `json:"@type"`
	Streams       []Stream      `json:"streams"`
	CommandStatus CommandStatus `json:"commandStatus"`
	ErrorCode     int           `json:"error_code"`
	Message       string        `json:"message"`
	Entities      []interface{} `json:"entities"`
	StatementText string        `json:"statementText"`
}

type CommandStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Stream struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
}
