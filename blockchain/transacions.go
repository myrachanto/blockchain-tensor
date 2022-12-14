package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type (
	Transaction struct {
		ID     []byte
		Inputs []TxInput
		Output []TxOutPut
	}
	TxOutPut struct {
		Value  int
		PubKey string
	}
	TxInput struct {
		Id  []byte
		Out int
		Sig string
	}
)

func CoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}
	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutPut{100, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutPut{txout}}
	tx.SetId()
	return &tx
}
func (tx *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutPut
	acc, validOutputs := chain.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("Error: not Enough funds")
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutPut{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutPut{acc - amount, from})
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetId()

	return &tx
}
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].Id) == 0 && tx.Inputs[0].Out == -1
}
func (in *TxInput) CanUnLock(data string) bool {
	return in.Sig == data
}
func (out *TxOutPut) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
