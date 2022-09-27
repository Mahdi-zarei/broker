package modules

import "time"

type RowObject struct {
	Id         int
	Body       string
	Creation   time.Time
	Expiration int
	Subject    string
}

func NewRowObject(id int, body string, creation time.Time, expiration int, subject string) *RowObject {
	return &RowObject{
		Id:         id,
		Body:       body,
		Creation:   creation,
		Expiration: expiration,
		Subject:    subject,
	}
}
