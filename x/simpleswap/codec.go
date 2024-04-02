package simpleswap

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "rps/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgAddLiquidity{}, "simpleswap/MsgAddLiquidity")
	legacy.RegisterAminoMsg(cdc, &MsgSwapLiquidity{}, "simpleswap/MsgSwapLiquidity")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveLiquidity{}, "simpleswap/MsgRemoveLiquidity")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgAddLiquidity{},
		&MsgSwapLiquidity{},
		&MsgRemoveLiquidity{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
