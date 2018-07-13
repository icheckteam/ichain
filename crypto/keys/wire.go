package keys

import (
	amino "github.com/tendermint/go-amino"
	tcrypto "github.com/tendermint/tendermint/crypto"
)

var cdc = amino.NewCodec()

func init() {
	tcrypto.RegisterAmino(cdc)
	cdc.RegisterInterface((*Info)(nil), nil)
	cdc.RegisterConcrete(localInfo{}, "crypto/keys/localInfo", nil)
	cdc.RegisterConcrete(offlineInfo{}, "crypto/keys/offlineInfo", nil)
}
