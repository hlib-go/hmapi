package hmapi

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token := GenToken("1", "1", 100)
	fmt.Println(token)
	uid, err := VerToken(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(uid)
}
