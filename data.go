package tds


// After price adjusted, it maybe a negative value
type Record struct {
	Date uint32
	Open int32
	Close int32
	High int32
	Low int32
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
