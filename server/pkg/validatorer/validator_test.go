package validatorer

import "testing"

func TestValidator(t *testing.T) {
	Setup()
	data1 := "zoujh99@qq.com"
	data2 := "qq.com"
	data3 := "1234"
	if msg, ok := E(V().Var(data1, "required,email")); !ok {
		t.Log(msg)
	}
	if msg, ok := E(V().Var(data2, "required,email")); !ok {
		t.Log(msg)
	}
	if msg, ok := E(V().Var(data3, "required,checkUsername")); !ok {
		t.Log(msg)
	}
}
