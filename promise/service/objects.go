package promises

import (
	"github.com/satori/go.uuid"
)

type Status int
type Privacy int

const (
	StatusCreated Status = iota
	StatusPromised Status = iota
	StatusDone Status = iota
)

const (
	PrivacyPrivate Privacy = iota
	PrivacyFriends = iota
	PrivacyPublic = iota
)

type LEPromise struct{
	PromiseId int
	Title string
	Description string
	UserId uuid.UUID
	PromisedTo uuid.UUID
	Status Status
	Privacy Privacy
}