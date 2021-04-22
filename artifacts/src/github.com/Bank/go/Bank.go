package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

type BankContract struct {
	contractapi.Contract
}

var logger = flogging.MustGetLogger("BankingLogger")
var contractType ContractType = Business

var contractPrefix string
var contractName string
var channelName string
var otherContractName string
var otherChannelName string

type ContractType uint8

const (
	Normal ContractType = iota
	Business
)

type Account struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Balance  uint64            `json:"balance"`
	Lock     map[string]string `json:"lock"`
	Received map[string]bool   `json:"received"`
}

type Transfer struct {
	SenderID    string `json:"senderID"`
	RecipientID string `json:"recipientID"`
	Amount      uint64 `json:"amount"`
	TX          string `json:"tx"`
}

func (s *BankContract) CreateAccount(ctx contractapi.TransactionContextInterface, user string) (string, error) {
	if len(user) == 0 {
		return "", fmt.Errorf("please pass the required args")
	}

	var account Account
	err := json.Unmarshal([]byte(user), &account)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling account. %s", err.Error())
	}
	if len(account.ID) == 0 || strings.Contains(account.ID, "~") {
		return "", fmt.Errorf("please pass a valid ID")
	}
	account.ID = contractPrefix + account.ID

	// TODO: check doesn't exists

	accountAsBytes, err := json.Marshal(account)
	if err != nil {
		return "", fmt.Errorf("failed while marshling account. %s", err.Error())
	}
	return account.ID, ctx.GetStub().PutState(account.ID, accountAsBytes)
}

func (s *BankContract) GetBalance(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	if len(id) == 0 {
		return "", fmt.Errorf("please pass the required args")
	}

	var accountAsByte []byte
	var err error
	if strings.HasPrefix(id, contractPrefix) {
		accountAsByte, err = ctx.GetStub().GetState(id)
	} else {
		params := []string{"GetBalance", id}
		queryArgs := make([][]byte, len(params))
		for i, arg := range params {
			queryArgs[i] = []byte(arg)
		}
		resp := ctx.GetStub().InvokeChaincode(otherContractName, queryArgs, otherChannelName)
		if resp.GetStatus() != 200 {
			return ctx.GetStub().GetTxID(), fmt.Errorf("failed to get the other channel %s method %s response with message %s", otherChannelName, otherContractName, resp.GetMessage())
		}
		accountAsByte = resp.GetPayload()
	}
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting account information %s", err.Error())
	}
	return string(accountAsByte), nil
}

func (s *BankContract) LockSenderMoney(ctx contractapi.TransactionContextInterface, transferData string) (string, error) {
	if len(transferData) == 0 {
		return "", fmt.Errorf("please pass the required args")
	}
	var transfer Transfer
	err := json.Unmarshal([]byte(transferData), &transfer)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}
	var senderAccountAsByte []byte

	if strings.HasPrefix(transfer.SenderID, contractPrefix) {
		senderAccountAsByte, err = ctx.GetStub().GetState(transfer.SenderID)
	} else {
		return ctx.GetStub().GetTxID(), fmt.Errorf("you should send the transaction to %s", otherChannelName)
	}
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting %s information %s", transfer.SenderID, err.Error())
	}
	var sender Account
	err = json.Unmarshal(senderAccountAsByte, &sender)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}
	if sender.Balance < transfer.Amount {
		return "", fmt.Errorf("failed insufficient credit for account %s", transfer.SenderID)
	}
	sender.Balance -= transfer.Amount
	sender.Lock = make(map[string]string)
	sender.Lock[ctx.GetStub().GetTxID()] = strings.Join([]string{transfer.RecipientID, strconv.FormatUint(transfer.Amount, 10)}, "~")

	senderAccountAsByte, err = json.Marshal(sender)
	if err != nil {
		return "", fmt.Errorf("failed while marshaling account: %s", err.Error())
	}
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(sender.ID, senderAccountAsByte)
}

