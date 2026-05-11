package types

import (
	"encoding/hex"
	fmt "fmt"
)

const (
	BtcAddressMinLen = 26
	BtcAddressMaxLen = 62
)

type BtcAddress struct {
	BtcAddress string
}

type BtcScript struct {
	BtcScript string
}

func (ea BtcAddress) GetBtcAddress() string {
	return ea.BtcAddress
}

func NewBtcAddress(Address string) (*BtcAddress, error) {
	if err := ValidateBtcAddress(Address); err != nil {
		return nil, fmt.Errorf("invalid input Address: %w", err)
	}
	addr := BtcAddress{Address}
	return &addr, nil
}

func ValidateBtcAddress(Address string) error {
	if Address == "" {
		return fmt.Errorf("empty")
	}
	if len(Address) < BtcAddressMinLen && len(Address) > BtcAddressMaxLen {
		return fmt.Errorf("address (%s) of the wrong length expected between (%d) and (%d) actual(%d)", Address, BtcAddressMinLen, BtcAddressMaxLen, len(Address))
	}
	return nil
}

func IsValidSignature(signature string) bool {
	signatureLen := len(signature)
	if (signatureLen < 140 || signatureLen > 144) && signatureLen != 128 {
		return false
	}
	_, err := hex.DecodeString(signature)
	return err == nil
}

func ValidateSignatures(signatures []string) bool {
	for _, signature := range signatures {
		if !IsValidSignature(signature) {
			return false
		}
	}
	return true
}

func NewBtcScript(script string) (*BtcScript, error) {
	return &BtcScript{script}, nil
}

func (ea BtcScript) GetBtcReserveScript() string {
	return ea.BtcScript
}

func ValidateBtcTransaction(tx string) error {
	if len(tx) == 0 {
		return fmt.Errorf("transaction cannot be empty")
	}
	txBytes, err := hex.DecodeString(tx)
	if err != nil {
		return fmt.Errorf("invalid transaction data: not a valid hex tx %s", tx)
	}
	if len(txBytes) < 50 || len(txBytes) > 100000 {
		return fmt.Errorf("invalid transaction size: must be between 50 and 100000 bytes")
	}
	return nil
}

func IsValidBtcTxHash(txHash string) bool {
	_, err := hex.DecodeString(txHash)
	if err != nil {
		return false
	}
	return len(txHash) == 64
}
