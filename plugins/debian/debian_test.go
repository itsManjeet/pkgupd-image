package debian

import (
	"fmt"
	"testing"
)

func TestDepends(t *testing.T) {
	deb := Debian{}
	pkg, err := deb.Prepare("main")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(pkg)
}
