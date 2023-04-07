package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/public-awesome/stargaze/v9/x/globalfee/types"
)

func (k Keeper) IterateCodeAuthorizations(ctx sdk.Context, cb func(types.CodeAuthorization) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.CodeAuthorizationPrefix)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ca types.CodeAuthorization
		k.cdc.MustUnmarshal(iter.Value(), &ca)
		// cb returns true to stop early
		if cb(ca) {
			return
		}
	}
}

func (k Keeper) GetCodeAuthorization(ctx sdk.Context, codeId uint64) (types.CodeAuthorization, bool) {
	store := ctx.KVStore(k.storeKey)

	var ca types.CodeAuthorization
	bz := store.Get(types.GetCodeAuthorizationPrefix(codeId))
	if bz == nil {
		return ca, false
	}

	k.cdc.MustUnmarshal(bz, &ca)
	return ca, true
}

func (k Keeper) SetCodeAuthorization(ctx sdk.Context, ca types.CodeAuthorization) error {
	if err := ca.Validate(); err != nil {
		return err
	}

	if k.wasmKeeper.GetCodeInfo(ctx, ca.GetCodeId()) == nil {
		return types.ErrCodeIdNotExist
	}

	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshal(&ca)

	store.Set(types.GetCodeAuthorizationPrefix(ca.CodeId), value)
	return nil
}

func (k Keeper) DeleteCodeAuthorization(ctx sdk.Context, codeId uint64) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetCodeAuthorizationPrefix(codeId))
}
