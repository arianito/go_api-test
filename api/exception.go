package api


type Exception struct {
	Status     int               `json:"status"`
	Message    string            `json:"message"`
	Errors     []interface{}     `json:"errors"`
	Validation map[string]string `json:"validation"`
	Stack      string            `json:"stack"`
}