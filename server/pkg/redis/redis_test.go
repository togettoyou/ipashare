package redis

import (
	"context"
	"testing"
)

func TestRedis(t *testing.T) {
	err := Setup(0, "127.0.0.1:6379", "")
	if err != nil {
		t.Fatal(err)
	}
	err = Set(context.Background(), "key", "value", 30)
	if err != nil {
		t.Fatal(err)
	}
	get, err := Get(context.Background(), "key")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(get)
}
