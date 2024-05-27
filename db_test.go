package flashdb

import (
	"os"
	"testing"
)

func assert(a bool) {
	if !a {
		panic("failed assert")
	}
}

func Test_FastDb(t *testing.T) {
	os.Remove("test1")

	db := New()
	key1 := "key1"
	key2 := "key22222"
	key3 := "key3333"

	db.Add(key1, nil)
	db.Add(key2, []byte("123"))
	db.Add(key3, []byte("12"))

	assert(db.Exist(key1))
	assert(db.Exist(key2))
	assert(db.Exist(key3))

	must := map[string]string{
		key1: "",
		key2: "123",
		key3: "12",
	}

	db.Iterate(func(k string, v []byte) bool {
		assert(must[k] == string(v))
		return true
	})

	assert(db.Count() == 3)

	assert(db.Save() != nil)

	err := db.SaveAs("test1")
	if err != nil {
		t.Fatalf("err: %s\n", err.Error())
	}
	assert(err == nil)

	// -----------------------------------
	db2 := New()
	err = db2.Load("test1")
	if err != nil {
		t.Logf("test1 error " + err.Error())
	}
	assert(err == nil)

	assert(db2.Count() == 3)

	assert(db2.Exist(key1))
	assert(db2.Exist(key2))
	assert(db2.Exist(key3))
}
