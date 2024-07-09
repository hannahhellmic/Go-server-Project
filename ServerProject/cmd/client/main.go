package main

import (
	"awesomeProject/accounts/dto"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type Command struct {
	Port      int
	Host      string
	Cmd       string
	Name      string
	NewName   string
	Amount    int
	NameFrom  string
	NameTo    string
	SumChange int
}

func (c *Command) Do() error {
	switch c.Cmd {
	case "create":
		return c.create()
	case "get":
		return c.get()
	case "delete":
		return c.delete()
	case "patch":
		return c.patch()
	case "change":
		return c.change()
	case "transfer":
		return c.transfer()
	case "get_all":
		return c.get_all()
	case "transactions":
		return c.transactions()
	default:
		return fmt.Errorf("unknown command: %s", c.Cmd)
	}
}

func (c *Command) create() error {
	request := dto.CreateAccountRequest{
		Name:   c.Name,
		Amount: c.Amount,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/account/create", c.Host, c.Port),
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("http post failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}

	return fmt.Errorf("resp error %s", string(body))
}

func (c *Command) get() error {
	resp, err := http.Get(
		fmt.Sprintf("http://%s:%d/account?name=%s", c.Host, c.Port, c.Name),
	)
	if err != nil {
		return fmt.Errorf("http get failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body failed: %w", err)
		}

		return fmt.Errorf("resp error %s", string(body))
	}

	var response dto.GetAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("json decode failed: %w", err)
	}

	fmt.Printf("response account name: %s and amount: %d\n", response.Name, response.Amount)
	return nil
}

func (c *Command) delete() error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s:%d/account?name=%s", c.Host, c.Port, c.Name), nil)
	if err != nil {
		return fmt.Errorf("http new request failed: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http delete failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}

	return fmt.Errorf("resp error %s", string(body))
}

func (c *Command) patch() error {
	request := dto.ChangeBalanceRequest{
		Name:      c.Name,
		SumChange: c.SumChange,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("http://%s:%d/account/changebalance", c.Host, c.Port), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("http new request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http patch failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}

	return fmt.Errorf("resp error %s", string(body))
}

func (c *Command) change() error {
	req, err := http.NewRequest("PATCH", fmt.Sprintf("http://%s:%d/account/change?name=%s&new_name=%s", c.Host, c.Port, c.Name, c.NewName), nil)
	if err != nil {
		return fmt.Errorf("http new request failed: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http change failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}

	return fmt.Errorf("resp error %s", string(body))
}

func (c *Command) transfer() error {
	request := dto.TransferAccountRequest{
		NameFrom: c.NameFrom,
		NameTo:   c.NameTo,
		Amount:   c.Amount,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("http://%s:%d/account/transfer", c.Host, c.Port), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("http new request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http transfer failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}

	return fmt.Errorf("resp error %s", string(body))
}

func (c *Command) get_all() error {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/account/all", c.Host, c.Port))
	if err != nil {
		return fmt.Errorf("http get failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body failed: %w", err)
		}
		return fmt.Errorf("resp error %s", string(body))
	}

	var response []dto.GetAllAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("json decode failed: %w", err)
	}

	for _, account := range response {
		fmt.Printf("account name: %s, amount: %d\n", account.Name, account.Amount)
	}

	return nil
}

func (c *Command) transactions() error {
	resp, err := http.Get(
		fmt.Sprintf("http://%s:%d/account/transactions?name=%s", c.Host, c.Port, c.Name),
	)
	if err != nil {
		return fmt.Errorf("http get failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body failed: %w", err)
		}
		return fmt.Errorf("resp error %s", string(body))
	}

	var responce dto.GetAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&responce); err != nil {
		return fmt.Errorf("json decode failed: %w", err)
	}

	for _, transaction := range responce.Transactions {
		fmt.Printf("%s: %s\n", responce.Name, transaction)
	}

	return nil
}

func main() {
	portVal := flag.Int("port", 1323, "server port")
	hostVal := flag.String("host", "0.0.0.0", "server host")
	cmdVal := flag.String("cmd", "", "command to execute")
	nameVal := flag.String("name", "", "name of account")
	newNameVal := flag.String("new_name", "", "new name of account")
	amountVal := flag.Int("amount", 0, "amount of account")
	nameFromVal := flag.String("name_from", "", "name of account")
	nameToVal := flag.String("name_to", "", "new name of account")
	SumChangeVal := flag.Int("sum_change", 0, "sum change of account")

	flag.Parse()

	cmd := Command{
		Port:      *portVal,
		Host:      *hostVal,
		Cmd:       *cmdVal,
		Name:      *nameVal,
		NewName:   *newNameVal,
		Amount:    *amountVal,
		NameFrom:  *nameFromVal,
		NameTo:    *nameToVal,
		SumChange: *SumChangeVal,
	}

	if err := cmd.Do(); err != nil {
		panic(err)
	}
}
