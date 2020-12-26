package sink

// SunkMessage is a message that has been passed to the sink for forwarding
type SunkMessage struct {
	MessageID   *string
	Host        *string
	ContentType *string
	UserAgent   *string
	URI         *string
	Content     *[]byte
}
