package file

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alexmullins/zip"
)

func WriteTxt(outputDir, fileName string, content []byte) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	filePath := filepath.Join(outputDir, fileName)

	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to text file %s: %w", filePath, err)
	}

	return nil
}

func WriteCsv(outputDir, fileName string, data [][]string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	filePath := filepath.Join(outputDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '|'
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to csv file %s: %w", filePath, err)
		}
	}

	return nil
}

func WriteZip(outputDir, zipFileName string, filesToZip []string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	zipFilePath := filepath.Join(outputDir, zipFileName)
	newZipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file %s: %w", zipFilePath, err)
	}

	zipDirName := zipFileName
	if ext := filepath.Ext(zipFileName); ext == ".zip" {
		zipDirName = zipFileName[:len(zipFileName)-len(ext)]
	}

	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	for _, fileName := range filesToZip {
		zipFilePath := filepath.Join(outputDir, fileName)
		fileToZip, err := os.Open(zipFilePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s for zipping: %w", zipFilePath, err)
		}

		defer fileToZip.Close()
		info, err := fileToZip.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", zipFilePath, err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header for %s: %w", zipFilePath, err)
		}

		header.Name = filepath.Join(zipDirName, fileName)
		writerInZip, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create entry in zip for %s: %w", zipFilePath, err)
		}

		if _, err = io.Copy(writerInZip, fileToZip); err != nil {
			return fmt.Errorf("failed to copy file %s to zip: %w", zipFilePath, err)
		}
	}

	return nil
}

func WriteEncryptedZip(outputDir string, zipFileName string, filesToZip []string, password string) error {
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	zipFilePath := filepath.Join(outputDir, zipFileName)
	newZipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file %s: %w", zipFilePath, err)
	}

	zipDirName := zipFileName
	if ext := filepath.Ext(zipFileName); ext == ".zip" {
		zipDirName = zipFileName[:len(zipFileName)-len(ext)]
	}

	defer newZipFile.Close()
	defer zipWriter.Close()
	for _, fileName := range filesToZip {
		filePath := filepath.Join(outputDir, fileName)
		fileToZip, _ := os.Open(filePath)
		w, err := zipWriter.Encrypt(filepath.Join(zipDirName, fileName), password)
		if err != nil {
			return fmt.Errorf("failed to create encrypted entry for %s: %w", filePath, err)
		}

		if _, err = io.Copy(w, fileToZip); err != nil {
			return fmt.Errorf("failed to copy file %s to zip: %w", zipFilePath, err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close zip writer: %w", err)
	}

	err = os.WriteFile(zipFilePath, zipBuffer.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Failed to write ZIP file %s: %v\n", zipFilePath, err)
		return err
	}

	fmt.Printf("Successfully created encrypted ZIP file at %s\n", zipFilePath)
	return nil
}
