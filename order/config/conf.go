package config

type Stocker interface {
	ModifyProductStock(targetID string, num int)
}

type Config struct {
	OrderDB     string
	OrderTable  string
	ItemTable   string
	ModifyStock Stocker
}
