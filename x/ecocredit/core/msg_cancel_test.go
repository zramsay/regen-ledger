package core

import (
	"testing"

	"github.com/regen-network/regen-ledger/types/testutil"
	"github.com/stretchr/testify/require"
)

func TestMsgCancel(t *testing.T) {
	t.Parallel()

	addr1 := testutil.GenAddress()

	tests := map[string]struct {
		src    MsgCancel
		expErr bool
	}{
		"valid msg": {
			src: MsgCancel{
				Owner: addr1,
				Credits: []*MsgCancel_CancelCredits{
					{
						BatchDenom: batchDenom,
						Amount:     "10",
					},
				},
				Reason: "reason",
			},
			expErr: false,
		},
		"invalid msg without holder": {
			src: MsgCancel{
				Credits: []*MsgCancel_CancelCredits{
					{
						BatchDenom: batchDenom,
						Amount:     "10",
					},
				},
			},
			expErr: true,
		},
		"invalid msg with wrong holder address": {
			src: MsgCancel{
				Owner: "wrong owner",
				Credits: []*MsgCancel_CancelCredits{
					{
						BatchDenom: batchDenom,
						Amount:     "10",
					},
				},
			},
			expErr: true,
		},
		"invalid msg without credits": {
			src: MsgCancel{
				Owner: addr1,
			},
			expErr: true,
		},
		"invalid msg without Credits.BatchDenom": {
			src: MsgCancel{
				Owner: addr1,
				Credits: []*MsgCancel_CancelCredits{
					{
						Amount: "10",
					},
				},
			},
			expErr: true,
		},
		"invalid msg without Credits.Amount": {
			src: MsgCancel{
				Owner: addr1,
				Credits: []*MsgCancel_CancelCredits{
					{
						BatchDenom: batchDenom,
					},
				},
			},
			expErr: true,
		},
		"invalid msg with wrong Credits.Amount": {
			src: MsgCancel{
				Owner: addr1,
				Credits: []*MsgCancel_CancelCredits{
					{
						BatchDenom: batchDenom,
						Amount:     "abc",
					},
				},
			},
			expErr: true,
		},
		"invalid msg reason is required": {
			src: MsgCancel{
				Owner: addr1,
				Credits: []*MsgCancel_CancelCredits{
					{
						BatchDenom: batchDenom,
						Amount:     "1",
					},
				},
			},
			expErr: true,
		},
	}

	for msg, test := range tests {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			err := test.src.ValidateBasic()
			if test.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
