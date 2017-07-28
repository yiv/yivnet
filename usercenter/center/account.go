package center

type Account struct {
	Id          int64
	Uid         int64
	Unionid     string
	Uuid        string
	Username    string
	Password    string
	Nick        string
	Gender      bool
	Addr        string
	Avatar      string
	Isguest     bool
	Condays     int32
	Signdate    int64
	Vipsigndate int64
	Status      bool
	Mtime       int64
	Ctime       int64
	Bankpwd     string
	Forbid      string
	Imsi        string
	Imei        string
	Mac         string
	Did         string
	Psystem     string
	Pmodel      string
	Others      map[string]int32
}
