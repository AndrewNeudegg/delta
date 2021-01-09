package http1

const (
	// ID for this collection of resources.
	ID = "http/v1"
)

type httpSinkServerResponse struct {
	ID     string `json:"id"`     // ID is the response ID for this accepted event.
	Reason string `json:"reason"` // Reason is why the response happened as it did.
	Status string `json:"status"` // Status states what happened to this event.
}
