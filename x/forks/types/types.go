package types

import (
	"crypto/md5"
	"encoding/binary"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HashString hashes a string using MD5 and returns 16 bytes
func HashString(input string) []byte {
	md5 := md5.New()
	md5.Write([]byte(input))
	return md5.Sum(nil)
}

// AppendBytes concatenates multiple byte slices
func AppendBytes(args ...[]byte) []byte {
	length := 0
	for _, v := range args {
		length += len(v)
	}

	res := make([]byte, length)
	length = 0
	for _, v := range args {
		copy(res[length:length+len(v)], v)
		length += len(v)
	}

	return res
}

// UInt64Bytes uses the SDK byte marshaling to encode a uint64
func UInt64Bytes(n uint64) []byte {
	return sdk.Uint64ToBigEndian(n)
}

// UInt64FromBytes creates uint64 from binary big endian representation
func UInt64FromBytes(s []byte) uint64 {
	return binary.BigEndian.Uint64(s)
}

// AttestationVotesPowerThreshold is the threshold of power needed for an attestation to be observed (67%)
var AttestationVotesPowerThreshold = math.NewInt(67)

// AttestationVoteCountThreshold is the threshold of vote count needed for proposals (67%)
var AttestationVoteCountThreshold = math.NewInt(67)
