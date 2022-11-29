package client

type Response []struct {
	Streams       []Stream      `json:"streams"`
	CommandStatus CommandStatus `json:"commandStatus"`
	ErrorCode     int           `json:"error_code"`
	Message       string        `json:"message"`
}

type CommandStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Stream struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
}
