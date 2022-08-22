package utils

import "testing"

func TestRand(t *testing.T) {
	t.Logf("Get rand string: %s.", RandString(25))
	t.Logf("Get rand byte: %v.", RandByte(10))

	nums, err := RandCountForDiff(1, 100, 10)
	if err != nil {
		t.Logf("Get rand count for diff err: %v.", err)
	} else {
		t.Logf("Get rand count for diff: %v.", nums)
	}
	t.Logf("Get rand by area: %d.", RandByArea(1, 10))
}
