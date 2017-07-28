package center

type AccountInfo struct {
	Accounts_uid int64
	Token        string
	Coin         int64
	Gem          int32
	Bank         int64
	Growth       int32
	Level        int32
	Viptype      int32
	Vipexpiry    int64
	Voucher      int32
	Online       int32
	Props        map[string]int32
	Gifts        map[string]int32
	Medals       map[string]int32
	Friends      []int64
	Tags         []string
	Records      map[string]int32
}
