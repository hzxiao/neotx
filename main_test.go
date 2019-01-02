package main

import (
	"testing"
)

func TestHexlify(t *testing.T) {
	println(hexlify("w"))
	println(unhexlify(hexlify("w")))

	println(computeTxid(`80000001838d0929ea7cbaa69956078f9471b996af4c4028662243ac6d67c481b22c899b000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c600065cd1d0000000005a0a304bac8edf51064a9165670cb39fb87439e`))
}
