package main

import "testing"

func TestTranseformURL(t *testing.T) {
	url, err := transeformURL("http://localhost.com")
	if err != nil {
		t.Fatalf("has error :%s", err.Error())
	}

	want := "http://localhost.com"
	if url != want {
		t.Fatalf("got %s want %s", url, want)
	}
}
