package dnsiterativeecs

import (
	"testing"
)

func TestLookupAFail(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := Lookup("1.1.1.1", "baidu.com.")
		if err != nil {
			t.Error(err)
		}
		err = Lookup("1.1.1.1", "www.google.com.")
		if err != nil {
			t.Error(err)
		}
		err = Lookup("1.1.1.1", "www.bilibili.com.")
		if err != nil {
			t.Error(err)
		}
	}
}
