package mt4SyncEt6

import (
	"fmt"
	"mt4SyncEt6/security"
	"testing"
)

func TestGetGroup(t *testing.T) {
	a:=security.GetGroup()
	fmt.Println(a)
}
