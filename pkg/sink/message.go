package sink

type SunkMessage struct {
	MessageID   *string
	Host        *string
	ContentType *string
	UserAgent   *string
	URI         *string
	Content     *[]byte
}
