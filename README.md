# Go PGP Utility

This project provides a set of Go utilities for performing PGP encryption and decryption, as well as some file manipulation tasks like writing to TXT, CSV, and creating ZIP archives.

## Features

* **PGP Encryption**: Encrypt files using a public PGP key.
* **PGP Decryption**:
  * Decrypt files using a private PGP key and a passphrase.
  * Decrypt files using a private PGP key from a keyring (without a passphrase if the key is not encrypted).
* **File Utilities**:
  * Write content to text files.
  * Write data to CSV files with a custom delimiter.
  * Create ZIP archives from a list of files.

## Prerequisites

* Go (version 1.18 or later recommended)
* PGP keys (public and private) for testing encryption and decryption.

## Dependencies

* [github.com/ProtonMail/gopenpgp/v3](https://github.com/ProtonMail/gopenpgp)

## Setup

1. **Clone the repository (if applicable) or ensure you have the project files.**
2. **Place your PGP keys**:

    * For encryption: `public.pgp` (or update the path in `main.go`)
    * For decryption with passphrase: `private.pgp` (or update the path in `main.go`)
    * For decryption with keyring: `secret.asc` (or update the path in `main.go`)

3. **Prepare input files**:

    * Create a file `decrypted/test.txt` with content to be encrypted.
    * Place an encrypted file at `encrypted/test.txt.zip.pgp` for decryption with passphrase.
    * Place an encrypted file at `encrypted/test.csv.gpg` for decryption with keyring.
    * Ensure the output directories (e.g., `decrypted/`, `encrypted/`) exist or create them.
