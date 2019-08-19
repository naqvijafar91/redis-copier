package main

import "testing"

func TestDestinationSafetyHook(t *testing.T) {
	if !shouldAllowDestinationAddress("localhost:6379") {
		t.Error("Should have passed")
	}
	if shouldAllowDestinationAddress("xyz.bigactualprod.omega-cloud.com:6379") {
		t.Error("Should have failed")
	}
}
