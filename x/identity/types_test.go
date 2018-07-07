package identity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	metadata := Metadata(`{"demo":"1"}`)
	b, err := json.Marshal(metadata)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(b), `{"demo":"1"}`)
	var newMeta Metadata
	err = json.Unmarshal(b, &newMeta)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(newMeta), string(metadata))
}
