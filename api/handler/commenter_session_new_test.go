package handler

import (
	"simple-commenting/test"
	"testing"
)

func TestCommenterTokenNewBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if _, err := commenterTokenNew(); err != nil {
		t.Errorf("unexpected error creating new commenterToken: %v", err)
		return
	}
}
