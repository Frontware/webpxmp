package xmp

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	fmt.Println("\u524b\u1533\u5460\u554a")
	fmt.Println([]byte("\u554a"))
	fmt.Printf("%+q\n", "啊")
	var b []byte
	b = strconv.AppendQuoteToASCII(b, "\u524b\u1533\u5460\u554a")
	fmt.Println(b)
	fmt.Println(string(b))
	tes, err := strconv.Unquote(string(b))
	if err != nil {
		panic(err)
	}
	fmt.Println(tes)

	p, err := Read("https://ipfs.weladee.com/ipfs/QmbkyX4j7qUN7fTyMc1RBH6KLD8LApxuur4PJqEfRFdnVG")
	if err != nil {
		t.Logf("error: %v\n", err)
	}
	t.Log(p)

	err = Write("test.webp", "อง", time.Now().UnixNano(), 0.12321331, 1.12312321, "test location")
	err = Write("test2.webp", "ของ", time.Now().UnixNano(), 0.12321331, 1.12312321, "test location")
	err = Write("extended.webp", "Testing unicode éïööçù", time.Now().UnixNano(), 0.12321331, 1.12312321, "test location")
	if err != nil {
		t.Logf("error: %v\n", err)
	}

	p, err = Read("test.webp")
	if err != nil {
		t.Logf("error: %v\n", err)
	}
	t.Log(p)
	p, err = Read("test2.webp")
	if err != nil {
		t.Logf("error: %v\n", err)
	}
	t.Log(p)
	p, err = Read("extended.webp")
	if err != nil {
		t.Logf("error: %v\n", err)
	}
	t.Log(p)
}
