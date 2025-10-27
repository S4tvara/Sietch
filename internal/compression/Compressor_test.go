package compression

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/substantialcattle5/sietch/internal/constants"
)

func TestCompressAndDecompressData(t *testing.T) {
	input := []byte("Some string test data")
	tests := []struct {
		name            string
		algorithm       string
		isErrorExpected bool
		err             error
	}{
		{
			name:            "With Gzip Algorithm",
			algorithm:       constants.CompressionTypeGzip,
			isErrorExpected: false,
		},
		{
			name:            "With Zstd Algorithm",
			algorithm:       constants.CompressionTypeZstd,
			isErrorExpected: false,
		},
		{
			name:            "With Lz4 Algorithm",
			algorithm:       constants.CompressionTypeLz4,
			isErrorExpected: false,
		},
		{
			name:            "With No Compression Algorithm",
			algorithm:       constants.CompressionTypeNone,
			isErrorExpected: false,
		},
		{
			name:            "Unsupported Compression Algorithm",
			algorithm:       "something else",
			isErrorExpected: true,
			err:             fmt.Errorf("unsupported compression algorithm: something else"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := CompressData(input, tt.algorithm)
			if tt.isErrorExpected {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if compressed != nil {
					t.Errorf("Expected nil compressed data on error")
				}
				if err != nil && err.Error() != tt.err.Error() {
					t.Errorf("Wrong error recieved. Expected: %v. Recieved: %v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error during compression: %v", err)
				return
			}
			if compressed == nil {
				t.Errorf("Compressed data is nil")
				return
			}

			decompressed, err := DecompressData(compressed, tt.algorithm)
			if err != nil {
				t.Errorf("Unexpected error during decompression: %v", err)
				return
			}
			if string(decompressed) != string(input) {
				t.Errorf("Decompressed data does not match original. Got: %s", string(decompressed))
			}
		})
	}
}

func TestDecompressBomb(t *testing.T) {
	raw := bytes.Repeat([]byte("a"), constants.MaxDecompressionSize+1)

	tests := []struct {
		name      string
		algorithm string
	}{
		// {
		// 	name:      "With Gzip Algorithm",
		// 	algorithm: constants.CompressionTypeGzip,
		// },
		{
			name:      "With LZ4 Algorithm",
			algorithm: constants.CompressionTypeLz4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := CompressData(raw, tt.algorithm)
			if err != nil {
				t.Errorf("Could not setup decompression bomb. Error: %v.", err)
				return
			}
			if compressed == nil {
				t.Errorf("Could not setup decompression bomb. Setup resulted in nil")
				return
			}

			decompressed, err := DecompressData(compressed, tt.algorithm)
			if decompressed != nil {
				t.Errorf("Unexpected response. Expected nil.")
				return
			}
			if err == nil || err.Error() != fmt.Sprintf("decompressed data exceeds maximum size limit (%d bytes) - potential decompression bomb", constants.MaxDecompressionSize) {
				t.Errorf("Incorrect error recieved. Error: %v", err)
			}
		})
	}
}
