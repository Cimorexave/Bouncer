package channels

// create channels class
type Channel struct {
	// define channels
	requestsChan chan mock_request
}

func (c *Channel) Init() {
	c.requestsChan = make(chan mock_request)
	return
}
