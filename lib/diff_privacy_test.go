package libdrynx

import (
	"github.com/stretchr/testify/assert"
	"go.dedis.ch/onet/v3/log"
	"testing"
)

func TestGenerateNoiseValues(t *testing.T) {
	aux := GenerateNoiseValues(0, 0, 1, 0.005, 0)
	assert.Empty(t, aux)

	aux = GenerateNoiseValues(500, 0, 1, 0.005, 0)
	log.Lvl1(aux)
	log.Lvl1(len(aux))
	assert.Equal(t, len(aux), 500)

	temp := make([]float64, 0)
	for i := 0; i < 100; i++ {
		temp = append(temp, 0)
	}

	assert.Equal(t, temp, aux[:100])

	temp = make([]float64, 0)
	for i := 0; i < 19; i++ {
		temp = append(temp, 1)
		temp = append(temp, -1)
	}

	assert.Equal(t, temp, aux[100:138])

	aux = GenerateNoiseValuesScale(500, 0, 1, 0.005, 100, 60)
	log.Lvl1(aux)
}
