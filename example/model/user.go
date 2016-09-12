package model

type GUID string

type VIPInfo struct {
	VIPLevel int64
	VIPType  string
}

type VIPExInfo struct {
	SumCost int64
}

type User struct {
	ID       int64
	Name     string
	Email    string
	Phone    string
	Address  string
	HomePage string
	Age      int64
	UUID     GUID
}
