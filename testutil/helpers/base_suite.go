package helpers

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tenderminttypes "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	localhost "github.com/cosmos/ibc-go/v8/modules/light-clients/09-localhost"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/app"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
)

type BaseSuite struct {
	suite.Suite
	MintValNumber int
	ValSet        *tmtypes.ValidatorSet
	ValAddr       []sdk.ValAddress
	App           *app.App
	Ctx           sdk.Context
}

func (s *BaseSuite) SetupTest() {
	valNumber := s.MintValNumber
	if valNumber <= 0 {
		valNumber = tmrand.Intn(10) + 1
	}
	valSet, valAccounts, valBalances := generateGenesisValidator(valNumber, sdk.Coins{})
	s.ValSet = valSet
	s.ValAddr = make([]sdk.ValAddress, valNumber)
	for i := 0; i < valNumber; i++ {
		s.ValAddr[i] = valAccounts[i].GetAddress().Bytes()
	}

	s.App = setupWithGenesisValSet(s.T(), valSet, valAccounts, valBalances...)
	s.Ctx = s.App.GetContextForFinalizeBlock(nil)
	s.Ctx = s.Ctx.WithProposer(s.ValSet.Proposer.Address.Bytes())
}

func (s *BaseSuite) AddTestSigner(amount ...int64) *Signer {
	signer := NewSigner(NewEthPrivKey())
	defAmount := int64(10000)
	if len(amount) > 0 {
		defAmount = amount[0]
	}
	s.MintToken(signer.AccAddress(), NewStakingCoin(defAmount, 18))
	return signer
}

func (s *BaseSuite) NewSigner() *Signer {
	signer := NewSigner(NewEthPrivKey())
	account := s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, signer.AccAddress())
	s.App.AccountKeeper.SetAccount(s.Ctx, account)
	return signer
}

func (s *BaseSuite) Commit(block ...int64) sdk.Context {
	ctx := s.Ctx
	lastBlockHeight := s.Ctx.BlockHeight()
	nextHeight := lastBlockHeight + 1
	if len(block) > 0 {
		nextHeight = lastBlockHeight + block[0]
	}
	commitInfo := abci.CommitInfo{
		Round: 1,
	}

	for _, val := range s.ValSet.Validators {
		pk, err := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		s.Require().NoError(err)
		commitInfo.Votes = append(commitInfo.Votes, abci.VoteInfo{
			Validator: abci.Validator{
				Address: pk.Address(),
				Power:   val.VotingPower,
			},
			BlockIdFlag: tenderminttypes.BlockIDFlagCommit,
		})

		signingInfo := slashingtypes.NewValidatorSigningInfo(
			sdk.ConsAddress(pk.Address()),
			ctx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
		s.Require().NoError(s.App.SlashingKeeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()), signingInfo))
	}

	for i := lastBlockHeight; i < nextHeight; i++ {
		// 1. try to finalize the block + commit finalizeBlockState
		if _, err := s.App.FinalizeBlock(&abci.RequestFinalizeBlock{
			Height:            i,
			Time:              tmtime.Now(),
			ProposerAddress:   s.Ctx.BlockHeader().ProposerAddress,
			DecidedLastCommit: commitInfo,
		}); err != nil {
			panic(err)
		}

		// 2. commit lastCommitInfo
		if _, err := s.App.Commit(); err != nil {
			panic(err)
		}

		// 3. prepare to process new blocks (myApp.GetContextForFinalizeBlock(nil))
		if _, err := s.App.ProcessProposal(&abci.RequestProcessProposal{
			Height:             i + 1,
			Time:               tmtime.Now(),
			ProposerAddress:    s.Ctx.BlockHeader().ProposerAddress,
			ProposedLastCommit: commitInfo,
		}); err != nil {
			panic(err)
		}

		// 4. get new ctx for finalizeBlockState
		ctx = s.App.GetContextForFinalizeBlock(nil)
	}
	s.Ctx = ctx
	return ctx
}

func (s *BaseSuite) AddTestSigners(accNum int, coin sdk.Coin) []*Signer {
	signers := make([]*Signer, accNum)
	for i := 0; i < accNum; i++ {
		signers[i] = NewSigner(NewEthPrivKey())
		s.MintToken(signers[i].AccAddress(), coin)
	}
	return signers
}

func (s *BaseSuite) AddTestAddress(accNum int, coin sdk.Coin) []sdk.AccAddress {
	accAddresses := make([]sdk.AccAddress, accNum)
	signers := s.AddTestSigners(accNum, coin)
	for i := 0; i < accNum; i++ {
		accAddresses[i] = signers[i].AccAddress()
	}
	return accAddresses
}

func (s *BaseSuite) MintToken(address sdk.AccAddress, amount ...sdk.Coin) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, sdk.NewCoins(amount...))
	s.Require().NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, address.Bytes(), sdk.NewCoins(amount...))
	s.Require().NoError(err)
}

func (s *BaseSuite) MintTokenToModule(module string, amount ...sdk.Coin) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, module, sdk.NewCoins(amount...))
	s.Require().NoError(err)
}

