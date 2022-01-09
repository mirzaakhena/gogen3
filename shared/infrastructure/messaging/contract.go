package messaging

type Publisher interface {
	Publish(topic string, data []byte) error
}

type Subscriber interface {
	Handle(topic string, onReceived func(data []byte) error)
	Run(url string)
}
