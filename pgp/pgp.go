package pgp

import (
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func Encrypt(publicKeyPath, inputFilePath, outputFilePath string) error {
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKeyObj, err := crypto.NewKeyFromArmored(string(publicKey))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	pgp := crypto.PGP()
	publicKeyRing, err := crypto.NewKeyRing(publicKeyObj)
	if err != nil {
		return fmt.Errorf("failed to create key ring from public key: %w", err)
	}

	encHandle, err := pgp.Encryption().Recipients(publicKeyRing).New()
	if err != nil {
		return fmt.Errorf("failed to handle encryption: %w", err)
	}

	message, err := os.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to read message file: %w", err)
	}

	encryptedPGPMessage, err := encHandle.Encrypt(message)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %w", err)
	}

	armoredOutput, err := encryptedPGPMessage.ArmorBytes()
	if err != nil {
		return fmt.Errorf("failed to get armored PGP message: %w", err)
	}

	err = os.WriteFile(outputFilePath, armoredOutput, 0644)
	if err != nil {
		return fmt.Errorf("failed to write encrypted message to file: %w", err)
	}

	return nil
}

func DecryptWithPassphrase(privateKeyPath, passphrase, encryptedFilePath, outputFilePath string) error {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("failed to read private key file '%s': %v\n", privateKeyPath, err)
		return err
	}

	privateKeyObj, err := crypto.NewPrivateKeyFromArmored(string(privateKey), []byte(passphrase))
	if err != nil {
		fmt.Printf("failed to parse private key from '%s': %v\n", privateKeyPath, err)
		return err
	}

	defer privateKeyObj.ClearPrivateParams()
	pgp := crypto.PGP()
	decryptionHandle, err := pgp.Decryption().DecryptionKey(privateKeyObj).New()
	if err != nil {
		fmt.Printf("failed to decryption handler: %v\n", err)
		return err
	}

	defer decryptionHandle.ClearPrivateParams()
	encryptedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		fmt.Printf("failed to open encrypted file '%s': %v\n", encryptedFilePath, err)
		return err
	}

	defer encryptedFile.Close()
	ptReader, err := decryptionHandle.DecryptingReader(encryptedFile, crypto.Auto)
	if err != nil {
		fmt.Printf("failed to decrypt stream from '%s' to '%s': %v\n", encryptedFilePath, outputFilePath, err)
		return err
	}

	decrypted, err := ptReader.ReadAllAndVerifySignature()
	if err != nil {
		fmt.Println("failed to read and verify signature: ", err)
		return err
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("failed to create output file '%s': %v\n", outputFilePath, err)
		return err
	}

	defer outputFile.Close()

	_, err = outputFile.Write(decrypted.Bytes())
	if err != nil {
		fmt.Printf("failed to write decrypted content to '%s': %v\n", outputFilePath, err)
		return err
	}

	fmt.Printf("Successfully decrypted '%s' to '%s'\n", encryptedFilePath, outputFilePath)
	return nil
}

func DecryptWithKeyRing(privateKeyPath, encryptedFilePath, outputFilePath string) error {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("failed to read private key file '%s': %v\n", privateKeyPath, err)
		return err
	}

	privateKeyObj, err := crypto.NewKeyFromArmored(string(privateKey))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	pgp := crypto.PGP()
	privateKeyRing, err := crypto.NewKeyRing(privateKeyObj)
	if err != nil {
		return fmt.Errorf("failed to create key ring from public key: %w", err)
	}

	decryptionHandle, err := pgp.Decryption().DecryptionKeys(privateKeyRing).New()
	if err != nil {
		fmt.Printf("failed to decryption handler: %v\n", err)
		return err
	}

	defer decryptionHandle.ClearPrivateParams()
	encryptedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		fmt.Printf("failed to open encrypted file '%s': %v\n", encryptedFilePath, err)
		return err
	}

	defer encryptedFile.Close()
	ptReader, err := decryptionHandle.DecryptingReader(encryptedFile, crypto.Auto)
	if err != nil {
		fmt.Printf("failed to decrypt stream from '%s' to '%s': %v\n", encryptedFilePath, outputFilePath, err)
		return err
	}

	decrypted, err := ptReader.ReadAllAndVerifySignature()
	if err != nil {
		fmt.Println("failed to read and verify signature: ", err)
		return err
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("failed to create output file '%s': %v\n", outputFilePath, err)
		return err
	}

	_, err = outputFile.WriteString(decrypted.String())
	if err != nil {
		fmt.Printf("failed to write decrypted content to '%s': %v\n", outputFilePath, err)
		return err
	}

	fmt.Printf("Successfully decrypted '%s' to '%s'\n", encryptedFilePath, outputFilePath)
	return nil
}