func (s *BankContract) ReceiveMoney(ctx contractapi.TransactionContextInterface, transferData string) (string, error) {
	if len(transferData) == 0 {
		return "", fmt.Errorf("please pass the required args")
	}
	var transfer Transfer
	err := json.Unmarshal([]byte(transferData), &transfer)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}
	var recipientAccountAsByte []byte

	if strings.HasPrefix(transfer.RecipientID, contractPrefix) {
		recipientAccountAsByte, err = ctx.GetStub().GetState(transfer.RecipientID)
	} else {
		return ctx.GetStub().GetTxID(), fmt.Errorf("you should send the transaction to %s", otherChannelName)
	}
	var recipient Account
	err = json.Unmarshal(recipientAccountAsByte, &recipient)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}

	var senderAccountAsByte []byte

	if strings.HasPrefix(transfer.SenderID, contractPrefix) {
		senderAccountAsByte, err = ctx.GetStub().GetState(transfer.SenderID)
	} else {
		params := []string{"GetBalance", transfer.SenderID}
		queryArgs := make([][]byte, len(params))
		for i, arg := range params {
			queryArgs[i] = []byte(arg)
		}
		resp := ctx.GetStub().InvokeChaincode(otherContractName, queryArgs, otherChannelName)
		if resp.GetStatus() != 200 {
			return ctx.GetStub().GetTxID(), fmt.Errorf("failed to get the other channel response with message %s", resp.GetMessage())
		}
		senderAccountAsByte = resp.GetPayload()
	}
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting %s information %s", transfer.SenderID, err.Error())
	}

	var senderAccount Account
	err = json.Unmarshal(senderAccountAsByte, &senderAccount)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling %s. %s", transfer.SenderID, err.Error())
	}

	if strings.Compare(senderAccount.Lock[transfer.TX], strings.Join([]string{transfer.RecipientID, strconv.FormatUint(transfer.Amount, 10)}, "~")) != 0 {
		return "", fmt.Errorf("the sender money should get locked first")
	}
	recipient.Received = make(map[string]bool)
	recipient.Received[transfer.TX] = true
	recipient.Balance += transfer.Amount
	recipientAccountAsByte, err = json.Marshal(recipient)
	if err != nil {
		return "", fmt.Errorf("failed while marshaling account: %s", err.Error())
	}
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(recipient.ID, recipientAccountAsByte)
}

func (s *BankContract) FinializeTransfer(ctx contractapi.TransactionContextInterface, transferData string) (string, error) {
	if len(transferData) == 0 {
		return "", fmt.Errorf("please pass the required args")
	}
	var transfer Transfer
	err := json.Unmarshal([]byte(transferData), &transfer)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}
	var senderAccountAsByte []byte

	if strings.HasPrefix(transfer.SenderID, contractPrefix) {
		senderAccountAsByte, err = ctx.GetStub().GetState(transfer.SenderID)
	} else {
		return ctx.GetStub().GetTxID(), fmt.Errorf("you should send the transaction to %s", otherChannelName)
	}
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting %s information %s", transfer.SenderID, err.Error())
	}
	var sender Account
	err = json.Unmarshal(senderAccountAsByte, &sender)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}

	var recipientAccountAsByte []byte

	if strings.HasPrefix(transfer.RecipientID, contractPrefix) {
		recipientAccountAsByte, err = ctx.GetStub().GetState(transfer.RecipientID)
	} else {
		params := []string{"GetBalance", transfer.RecipientID}
		queryArgs := make([][]byte, len(params))
		for i, arg := range params {
			queryArgs[i] = []byte(arg)
		}
		resp := ctx.GetStub().InvokeChaincode(otherContractName, queryArgs, otherChannelName)
		if resp.GetStatus() != 200 {
			return ctx.GetStub().GetTxID(), fmt.Errorf("failed to get the other channel response with message %s", resp.GetMessage())
		}
		recipientAccountAsByte = resp.GetPayload()
	}
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting %s information %s", transfer.SenderID, err.Error())
	}

	var recipient Account
	err = json.Unmarshal(recipientAccountAsByte, &recipient)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling %s. %s", transfer.RecipientID, err.Error())
	}

	if ok := recipient.Received[transfer.TX]; !ok {
		return "", fmt.Errorf("the transfer should get confirm by recipient")
	}

	delete(sender.Lock, transfer.TX)
	senderAccountAsByte, err = json.Marshal(sender)
	if err != nil {
		return "", fmt.Errorf("failed while marshaling account: %s", err.Error())
	}
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(sender.ID, senderAccountAsByte)
}

