package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

func RegisterInterfaces(registry types.InterfaceRegistry) {}

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
