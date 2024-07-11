package token

import "time"

type Maker interface {
	generateToken(username string, role string, duration time.Duration) (string, *Payload, error)
	verifyToken(token string) (*Payload, error)
}
