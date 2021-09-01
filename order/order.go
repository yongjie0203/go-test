package order

type Order struct {
	OID   int64   `json:"oid" gorm:"primary_key;column:oid"`
	UID   int64   `json:"uid" gorm:"index:idx_uid"`
	Price float64 `json:"price"`
	//time added
	Time int64 `json:"time"`
	//buy or sell
	Direction int    `json:"direction"`
	TokenPair string `json:"token_pair"`
	//trade pair Id
	PairId     int     `json:"pair_id" gorm:"index:idx_pair_id"`
	Priority   float64 `json:"priority"`
	Num        float64 `json:"num"`
	Status     int     `json:"status"`
	UpdateTime int64
}

func (o Order) GetPriority() float64 {
	return o.Priority
}
