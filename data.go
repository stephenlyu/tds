package tds


type Record struct {
	Date uint32
	Open uint32
	Close uint32
	High uint32
	Low uint32
	Volume float32
	Amount float32
}

type InfoExItem struct {
	Date uint32					`json:"date"`
	Bonus float32				`json:"bonus"`
	DeliveredShares float32		`json:"delivered_shares"`
	RationedSharePrice float32	`json:"rationed_share_price"`
	RationedShares float32		`json:"rationed_shares"`
}
