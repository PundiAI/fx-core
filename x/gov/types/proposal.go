package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type FXMetadata struct {
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Metadata string `json:"metadata"`
}

func (md FXMetadata) String() string {
	bz, _ := json.Marshal(md)
	return base64.StdEncoding.EncodeToString(bz)
}

func NewFXMetadata(title, summary, metadata string) FXMetadata {
	return FXMetadata{
		Title:    title,
		Summary:  summary,
		Metadata: metadata,
	}
}

func ParseFXMetadata(fxMDStr string) (fxMD FXMetadata, err error) {
	if len(strings.TrimSpace(fxMDStr)) == 0 {
		return FXMetadata{}, fmt.Errorf("fx metadata cannot be empty")
	}
	bz, err := base64.StdEncoding.DecodeString(fxMDStr)
	if err != nil {
		return FXMetadata{}, err
	}
	err = json.Unmarshal(bz, &fxMD)
	return fxMD, err
}
