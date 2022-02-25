package dao

import (
	"testing"
)

func TestDB(t *testing.T) {
	store, err := NewSqlite()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(store.AppleDeveloper.GetUsable())
}
