package token

import "time"

type Creator interface {
	Create(username string, duration time.Duration) (string, error)
	Verify(token string) (*Payload, error)
}
