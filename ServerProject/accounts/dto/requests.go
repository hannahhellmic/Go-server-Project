package dto

type CreateAccountRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type PatchAccountRequest struct {
	Name string `json:"name"`
}

type ChangeAccountRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type DeleteAccountRequest struct {
	Name string `json:"name"`
}

type ChangeBalanceRequest struct {
	Name      string `json:"name"`
	SumChange int    `json:"sum_change"`
}

type TransferAccountRequest struct {
	NameFrom string `json:"name_from"`
	NameTo   string `json:"name_to"`
	Amount   int    `json:"amount"`
}
