package encoding_test

import (
	"github.com/lca1/drynx/lib"
	"github.com/lca1/drynx/lib/encoding"
	"github.com/lca1/unlynx/lib"
	"github.com/stretchr/testify/assert"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/onet/v3/log"
	"math"
	"testing"
)

func TestEncodeDecodeCosim(t *testing.T) {
	secKey, pubKey := libunlynx.GenKey()

	limit := int64(10000)
	log.Lvl1("Preparing decryption up to:", limit)

	// Decrpytion hashtable creation
	libdrynx.CreateDecryptionTable(limit, pubKey, secKey)

	//data
	rijs := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12}
	riks := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12}

	// excpected results
	rijs_sum := int64(0)
	riks_sum := int64(0)
	rijs_2_sum := int64(0)
	riks_2_sum := int64(0)
	rijs_x_rijks_sum := int64(0)

	for i, el := range rijs {
		el2 := riks[i]
		rijs_sum = rijs_sum + el
		riks_sum = riks_sum + el2
		rijs_2_sum = rijs_2_sum + el*el
		riks_2_sum = riks_2_sum + el2*el2
		rijs_x_rijks_sum = rijs_x_rijks_sum + el*el2

	}
	resultClear := []int64{rijs_sum, riks_sum, rijs_2_sum, riks_2_sum, rijs_x_rijks_sum}
	log.Lvl1("Preliminary Results ", resultClear)

	//expected results
	expect := float64(resultClear[4]) / (math.Sqrt(float64(resultClear[2])) * math.Sqrt(float64(resultClear[3])))
	log.Lvl1("Expected Preliminary Results ", resultClear)

	resultEncrypted, _ := encoding.EncodeCosim(rijs, riks, pubKey)
	result := encoding.DecodeCosim(resultEncrypted, secKey)
	log.Lvl1("Final Results ", result)
	assert.Equal(t, expect, result)
}

func TestEncodeDecodeCosimWithProofs(t *testing.T) {
	// key
	secKey, pubKey := libunlynx.GenKey()

	limit := int64(10000)
	log.Lvl1("Preparing decryption up to:", limit)

	// Decrpytion hashtable creation
	libdrynx.CreateDecryptionTable(limit, pubKey, secKey)

	//data
	rijs := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1}
	riks := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1}

	// excpected results
	rijsSum := int64(0)
	riksSum := int64(0)
	rijs2Sum := int64(0)
	riks2Sum := int64(0)
	rijsXRijksSum := int64(0)

	for i, el := range rijs {
		el2 := riks[i]
		rijsSum = rijsSum + el
		riksSum = riksSum + el2
		rijs2Sum = rijs2Sum + el*el
		riks2Sum = riks2Sum + el2*el2
		rijsXRijksSum = rijsXRijksSum + el*el2

	}
	resultClear := []int64{rijsSum, riksSum, rijs2Sum, riks2Sum, rijsXRijksSum}

	//expected results
	expect := float64(resultClear[4]) / (math.Sqrt(float64(resultClear[2])) * math.Sqrt(float64(resultClear[3])))

	//proofs signatures
	//signatures needed to check the proof
	u := int64(2)
	l := int64(10)
	ps := make([][]libdrynx.PublishSignature, 2)

	ranges := make([]*[]int64, len(resultClear))
	ps[0] = make([]libdrynx.PublishSignature, len(resultClear))
	ps[1] = make([]libdrynx.PublishSignature, len(resultClear))
	ys := make([][]kyber.Point, 2)
	ys[0] = make([]kyber.Point, len(resultClear))
	ys[1] = make([]kyber.Point, len(resultClear))
	for i := range ps[0] {
		ps[0][i] = libdrynx.PublishSignatureBytesToPublishSignatures(libdrynx.InitRangeProofSignature(u))
		ps[1][i] = libdrynx.PublishSignatureBytesToPublishSignatures(libdrynx.InitRangeProofSignature(u))
		ys[0][i] = ps[0][i].Public
		ys[1][i] = ps[1][i].Public
		ranges[i] = &[]int64{u, l}
	}

	yss := make([][]kyber.Point, len(resultClear))
	for i := range yss {
		yss[i] = make([]kyber.Point, 2)
		for j := range ys {
			yss[i][j] = ys[j][i]
		}
	}

	resultEncrypted, _, prf := encoding.EncodeCosimWithProofs(rijs, riks, pubKey, ps, ranges)
	for i, v := range prf {
		yss := make([]kyber.Point, 2)
		for j := range ys {
			yss[j] = ys[j][i]
		}
		assert.True(t, libdrynx.RangeProofVerification(libdrynx.CreatePredicateRangeProofForAllServ(v), u, l, yss, pubKey))
	}

	result := encoding.DecodeCosim(resultEncrypted, secKey)
	assert.Equal(t, expect, result)
}
