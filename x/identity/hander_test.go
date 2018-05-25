package identity

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	addr := sdk.Address([]byte("input"))
	addr1 := sdk.Address([]byte("input1"))
	ctx, keeper := createTestInput(t, false)
	creatTime, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")
	expiration, _ := time.Parse(time.RFC3339Nano, "2018-05-11T16:28:45.78807557+07:00")
	var msg = MsgCreateClaim{
		ID:      "1212",
		Context: "claim:identity",
		Content: []byte(`{"demo": 1}`),
		Metadata: ClaimMetadata{
			CreateTime: creatTime,
			Expires:    expiration,
			Issuer:     addr,
			Recipient:  addr1,
		},
	}

	got := handleCreate(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected no error on TestHandle")
	claim, _ := keeper.GetClaim(ctx, msg.ID)
	require.True(t, claim != nil)

	got = handleRevokeMsg(ctx, keeper, MsgRevokeClaim{
		ClaimID:    claim.ID,
		Owner:      addr1,
		Revocation: "1212",
	})
	require.False(t, got.IsOK(), "expected no error on handleRevokeMsg")

	got = handleRevokeMsg(ctx, keeper, MsgRevokeClaim{
		ClaimID:    claim.ID,
		Owner:      addr,
		Revocation: "1212",
	})
	require.True(t, got.IsOK(), "expected no error on handleRevokeMsg")
	claim, _ = keeper.GetClaim(ctx, msg.ID)
	require.True(t, claim.Metadata.Revocation == "1212")
}
