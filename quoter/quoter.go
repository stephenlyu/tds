package quoter

import "github.com/stephenlyu/tds/entity"

type QuoterCallback interface {
	OnTickItem(tick *entity.TickItem)
}

type Quoter interface {
	Subscribe(security *entity.Security)
	SetCallback(callback QuoterCallback)

	Destroy()
}
