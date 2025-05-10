package main

import (
	"fmt"
	"go-pgp/pgp"
)

func main() {
	err := pgp.Encrypt("public.pgp", "decrypted/006_Merchant_Report_20250429_1.txt", "006_Merchant_Report_20250429_1.txt.pgp")
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	err = pgp.DecryptWithPassphrase("private.pgp", "tungngern1234", "encrypted/006_Merchant_Report_20250429_1.txt.zip.pgp", "006_Merchant_Report_20250429_1.txt.zip")
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	err = pgp.DecryptWithKeyRing("secret.asc", "encrypted/data_casa_dga_to_fi_006_25680424_1.csv.gpg", "data_casa_dga_to_fi_006_25680424_1.csv")
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
}
