package systeminfo

import (
	"fmt"
	"testing"
)

func TestNewSystem(t *testing.T) {
	s := New()

	fmt.Println(s.GetHostUniqueID())
}
