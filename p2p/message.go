package p2p

type Message struct {
	Payload any
	From    string
}

type Person struct {
	Name string
	Age  int64
}
