package tests

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/app"
	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	fxgovtypes "github.com/functionx/fx-core/v7/x/gov/types"
	ibctransfertypes "github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
)

func TestAminoEncode(t *testing.T) {
	_ = app.MakeEncodingConfig()
	oneDec := sdk.NewDec(1)
	oneInt := sdk.NewInt(1)

	testcases := []struct {
		name     string
		expected string
		msg      interface{}
		amino    *codec.LegacyAmino
	}{
		{
			name:     "upgrade-SoftwareUpgradeProposal",
			expected: `{"type":"cosmos-sdk/MsgSubmitProposal","value":{"content":{"type":"cosmos-sdk/SoftwareUpgradeProposal","value":{"description":"foo","plan":{"height":"123","info":"foo","name":"foo","time":"0001-01-01T00:00:00Z"},"title":"v2"}},"initial_deposit":[]}}`,
			amino:    govv1betal.ModuleCdc.LegacyAmino,
			msg: govv1betal.MsgSubmitProposal{
				Content: NewAnyWithValue(&upgradetypes.SoftwareUpgradeProposal{ // nolint:staticcheck
					Title:       "v2",
					Description: "foo",
					Plan: upgradetypes.Plan{
						Name:   "foo",
						Time:   time.Time{},
						Height: 123,
						Info:   "foo",
					},
				}),
				InitialDeposit: nil,
				Proposer:       "",
			},
		},
		{
			name:     "upgrade-MsgSoftwareUpgrade",
			expected: `{"type":"cosmos-sdk/MsgSubmitProposal","value":{"content":{"type":"cosmos-sdk/MsgSoftwareUpgrade","value":{"authority":"foo","plan":{"height":"123","info":"foo","name":"foo","time":"0001-01-01T00:00:00Z"}}},"initial_deposit":[]}}`,
			amino:    govv1betal.ModuleCdc.LegacyAmino,
			msg: govv1betal.MsgSubmitProposal{
				Content: NewAnyWithValue(&upgradetypes.MsgSoftwareUpgrade{
					Authority: "foo",
					Plan: upgradetypes.Plan{
						Name:   "foo",
						Time:   time.Time{},
						Height: 123,
						Info:   "foo",
					},
				}),
				InitialDeposit: nil,
				Proposer:       "",
			},
		},
		{
			name:     "ibc-MsgTransfer",
			expected: `{"type":"fxtransfer/MsgTransfer","value":{"fee":{"amount":"0","denom":"FX"},"receiver":"0x001","sender":"fx1001","source_channel":"channel-0","source_port":"transfer","timeout_height":{},"timeout_timestamp":"1675063442000000000","token":{"amount":"1","denom":"FX"}}}`,
			amino:    ibctransfertypes.AminoCdc.LegacyAmino,
			msg: ibctransfertypes.MsgTransfer{
				SourcePort:       "transfer",
				SourceChannel:    "channel-0",
				Token:            sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1)),
				Sender:           "fx1001",
				Receiver:         "0x001",
				TimeoutHeight:    clienttypes.Height{},
				TimeoutTimestamp: 1675063442000000000,
				Router:           "",
				Fee:              sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.ZeroInt()),
				Memo:             "",
			},
		},
		{
			name:     "erc20-RegisterCoinProposal",
			expected: `{"type":"cosmos-sdk/MsgSubmitProposal","value":{"content":{"type":"erc20/RegisterCoinProposal","value":{"description":"foo","metadata":{"base":"test","denom_units":[{"aliases":["ethtest"],"denom":"test"},{"denom":"TEST","exponent":18}],"description":"test","display":"test","name":"test name","symbol":"TEST"},"title":"v2"}},"initial_deposit":[]}}`,
			amino:    govv1betal.ModuleCdc.LegacyAmino,
			msg: govv1betal.MsgSubmitProposal{
				Content: NewAnyWithValue(&erc20types.RegisterCoinProposal{ // nolint:staticcheck
					Title:       "v2",
					Description: "foo",
					Metadata: types.Metadata{
						Description: "test",
						DenomUnits: []*types.DenomUnit{
							{
								Denom:    "test",
								Exponent: 0,
								Aliases: []string{
									"ethtest",
								},
							},
							{
								Denom:    "TEST",
								Exponent: 18,
								Aliases:  []string{},
							},
						},
						Base:    "test",
						Display: "test",
						Name:    "test name",
						Symbol:  "TEST",
					},
				}),
				InitialDeposit: nil,
				Proposer:       "",
			},
		},
		{
			name:     "staking-MsgEditValidator",
			expected: `{"type":"cosmos-sdk/MsgEditValidator","value":{"commission_rate":"1.000000000000000000","description":{"details":"foo","identity":"foo","moniker":"foo","security_contact":"foo","website":"foo"},"min_self_delegation":"1","validator_address":"cosmosvaloper1h6lrm4uusd46tu4slkg620hylv46hhff7a8su6"}}`,
			amino:    stakingtypes.ModuleCdc.LegacyAmino,
			msg: stakingtypes.MsgEditValidator{
				Description: stakingtypes.Description{
					Moniker:         "foo",
					Identity:        "foo",
					Website:         "foo",
					SecurityContact: "foo",
					Details:         "foo",
				},
				ValidatorAddress:  "cosmosvaloper1h6lrm4uusd46tu4slkg620hylv46hhff7a8su6",
				CommissionRate:    &oneDec,
				MinSelfDelegation: &oneInt,
			},
		},
		{
			name:     "staking-MsgEditValidator",
			expected: `{"type":"cosmos-sdk/MsgEditValidator","value":{"description":{"details":"foo","moniker":"foo","security_contact":"foo","website":"foo"},"validator_address":"cosmosvaloper1h6lrm4uusd46tu4slkg620hylv46hhff7a8su6"}}`,
			amino:    stakingtypes.ModuleCdc.LegacyAmino,
			msg: stakingtypes.MsgEditValidator{
				Description: stakingtypes.Description{
					Moniker:         "foo",
					Identity:        "",
					Website:         "foo",
					SecurityContact: "foo",
					Details:         "foo",
				},
				ValidatorAddress:  "cosmosvaloper1h6lrm4uusd46tu4slkg620hylv46hhff7a8su6",
				CommissionRate:    nil,
				MinSelfDelegation: nil,
			},
		},
		{
			name:     "gov-TextProposal",
			expected: `{"type":"cosmos-sdk/MsgSubmitProposal","value":{"content":{"type":"cosmos-sdk/TextProposal","value":{"description":"foo desc","title":"foo title"}},"initial_deposit":[]}}`,
			amino:    govv1betal.ModuleCdc.LegacyAmino,
			msg: govv1betal.MsgSubmitProposal{
				Content: NewAnyWithValue(&govv1betal.TextProposal{
					Title:       "foo title",
					Description: "foo desc",
				}),
				InitialDeposit: nil,
				Proposer:       "",
			},
		},
		{
			name:     "gov-v1-MsgSubmitProposal-crosschain-MsgUpdateParams",
			expected: `{"type":"cosmos-sdk/v1/MsgSubmitProposal","value":{"initial_deposit":null,"messages":[{"type":"crosschain/MsgUpdateParams","value":{"authority":"1","chain_name":"1","params":{"average_block_time":"1","average_external_block_time":"1","delegate_multiple":"1","delegate_threshold":{"amount":"1","denom":"FX"},"external_batch_timeout":"1","gravity_id":"1","ibc_transfer_timeout_height":"1","oracle_set_update_power_change_percent":"1.000000000000000000","signed_window":"1","slash_fraction":"1.000000000000000000"}}}]}}`,
			amino:    govv1.ModuleCdc.LegacyAmino,
			msg: govv1.MsgSubmitProposal{
				Messages: []*codectypes.Any{
					NewAnyWithValue(&crosschaintypes.MsgUpdateParams{
						ChainName: "1",
						Authority: "1",
						Params: crosschaintypes.Params{
							GravityId:                         "1",
							AverageBlockTime:                  1,
							ExternalBatchTimeout:              1,
							AverageExternalBlockTime:          1,
							SignedWindow:                      1,
							SlashFraction:                     sdk.NewDec(1),
							OracleSetUpdatePowerChangePercent: sdk.NewDec(1),
							IbcTransferTimeoutHeight:          1,
							Oracles:                           nil,
							DelegateThreshold:                 sdk.NewCoin("FX", sdk.NewInt(1)),
							DelegateMultiple:                  1,
						},
					}),
				},
				InitialDeposit: nil,
				Proposer:       "",
				Metadata:       "",
			},
		},
		{
			name:     "gov-v1-MsgSubmitProposal-erc20-MsgUpdateParams",
			expected: `{"type":"cosmos-sdk/v1/MsgSubmitProposal","value":{"initial_deposit":null,"messages":[{"type":"erc20/MsgUpdateParams","value":{"authority":"1","params":{"enable_erc20":true,"enable_evm_hook":true,"ibc_timeout":"1"}}}]}}`,
			amino:    govv1.ModuleCdc.LegacyAmino,
			msg: govv1.MsgSubmitProposal{
				Messages: []*codectypes.Any{
					NewAnyWithValue(&erc20types.MsgUpdateParams{
						Authority: "1",
						Params: erc20types.Params{
							EnableErc20:   true,
							EnableEVMHook: true,
							IbcTimeout:    1,
						},
					}),
				},
				InitialDeposit: nil,
				Proposer:       "",
				Metadata:       "",
			},
		},
		{
			name:     "gov-v1-MsgSubmitProposal-gov-MsgUpdateParams",
			expected: `{"type":"cosmos-sdk/v1/MsgSubmitProposal","value":{"initial_deposit":null,"messages":[{"type":"gov/MsgUpdateParams","value":{"authority":"1","params":{"max_deposit_period":"172800000000000","min_deposit":[{"amount":"10000000","denom":"stake"}],"min_initial_deposit":{"amount":"1000000000000000000000","denom":"FX"},"msg_type":"/fx.evm.v1.MsgCallContract","quorum":"0.334000000000000000","threshold":"0.500000000000000000","veto_threshold":"0.334000000000000000","voting_period":"172800000000000"}}}]}}`,
			amino:    govv1.ModuleCdc.LegacyAmino,
			msg: govv1.MsgSubmitProposal{
				Messages: []*codectypes.Any{
					NewAnyWithValue(&fxgovtypes.MsgUpdateParams{
						Authority: "1",
						Params:    *fxgovtypes.DefaultParams(),
					}),
				},
				InitialDeposit: nil,
				Proposer:       "",
				Metadata:       "",
			},
		},
		{
			name:     "gov-v1-MsgSubmitProposal-gov-MsgUpdateEGFParams",
			expected: `{"type":"cosmos-sdk/v1/MsgSubmitProposal","value":{"initial_deposit":null,"messages":[{"type":"gov/MsgUpdateEGFParams","value":{"authority":"1","params":{"claim_ratio":"0.100000000000000000","egf_deposit_threshold":{"amount":"10000000000000000000000","denom":"FX"}}}}]}}`,
			amino:    govv1.ModuleCdc.LegacyAmino,
			msg: govv1.MsgSubmitProposal{
				Messages: []*codectypes.Any{
					NewAnyWithValue(&fxgovtypes.MsgUpdateEGFParams{
						Authority: "1",
						Params:    *fxgovtypes.DefaultEGFParams(),
					}),
				},
				InitialDeposit: nil,
				Proposer:       "",
				Metadata:       "",
			},
		},
		{
			name:     "gov-v1-MsgVote",
			expected: `{"type":"cosmos-sdk/v1/MsgVote","value":{"metadata":"foo","option":1,"proposal_id":"1","voter":"foo"}}`,
			amino:    govv1.ModuleCdc.LegacyAmino,
			msg: govv1.MsgVote{
				ProposalId: 1,
				Voter:      "foo",
				Option:     1,
				Metadata:   "foo",
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			aminoJson, err := testcase.amino.MarshalJSON(testcase.msg)
			require.NoError(t, err)
			require.Equal(t, testcase.expected, string(sdk.MustSortJSON(aminoJson)))
		})
	}
}

func NewAnyWithValue(msg proto.Message) *codectypes.Any {
	value, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}
	return value
}
