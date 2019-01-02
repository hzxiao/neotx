package req

import (
	"testing"
)

func TestGetUtxoOutputValue(t *testing.T) {
	err := SetNetwork(TestNet, "")
	if err != nil {
		t.Error(err)
	}

	value, err := GetUtxoOutputValue("1aaa92ad08c7ee2b8f67d76cde4893096ccafcaa1703507cec3d5ed087368b45", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(value)
}
