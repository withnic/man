package main

import "testing"

func TestTranseformURL(t *testing.T) {
	url, err := transeformURL("http://localhost.com")
	if err != nil {
		t.Errorf("has error :%s", err.Error())
	}

	expected := "http://localhost.com"
	if url != expected {
		t.Errorf("got %s expected %s", url, expected)
	}
}
