package game

type Account struct {
	id                ClientId
	directReply       func(mesage string)
	reply             func(message string)
	broadcast         func(message string)
	loggedInCharacter *Character
}

func NewAccount(
	clientId ClientId,
	directReply func(mesage string),
	reply func(message string),
	broadcast func(message string),
) *Account {
	return &Account{
		id:                clientId,
		directReply:       directReply,
		reply:             reply,
		broadcast:         broadcast,
		loggedInCharacter: nil,
	}
}