func (s *BaseSuite) Balance(acc sdk.AccAddress) sdk.Coins {
	return s.App.BankKeeper.GetAllBalances(s.Ctx, acc)
}

func (s *BaseSuite) CheckBalance(addr sdk.AccAddress, expBal ...sdk.Coin) {
	balances := s.App.BankKeeper.GetAllBalances(s.Ctx, addr)
	for _, bal := range expBal {
		s.Equal(bal.Amount, balances.AmountOf(bal.Denom), bal.Denom)
	}
}

func (s *BaseSuite) CheckAllBalance(addr sdk.AccAddress, expBal ...sdk.Coin) {
	balances := s.App.BankKeeper.GetAllBalances(s.Ctx, addr)
	s.Equal(sdk.NewCoins(expBal...).String(), balances.String())
}

func (s *BaseSuite) GetStakingBalance(addr sdk.AccAddress) sdkmath.Int {
	balances := s.App.BankKeeper.GetAllBalances(s.Ctx, addr)
	return balances.AmountOf(fxtypes.DefaultDenom)
}

func (s *BaseSuite) NewCoin(amounts ...sdkmath.Int) sdk.Coin {
	amount := NewRandAmount()
	if len(amounts) > 0 && amounts[0].IsPositive() {
		amount = amounts[0]
	}
	denom := NewRandDenom()
	return sdk.NewCoin(denom, amount)
}

func (s *BaseSuite) NewBridgeCoin(module string, amounts ...sdkmath.Int) (sdk.Coin, string) {
	amount := NewRandAmount()
	if len(amounts) > 0 && amounts[0].IsPositive() {
		amount = amounts[0]
	}
	tokenAddr := GenExternalAddr(module)
	bridgeDenom := crosschaintypes.NewBridgeDenom(module, tokenAddr)
	return sdk.NewCoin(bridgeDenom, amount), tokenAddr
}

func (s *BaseSuite) GenIBCTransferChannel() (portID, channelID string) {
	portID = "transfer"

	channelSequence := s.App.IBCKeeper.ChannelKeeper.GetNextChannelSequence(s.Ctx)
	channelID = fmt.Sprintf("channel-%d", channelSequence)
	connectionID := connectiontypes.FormatConnectionIdentifier(uint64(tmrand.Intn(100)))
	clientID := clienttypes.FormatClientIdentifier(exported.Localhost, uint64(tmrand.Intn(100)))

	revision := clienttypes.ParseChainID(s.Ctx.ChainID())
	localHostClient := localhost.NewClientState(clienttypes.NewHeight(revision, uint64(s.Ctx.BlockHeight())))
	s.App.IBCKeeper.ClientKeeper.SetClientState(s.Ctx, clientID, localHostClient)

	params := s.App.IBCKeeper.ClientKeeper.GetParams(s.Ctx)
	params.AllowedClients = append(params.AllowedClients, localHostClient.ClientType())
	s.App.IBCKeeper.ClientKeeper.SetParams(s.Ctx, params)

	prevConsState := &ibctm.ConsensusState{
		Timestamp:          s.Ctx.BlockTime(),
		NextValidatorsHash: s.Ctx.BlockHeader().NextValidatorsHash,
	}
	height := clienttypes.NewHeight(0, uint64(s.Ctx.BlockHeight()))
	s.App.IBCKeeper.ClientKeeper.SetClientConsensusState(s.Ctx, clientID, height, prevConsState)

	channelCapability, err := s.App.ScopedIBCKeeper.NewCapability(s.Ctx, host.ChannelCapabilityPath(portID, channelID))
	s.Require().NoError(err)
	err = s.App.ScopedTransferKeeper.ClaimCapability(s.Ctx, capabilitytypes.NewCapability(channelCapability.Index), host.ChannelCapabilityPath(portID, channelID))
	s.Require().NoError(err)

	connectionEnd := connectiontypes.NewConnectionEnd(connectiontypes.OPEN, clientID, connectiontypes.Counterparty{ClientId: "clientId", ConnectionId: "connection-1", Prefix: commitmenttypes.NewMerklePrefix([]byte("prefix"))}, connectiontypes.GetCompatibleVersions(), 500)
	s.App.IBCKeeper.ConnectionKeeper.SetConnection(s.Ctx, connectionID, connectionEnd)

	channel := channeltypes.NewChannel(channeltypes.OPEN, channeltypes.ORDERED, channeltypes.NewCounterparty(portID, channelID), []string{connectionID}, "mock-version")
	s.App.IBCKeeper.ChannelKeeper.SetChannel(s.Ctx, portID, channelID, channel)
	s.App.IBCKeeper.ChannelKeeper.SetNextSequenceSend(s.Ctx, portID, channelID, uint64(tmrand.Intn(10000)+1))
	s.App.IBCKeeper.ChannelKeeper.SetNextChannelSequence(s.Ctx, channelSequence+1)
	return portID, channelID
}
