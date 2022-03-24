package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func accountToHex(addr string) (string, error) {
	if strings.HasPrefix(addr, sdk.GetConfig().GetBech32AccountAddrPrefix()) {
		// Check to see if address is Cosmos bech32 formatted
		toAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return "", errors.Wrap(err, "must provide a valid Bech32 address")
		}
		ethAddr := common.BytesToAddress(toAddr.Bytes())
		return ethAddr.Hex(), nil
	}

	if !strings.HasPrefix(addr, "0x") {
		addr = "0x" + addr
	}

	valid := common.IsHexAddress(addr)
	if !valid {
		return "", fmt.Errorf("%s is not a valid Ethereum or Cosmos address", addr)
	}

	ethAddr := common.HexToAddress(addr)

	return ethAddr.Hex(), nil
}

func formatKeyToHash(key string) string {
	if !strings.HasPrefix(key, "0x") {
		key = "0x" + key
	}

	ethkey := common.HexToHash(key)

	return ethkey.Hex()
}

func ParseMetadata(cdc codec.Codec, metadataFile string) (banktypes.Metadata, error) {
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

func ReadMetadataFromPath(cdc codec.Codec, path string) ([]banktypes.Metadata, error) {
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
