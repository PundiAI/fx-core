package cli

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// ParseRegisterCoinProposal reads and parses a ParseRegisterCoinProposal from a file.
func ParseMetadata(cdc codec.JSONMarshaler, metadataFile string) (banktypes.Metadata, error) {
	metadata := banktypes.Metadata{}

	contents, err := ioutil.ReadFile(filepath.Clean(metadataFile))
	if err != nil {
		return metadata, err
	}

	if err = cdc.UnmarshalJSON(contents, &metadata); err != nil {
		return metadata, err
	}

	return metadata, nil
}

func ReadMetadataFromPath(cdc codec.JSONMarshaler, path string) ([]banktypes.Metadata, error) {
	metadatas := make([]banktypes.Metadata, 0, 10)
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path %s error %v", path, err)
	}
	if stat.IsDir() {
		if err := filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			metadata, err := ParseMetadata(cdc, p)
			if err != nil {
				return fmt.Errorf("parse metadata file %s error %v", p, err)
			}
			metadatas = append(metadatas, metadata)
			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		metadata, err := ParseMetadata(cdc, path)
		if err != nil {
			return nil, fmt.Errorf("parse metadata file %s error %v", path, err)
		}
		metadatas = append(metadatas, metadata)
	}

	return metadatas, nil
}
