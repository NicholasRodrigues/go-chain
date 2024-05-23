package transactions

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/NicholasRodrigues/go-chain/pkg/crypto"
)

func TestNewTransaction(t *testing.T) {
	inputs := []TransactionInput{
		{Txid: []byte("somepreviousid"), Vout: 1, ScriptSig: "signature"},
	}
	outputs := []TransactionOutput{
		{Value: 10, ScriptPubKey: "pubkey1"},
	}
	tx := NewTransaction(inputs, outputs)

	if len(tx.ID) == 0 {
		t.Errorf("Expected transaction ID to be set")
	}
	if len(tx.Vin) != 1 {
		t.Errorf("Expected 1 input, got %d", len(tx.Vin))
	}
	if len(tx.Vout) != 1 {
		t.Errorf("Expected 1 output, got %d", len(tx.Vout))
	}
}

func TestTransaction_SerializeDeserialize(t *testing.T) {
	inputs := []TransactionInput{
		{Txid: []byte("somepreviousid"), Vout: 1, ScriptSig: "signature"},
	}
	outputs := []TransactionOutput{
		{Value: 10, ScriptPubKey: "pubkey1"},
	}
	tx := NewTransaction(inputs, outputs)

	serialized := tx.Serialize()
	deserialized := DeserializeTransaction(serialized)

	if !bytes.Equal(tx.ID, deserialized.ID) {
		t.Errorf("Expected transaction IDs to be equal, got %x and %x", tx.ID, deserialized.ID)
	}
	if len(deserialized.Vin) != len(tx.Vin) {
		t.Errorf("Expected %d inputs, got %d", len(tx.Vin), len(deserialized.Vin))
	}
	if len(deserialized.Vout) != len(tx.Vout) {
		t.Errorf("Expected %d outputs, got %d", len(tx.Vout), len(deserialized.Vout))
	}
}

func TestTransaction_SetID(t *testing.T) {
	inputs := []TransactionInput{
		{Txid: []byte("somepreviousid"), Vout: 1, ScriptSig: "signature"},
	}
	outputs := []TransactionOutput{
		{Value: 10, ScriptPubKey: "pubkey1"},
	}
	tx := NewTransaction(inputs, outputs)
	txIDBefore := tx.ID

	tx.SetID()
	if bytes.Equal(tx.ID, txIDBefore) {
		t.Errorf("Expected transaction ID to change after setting")
	}
}

func TestTransaction_SignAndVerify(t *testing.T) {
	privKey, err := crypto.NewPrivateKey()
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	inputs := []TransactionInput{
		{Txid: []byte("somepreviousid"), Vout: 1, ScriptSig: "signature"},
	}
	outputs := []TransactionOutput{
		{Value: 10, ScriptPubKey: "pubkey1"},
	}
	tx := NewTransaction(inputs, outputs)

	tx.Sign(privKey)

	for _, vin := range tx.Vin {
		pubKey, err := crypto.PublicKeyFromString(hex.EncodeToString(vin.PubKey))
		if err != nil {
			t.Fatalf("Failed to parse public key: %v", err)
		}
		if !pubKey.Verify(tx.Hash(), vin.Signature) {
			t.Errorf("Failed to verify transaction input signature")
		}
	}
}

func TestTransaction_Validate(t *testing.T) {
	privKey, err := crypto.NewPrivateKey()
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	// Example UTXO set
	utxoSet := make(map[string]TransactionOutput)
	utxoSet[utxoKey([]byte("somepreviousid"), 0)] = TransactionOutput{Value: 10, ScriptPubKey: "pubkey1"}

	inputs := []TransactionInput{
		{Txid: []byte("somepreviousid"), Vout: 0, ScriptSig: "signature"},
	}
	outputs := []TransactionOutput{
		{Value: 10, ScriptPubKey: "pubkey1"},
	}
	tx := NewTransaction(inputs, outputs)

	// Sign the transaction
	tx.Sign(privKey)

	// Validate the transaction
	isValid := tx.Validate(utxoSet)
	if !isValid {
		t.Errorf("Expected transaction to be valid")
	}

	// Modify the transaction to make it invalid
	tx.Vin[0].Txid = []byte("invalidid")
	isValid = tx.Validate(utxoSet)
	if isValid {
		t.Errorf("Expected transaction to be invalid")
	}
}

func TestTransaction_IsCoinbase(t *testing.T) {
	coinbaseTx := NewTransaction(
		[]TransactionInput{
			{Txid: []byte{}, Vout: -1, ScriptSig: "coinbase"},
		},
		[]TransactionOutput{
			{Value: 10, ScriptPubKey: "miner1"},
		},
	)

	if !coinbaseTx.IsCoinbase() {
		t.Errorf("Expected transaction to be a coinbase transaction")
	}

	regularTx := NewTransaction(
		[]TransactionInput{
			{Txid: []byte("somepreviousid"), Vout: 0, ScriptSig: "signature"},
		},
		[]TransactionOutput{
			{Value: 10, ScriptPubKey: "pubkey1"},
		},
	)

	if regularTx.IsCoinbase() {
		t.Errorf("Expected transaction to be a regular transaction")
	}
}
