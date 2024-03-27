package controller

const (
	PrimaryName   = "primary"
	SecondaryName = "secondary"
)

type Router struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}
