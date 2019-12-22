package quoter

import "github.com/stephenlyu/tds/entity"

type QuoterCallback interface {
	OnTickItem(tick *entity.TickItem)
	OnError(error)
}

type Quoter interface {
	Subscribe(security *entity.Security)
	SetCallback(callback QuoterCallback)

	Destroy()
}

type QuoteFactory interface {
	CreateQuoter(config interface{}) Quoter
}