func (s *BankContract) UnlockSenderMoney(ctx contractapi.TransactionContextInterface, transferData string) (string, error) {
	if len(transferData) == 0 {
		return "", fmt.Errorf("please pass the required args")
	}
	var transfer Transfer
	err := json.Unmarshal([]byte(transferData), &transfer)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}
	var senderAccountAsByte []byte

	if strings.HasPrefix(transfer.SenderID, contractPrefix) {
		senderAccountAsByte, err = ctx.GetStub().GetState(transfer.SenderID)
	} else {
		return ctx.GetStub().GetTxID(), fmt.Errorf("you should send the transaction to %s", otherChannelName)
	}
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting %s information %s", transfer.SenderID, err.Error())
	}
	var sender Account
	err = json.Unmarshal(senderAccountAsByte, &sender)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}

	var recipientAccountAsByte []byte

	if strings.HasPrefix(transfer.RecipientID, contractPrefix) {
		recipientAccountAsByte, err = ctx.GetStub().GetState(transfer.RecipientID)
	} else {
		params := []string{"GetBalance", transfer.RecipientID}
		queryArgs := make([][]byte, len(params))
		for i, arg := range params {
			queryArgs[i] = []byte(arg)
		}
		resp := ctx.GetStub().InvokeChaincode(otherContractName, queryArgs, otherChannelName)
		if resp.GetStatus() != 200 {
			return ctx.GetStub().GetTxID(), fmt.Errorf("failed to get the other channel response with message %s", resp.GetMessage())
		}
		recipientAccountAsByte = resp.GetPayload()
	}
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting %s information %s", transfer.SenderID, err.Error())
	}

	var recipient Account
	err = json.Unmarshal(recipientAccountAsByte, &recipient)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling %s. %s", transfer.RecipientID, err.Error())
	}

	if ok := recipient.Received[transfer.TX]; ok {
		return "", fmt.Errorf("the transfer should has been confirm by recipient")
	}

	amount, _ := strconv.ParseUint(strings.Split(sender.Lock[transfer.TX], "~")[0], 10, 64)
	delete(sender.Lock, transfer.TX)
	sender.Balance += amount
	senderAccountAsByte, err = json.Marshal(sender)
	if err != nil {
		return "", fmt.Errorf("failed while marshaling account: %s", err.Error())
	}
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(sender.ID, senderAccountAsByte)
}

// same channel transfer
func (s *BankContract) SendAmount(ctx contractapi.TransactionContextInterface, transferData string) (string, error) {
	if len(transferData) == 0 {
		return "", fmt.Errorf("please pass the required args")
	}
	var transfer Transfer
	err := json.Unmarshal([]byte(transferData), &transfer)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}
	if !strings.HasPrefix(transfer.RecipientID, contractPrefix) || !strings.HasPrefix(transfer.SenderID, contractPrefix) {
		return "", fmt.Errorf("you should use cross-channel transfer methods")
	}

	var senderAccountAsByte []byte

	senderAccountAsByte, err = ctx.GetStub().GetState(transfer.SenderID)
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting %s information %s", transfer.SenderID, err.Error())
	}
	var sender Account
	err = json.Unmarshal(senderAccountAsByte, &sender)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling transfer. %s", err.Error())
	}

	var recipientAccountAsByte []byte
	recipientAccountAsByte, err = ctx.GetStub().GetState(transfer.RecipientID)
	if err != nil {
		return ctx.GetStub().GetTxID(), fmt.Errorf("error while getting %s information %s", transfer.SenderID, err.Error())
	}
	var recipient Account
	err = json.Unmarshal(recipientAccountAsByte, &recipient)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshaling %s. %s", transfer.RecipientID, err.Error())
	}
	if sender.Balance < transfer.Amount {
		return "", fmt.Errorf("failed insufficient credit for account %s", transfer.SenderID)
	}
	sender.Balance -= transfer.Amount
	recipient.Balance += transfer.Amount
	senderAccountAsByte, err = json.Marshal(sender)
	if err != nil {
		return "", fmt.Errorf("failed while marshaling account: %s", err.Error())
	}
	recipientAccountAsByte, err = json.Marshal(recipient)
	if err != nil {
		return "", fmt.Errorf("failed while marshaling account: %s", err.Error())
	}
	err = ctx.GetStub().PutState(sender.ID, senderAccountAsByte)
	if err == nil {
		err = ctx.GetStub().PutState(recipient.ID, recipientAccountAsByte)
		return ctx.GetStub().GetTxID(), err
	}
	return ctx.GetStub().GetTxID(), err
}

func main() {

	if contractType == Normal {
		contractPrefix = "A"
		contractName = "banka"
		channelName = "bankachannel"
		otherContractName = "bankb"
		otherChannelName = "bankbchannel"
	} else {
		contractPrefix = "B"
		contractName = "bankb"
		channelName = "bankbchannel"
		otherContractName = "banka"
		otherChannelName = "bankachannel"
	}

	chaincode, err := contractapi.NewChaincode(new(BankContract))
	if err != nil {
		fmt.Printf("Error create %s chaincode: %s", contractName, err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting %s chaincode: %s", contractName, err.Error())
	}
}
