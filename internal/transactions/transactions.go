package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/NicholasRodrigues/go-chain/pkg/crypto"
)

type Transaction struct {
	ID   []byte
	Vin  []TransactionInput
	Vout []TransactionOutput
}

type TransactionInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
	Signature []byte
	PubKey    []byte
}

type TransactionOutput struct {
	Value        int
	ScriptPubKey string
}

// NewTransaction creates a new transaction with the given inputs and outputs.
func NewTransaction(vin []TransactionInput, vout []TransactionOutput) *Transaction {
	tx := &Transaction{Vin: vin, Vout: vout}
	tx.SetID()
	return tx
}

// Sign signs each input of the transaction using the provided private key.
func (tx *Transaction) Sign(privKey *crypto.PrivateKey) {
	for i := range tx.Vin {
		msg := tx.Hash()
		tx.Vin[i].Signature = privKey.Sign(msg)
		tx.Vin[i].PubKey = privKey.PublicKey().Bytes()
	}
}

// Hash returns the hash of the transaction, excluding the signature to avoid circular dependencies.
func (tx *Transaction) Hash() []byte {
	txCopy := *tx
	for i := range txCopy.Vin {
		txCopy.Vin[i].Signature = nil
		txCopy.Vin[i].PubKey = nil
	}

	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(txCopy)
	if err != nil {
		panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

// SetID sets the ID of the transaction by hashing its contents.
func (tx *Transaction) SetID() {
	tx.ID = tx.Hash()
}

// IsCoinbase checks whether the transaction is a coinbase transaction.
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func utxoKey(txid []byte, vout int) string {
	return fmt.Sprintf("%x:%d", txid, vout)
}

// / Validate ensures that the transaction is valid.
func (tx *Transaction) Validate(utxoSet map[string]TransactionOutput) bool {
	if tx.IsCoinbase() {
		return true
	}

	inputValue := 0
	for _, vin := range tx.Vin {
		// Look up the referenced transaction output in the UTXO set
		key := utxoKey(vin.Txid, vin.Vout)
		utxo, ok := utxoSet[key]
		if !ok {
			fmt.Printf("Referenced output not found in UTXO set: %s\n", key)
			return false // Referenced output not found, transaction is invalid
		}

		// Verify the signature
		pubKey, err := crypto.PublicKeyFromString(hex.EncodeToString(vin.PubKey))
		if err != nil {
			fmt.Printf("Failed to parse public key: %v\n", err)
			return false
		}
		if !pubKey.Verify(tx.Hash(), vin.Signature) {
			fmt.Println("Signature is invalid")
			return false // Signature is invalid
		}

		inputValue += utxo.Value
	}

	outputValue := 0
	for _, vout := range tx.Vout {
		outputValue += vout.Value
	}

	fmt.Printf("Input value: %d, Output value: %d\n", inputValue, outputValue)

	// Check if input value matches output value
	return inputValue >= outputValue
}

// Serialize serializes the transaction into a byte slice.
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		panic(err)
	}
	return encoded.Bytes()
}

// DeserializeTransaction deserializes a transaction from a byte slice.
func DeserializeTransaction(data []byte) *Transaction {
	var transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		panic(err)
	}
	return &transaction
}

// String returns a human-readable representation of the transaction.
func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("---- Transaction %x:", tx.ID))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       ScriptSig: %s", input.ScriptSig))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:        %d", output.Value))
		lines = append(lines, fmt.Sprintf("       ScriptPubKey: %s", output.ScriptPubKey))
	}

	return strings.Join(lines, "\n")
}
