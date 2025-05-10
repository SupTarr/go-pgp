package main

import (
	"fmt"
	"go-pgp/pgp"
)

func main() {
	err := pgp.Encrypt("public.pgp", "decrypted/test.txt", "test.txt.pgp")
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	err = pgp.DecryptWithPassphrase("private.pgp", "test", "encrypted/test.txt.zip.pgp", "test.txt.zip")
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	err = pgp.DecryptWithKeyRing("secret.asc", "encrypted/test.csv.gpg", "test.csv")
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
}
