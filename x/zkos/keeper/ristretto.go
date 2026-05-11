package keeper

import (
	"encoding/binary"
	"encoding/hex"

	r255 "github.com/gtank/ristretto255"
)

// PublicKey struct (g^r, g^r^sk)
type PublicKey struct {
	Gr   *r255.Element
	Grsk *r255.Element
}

// Commitment struct
type Commitment struct {
	C *r255.Element
	D *r255.Element
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func DecodePoints(bytes []byte) (*r255.Element, *r255.Element) {
	xBytes := make([]byte, 32)
	copy(xBytes[:], bytes[0:32])

	yBytes := make([]byte, 32)
	copy(yBytes[:], bytes[32:64])

	x := r255.NewElement()
	x.Decode(xBytes[:])

	y := r255.NewElement()
	y.Decode(yBytes[:])

	return x, y
}

func DecodeAccount(accStr string) (*PublicKey, *Commitment, error) {
	accountBytes, err := hex.DecodeString(accStr)
	if err != nil {
		return nil, nil, err
	}

	extractedPkBytes := accountBytes[1:65]
	g, h := DecodePoints(extractedPkBytes)

	extractedCommBytes := accountBytes[69:]
	c, d := DecodePoints(extractedCommBytes)

	return &PublicKey{Gr: g, Grsk: h}, &Commitment{C: c, D: d}, nil
}

func ScalarFromBytes(b [32]byte) (*r255.Scalar, error) {
	s := r255.NewScalar()
	err := s.Decode(b[:])
	if err != nil {
		return nil, err
	}
	return s, nil
}

func uintToScalar(bl uint64) (*r255.Scalar, error) {
	blslice := make([]byte, 32)
	binary.LittleEndian.PutUint64(blslice, bl)

	blArray := [32]byte{}
	copy(blArray[:], blslice[:32])

	intScalar, err := ScalarFromBytes(blArray)
	if err != nil {
		return nil, err
	}
	return intScalar, nil
}

func GenerateCommitment(p *PublicKey, rscalar *r255.Scalar, intScalar *r255.Scalar) *Commitment {
	e := r255.NewElement()
	c := e.ScalarMult(rscalar, p.Gr)

	b := r255.NewElement()
	gv := b.ScalarBaseMult(intScalar)

	k := r255.NewElement()
	kh := k.ScalarMult(rscalar, p.Grsk)

	n := r255.NewElement()
	d := n.Add(gv, kh)

	return &Commitment{C: c, D: d}
}

func CompareCommitment(u *Commitment, v *Commitment) bool {
	return u.C.Equal(v.C) == 1 && u.D.Equal(v.D) == 1
}

func (k msgServer) RevealCommitment(accStr string, scalarStr string, value uint64) (bool, error) {
	pk, comm, err := DecodeAccount(accStr)
	if err != nil {
		return false, err
	}

	scalarBytes, err := hex.DecodeString(scalarStr)
	if err != nil {
		return false, err
	}
	scalarSlice := [32]byte{}
	copy(scalarSlice[:], scalarBytes[:32])
	scalarComm, err := ScalarFromBytes(scalarSlice)
	if err != nil {
		return false, err
	}

	uintScalar, err := uintToScalar(value)
	if err != nil {
		return false, err
	}

	newCommitment := GenerateCommitment(pk, scalarComm, uintScalar)

	checkComm := CompareCommitment(comm, newCommitment)
	if checkComm {
		return true, nil
	}
	return false, nil
}
