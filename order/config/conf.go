package config
//to do fix
type Stocker interface {
	ModifyProductStock(targetID uint32, num int) error
}

type UserChecker interface {
	UserCheck(userid uint64) error
}

type Config struct {
	OrderDB        string
	OrderTable     string
	ItemTable      string
	Stock          Stocker
	User           UserChecker
	ClosedInterval int
}
