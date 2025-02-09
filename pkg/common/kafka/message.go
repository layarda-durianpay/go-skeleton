package kafka

type ResponseMessage[T any] struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	SubType string `json:"subtype"`
	Data    T      `json:"data"`
}
