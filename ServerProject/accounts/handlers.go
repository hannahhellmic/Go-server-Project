package accounts

import (
	"awesomeProject/accounts/dto"
	"awesomeProject/accounts/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"sync"
)

func New() *Handler {
	return &Handler{
		accounts: make(map[string]*models.Account),
		guard:    &sync.RWMutex{},
	}
}

type Handler struct {
	accounts map[string]*models.Account
	guard    *sync.RWMutex
}

func (h *Handler) CreateAccount(c echo.Context) error {
	var request dto.ChangeAccountRequest // {"name": "alice", "amount": 50}
	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err)

		return c.String(http.StatusBadRequest, "invalid request")
	}

	if len(request.Name) == 0 {
		return c.String(http.StatusBadRequest, "empty name")
	}

	h.guard.Lock()

	if _, ok := h.accounts[request.Name]; ok {
		h.guard.Unlock()

		return c.String(http.StatusForbidden, "account already exists")
	}

	h.accounts[request.Name] = &models.Account{
		Name:   request.Name,
		Amount: request.Amount,
	}

	h.guard.Unlock()

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) GetAccount(c echo.Context) error {
	name := c.QueryParams().Get("name")

	h.guard.RLock()

	account, ok := h.accounts[name]

	h.guard.RUnlock()

	if !ok {
		return c.String(http.StatusNotFound, "account not found")
	}

	response := dto.GetAccountResponse{
		Name:   account.Name,
		Amount: account.Amount,
	}

	return c.JSON(http.StatusOK, response)
}

// Удаляет аккаунт
func (h *Handler) DeleteAccount(c echo.Context) error {
	name := c.QueryParams().Get("name")
	h.guard.Lock()
	if _, ok := h.accounts[name]; !ok {
		h.guard.Unlock()
		return c.String(http.StatusNotFound, "account not found")
	}
	delete(h.accounts, name)
	h.guard.Unlock()
	return c.NoContent(http.StatusNoContent)
}

// Меняет баланс
func (h *Handler) PathAccount(c echo.Context) error {
	var request dto.ChangeBalanceRequest
	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "invalid request")
	}

	h.guard.Lock()
	defer h.guard.Unlock()

	account, ok := h.accounts[request.Name]
	if !ok {
		return c.String(http.StatusNotFound, "account not found")
	}

	if request.SumChange < 0 && account.Amount < -1*request.SumChange {
		return c.String(http.StatusBadRequest, "account balance is too lower to withdraw the required amount of money")
	}

	account.Amount += request.SumChange
	return c.NoContent(http.StatusOK)
}

// Меняет имя
func (h *Handler) ChangeAccount(c echo.Context) error {
	name := c.QueryParams().Get("name")
	newName := c.QueryParams().Get("new_name")

	h.guard.Lock()
	defer h.guard.Unlock()

	if _, ok := h.accounts[name]; !ok {
		return c.String(http.StatusNotFound, "account not found")
	}

	h.accounts[newName] = h.accounts[name]
	h.accounts[newName].Name = newName
	delete(h.accounts, name)

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) TransferAccount(c echo.Context) error {
	var request dto.TransferAccountRequest
	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "invalid request")
	}

	h.guard.Lock()
	defer h.guard.Unlock()

	nameFrom, okFrom := h.accounts[request.NameFrom]
	if !okFrom {
		return c.String(http.StatusNotFound, "account not found")
	}

	nameTo, okTo := h.accounts[request.NameTo]
	if !okTo {
		return c.String(http.StatusNotFound, "account not found")
	}

	if request.Amount < 0 {
		return c.String(http.StatusBadRequest, "cannot transfer negative amount")
	}

	if nameFrom.Amount < request.Amount {
		return c.String(http.StatusBadRequest, "account balance is too low to transfer the required amount of money")
	}

	nameFrom.Amount -= request.Amount
	nameTo.Amount += request.Amount
	return c.NoContent(http.StatusOK)
}

func (h *Handler) ListAccounts(c echo.Context) error {
	h.guard.RLock()
	defer h.guard.RUnlock()

	allAccounts := make([]dto.GetAccountResponse, 0, len(h.accounts))
	for _, account := range h.accounts {
		allAccounts = append(allAccounts, dto.GetAccountResponse{
			Name:   account.Name,
			Amount: account.Amount,
		})
	}

	sort.Slice(allAccounts, func(i, j int) bool {
		return allAccounts[i].Amount < allAccounts[j].Amount
	})

	return c.JSON(http.StatusOK, allAccounts)
}
