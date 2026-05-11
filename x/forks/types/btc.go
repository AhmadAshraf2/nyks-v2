package types

import (
	"bytes"
	"encoding/hex"
	fmt "fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

const (
	BtcPublicKeyLen = 66
)

// BtcPublicKey represents a BTC public key
type BtcPublicKey struct {
	BtcPublicKey string
}

func (ea BtcPublicKey) GetBtcPublicKey() string {
	return ea.BtcPublicKey
}

func (ea BtcPublicKey) SetBtcPublicKey(PublicKey string) error {
	if err := ValidateBtcPublicKey(PublicKey); err != nil {
		return err
	}
	ea.BtcPublicKey = PublicKey
	return nil
}

func NewBtcPublicKey(PublicKey string) (*BtcPublicKey, error) {
	if err := ValidateBtcPublicKey(PublicKey); err != nil {
		return nil, fmt.Errorf("invalid input PublicKey: %w", err)
	}
	addr := BtcPublicKey{PublicKey}
	return &addr, nil
}

func ValidateBtcPublicKey(PublicKey string) error {
	if PublicKey == "" {
		return fmt.Errorf("empty")
	}
	return nil
}

func (ea BtcPublicKey) ValidateBasic() error {
	return ValidateBtcPublicKey(ea.BtcPublicKey)
}

func BtcAddrLessThan(e BtcPublicKey, o BtcPublicKey) bool {
	return bytes.Compare([]byte(e.GetBtcPublicKey()), []byte(o.GetBtcPublicKey())) == -1
}

func NewBtcPublicKeyFromBytes(publicKeyBytes []byte) (*BtcPublicKey, error) {
	if err := ValidateBtcPublicKey(hex.EncodeToString(publicKeyBytes)); err != nil {
		return nil, fmt.Errorf("invalid input publicKeyBytes: %w", err)
	}
	pk := BtcPublicKey{hex.EncodeToString(publicKeyBytes)}
	return &pk, nil
}

// CreateTxFromHex creates a btc transaction object from a hex string
func CreateTxFromHex(txHex string) (*wire.MsgTx, error) {
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %v", err)
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	err = tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}
	return tx, nil
}

// CreateTxHashFromHex creates a btc transaction hash from a hex string
func CreateTxHashFromHex(txHex string) (*chainhash.Hash, error) {
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return &chainhash.Hash{}, fmt.Errorf("failed to decode hex string: %v", err)
	}
	txHash := chainhash.DoubleHashH(txBytes)
	return &txHash, nil
}
