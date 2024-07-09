package dto

type GetAccountResponse struct {
	Name         string   `json:"name"`
	Amount       int      `json:"amount"`
	Transactions []string `json:"transactions"`
}

type GetAllAccountResponse struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}
