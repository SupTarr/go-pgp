package file

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, filePath := range filesToZip {
		fileToZip, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s for zipping: %w", filePath, err)
		}
		defer fileToZip.Close()

		info, err := fileToZip.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", filePath, err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header for %s: %w", filePath, err)
		}

		header.Name = filepath.Base(filePath)
		writerInZip, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create entry in zip for %s: %w", filePath, err)
		}

		if _, err = io.Copy(writerInZip, fileToZip); err != nil {
			return fmt.Errorf("failed to copy file %s to zip: %w", filePath, err)
		}
	}

	return nil
}
