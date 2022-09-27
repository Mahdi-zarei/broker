package modules

type Subscribers map[connection]struct{}

func NewSubscribersList() *Subscribers {
	return &Subscribers{}
}

func (s Subscribers) AddNewClient(clt connection) {
	s[clt] = struct{}{}
}
