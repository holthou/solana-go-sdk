package token_metadata

import (
	"encoding/binary"
	"fmt"
	"github.com/blocto/solana-go-sdk/program/token"
	"reflect"
	"strings"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/near/borsh-go"
)

const EDITION_MARKER_BIT_SIZE uint64 = 248

type Key borsh.Enum

const (
	KeyUninitialized Key = iota
	KeyEditionV1
	KeyMasterEditionV1
	KeyReservationListV1
	KeyMetadataV1
	KeyReservationListV2
	KeyMasterEditionV2
	KeyEditionMarker
	KeyUseAuthorityRecord
	KeyCollectionAuthorityRecord
)

type Creator struct {
	Address  common.PublicKey
	Verified bool
	Share    uint8
}

type Data struct {
	Name                 string
	Symbol               string
	Uri                  string
	SellerFeeBasisPoints uint16
	Creators             *[]Creator
}

type DataV2 struct {
	Name                 string
	Symbol               string
	Uri                  string
	SellerFeeBasisPoints uint16
	Creators             *[]Creator
	Collection           *Collection
	Uses                 *Uses
}

type metadataPreV11 struct {
	Key                 Key
	UpdateAuthority     common.PublicKey
	Mint                common.PublicKey
	Data                Data
	PrimarySaleHappened bool
	IsMutable           bool
	EditionNonce        *uint8
}

type Metadata struct {
	Key                 Key
	UpdateAuthority     common.PublicKey
	Mint                common.PublicKey
	Data                Data
	PrimarySaleHappened bool
	IsMutable           bool
	EditionNonce        *uint8
	TokenStandard       *TokenStandard
	Collection          *Collection
	Uses                *Uses
	CollectionDetails   *CollectionDetails
	ProgrammableConfig  *ProgrammableConfig
}

type TokenStandard borsh.Enum

const (
	NonFungible TokenStandard = iota
	FungibleAsset
	Fungible
	NonFungibleEdition
	ProgrammableNonFungible
)

type Collection struct {
	Verified bool
	Key      common.PublicKey
}

type Uses struct {
	UseMethod UseMethod
	Remaining uint64
	Total     uint64
}

type UseMethod borsh.Enum

const (
	Burn UseMethod = iota
	Multiple
	Single
)

type CollectionDetails struct {
	Enum borsh.Enum `borsh_enum:"true"`
	V1   CollectionDetailsV1
}

type CollectionDetailsV1 struct {
	Size uint64
}

type ProgrammableConfig struct {
	Enum borsh.Enum `borsh_enum:"true"`
	V1   ProgrammableConfigV1
}

type ProgrammableConfigV1 struct {
	RuleSet *common.PublicKey
}

func MetadataDeserialize(data []byte) (Metadata, error) {
	var metadata Metadata
	err := borsh.Deserialize(&metadata, data)
	if err != nil {
		// https://github.com/samuelvanderwaal/metaboss/issues/121
		// https://github.com/metaplex-foundation/metaplex-program-library/pull/407
		// C.f. https://github.com/metaplex-foundation/metaplex-program-library/blob/master/token-metadata/program/src/deser.rs#L12
		var metadataPreV11 metadataPreV11
		err := borsh.Deserialize(&metadataPreV11, data)
		if err != nil {
			return Metadata{}, fmt.Errorf("failed to deserialize data, err: %v", err)
		} else {
			metadata.Key = metadataPreV11.Key
			metadata.UpdateAuthority = metadataPreV11.UpdateAuthority
			metadata.Mint = metadataPreV11.Mint
			metadata.Data = metadataPreV11.Data
			metadata.PrimarySaleHappened = metadataPreV11.PrimarySaleHappened
			metadata.IsMutable = metadataPreV11.IsMutable
			metadata.EditionNonce = metadataPreV11.EditionNonce
		}
	}
	// trim null byte
	metadata.Data.Name = strings.TrimRight(metadata.Data.Name, "\x00")
	metadata.Data.Symbol = strings.TrimRight(metadata.Data.Symbol, "\x00")
	metadata.Data.Uri = strings.TrimRight(metadata.Data.Uri, "\x00")
	return metadata, nil
}

type MasterEditionV2 struct {
	Key       Key
	Supply    uint64
	MaxSupply *uint64
}

type TokenMetadata struct {
	Mint            *common.PublicKey
	UpdateAuthority *common.PublicKey
	Name            string
	Symbol          string
	Uri             string
}

type Metadata2022 struct {
	token.MintAccount
	TokenMetadata
}

func readUInt16LE(aa []byte) uint16 {
	for i, j := 0, len(aa)-1; i < j; i, j = i+1, j-1 {
		aa[i], aa[j] = aa[j], aa[i]
	}
	return binary.BigEndian.Uint16(aa)
}

func Metadata2022Deserialize(data []byte) (*Metadata2022, error) {
	const (
		TypeSize = 2
		LengthSize
	)

	if len(data) < token.MintAccountSize {
		return nil, fmt.Errorf("invalid mint data size %d shoud >= %d", len(data), token.MintAccountSize)
	}

	mintAccount, err := token.MintAccountFromData(data[:token.MintAccountSize])
	if err != nil {
		return nil, fmt.Errorf("failed to parse data to a mint account, err: %w", err)
	}

	var tokenMetadata TokenMetadata

	extensionTypeIndex := token.MintAccountSize
	for extensionTypeIndex+TypeSize+LengthSize <= len(data) {
		entryType := readUInt16LE(data[extensionTypeIndex : extensionTypeIndex+TypeSize])
		extensionTypeIndex += TypeSize
		entryLength := readUInt16LE(data[extensionTypeIndex : extensionTypeIndex+LengthSize])
		extensionTypeIndex += LengthSize

		if entryType == 0 && entryLength == 256 {
			//这种情况不合乎规则，直接忽略
			continue
		}

		if entryType == 19 {
			insData := struct {
				UpdateAuthority common.PublicKey
				Mint            common.PublicKey
				Name            string
				Symbol          string
				Uri             string
			}{}
			err = borsh.Deserialize(&insData, data[extensionTypeIndex:])
			if err != nil {
				return nil, fmt.Errorf("borsh.Deserialize %w", err)
			}

			if !reflect.DeepEqual(common.SystemProgramID, insData.UpdateAuthority) {
				tokenMetadata.UpdateAuthority = &insData.UpdateAuthority
			}
			if !reflect.DeepEqual(common.SystemProgramID, insData.Mint) {
				tokenMetadata.Mint = &insData.Mint
			}
			tokenMetadata.Name = insData.Name
			tokenMetadata.Symbol = insData.Symbol
			tokenMetadata.Uri = insData.Uri
			break
		}
		extensionTypeIndex += int(entryLength)
	}

	return &Metadata2022{
		MintAccount:   mintAccount,
		TokenMetadata: tokenMetadata,
	}, nil
}
