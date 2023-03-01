// nolint:staticcheck
package v2

import (
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	migratetypes "github.com/functionx/fx-core/v3/x/migrate/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func GetModuleKey(chainId string) map[string]map[byte][2]int {
	return map[string]map[byte][2]int{
		bsctypes.ModuleName:     getBscKey(chainId),
		polygontypes.ModuleName: getPolygonKey(chainId),
		trontypes.ModuleName:    getTronKey(chainId),
		erc20types.ModuleName:   getErc20Key(chainId),
		migratetypes.ModuleName: getMigrateKey(chainId),
	}
}

func getBscKey(chainId string) map[byte][2]int {
	if chainId == fxtypes.TestnetChainId {
		return map[byte][2]int{
			crosschaintypes.OracleKey[0]:                          {2, 0},
			crosschaintypes.OracleAddressByExternalKey[0]:         {2, 0},
			crosschaintypes.OracleAddressByBridgerKey[0]:          {2, 0},
			crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
			crosschaintypes.OracleSetConfirmKey[0]:                {3, 0},
			crosschaintypes.OutgoingTxPoolKey[0]:                  {2, 0},
			crosschaintypes.OutgoingTxBatchKey[0]:                 {0, 0},
			crosschaintypes.OutgoingTxBatchBlockKey[0]:            {0, 0},
			crosschaintypes.LastEventNonceByOracleKey[0]:          {2, 0},
			crosschaintypes.LastObservedEventNonceKey[0]:          {1, 0},
			crosschaintypes.SequenceKeyPrefix[0]:                  {2, 0},
			crosschaintypes.DenomToTokenKey[0]:                    {2, 0},
			crosschaintypes.TokenToDenomKey[0]:                    {2, 0},
			crosschaintypes.LastSlashedOracleSetNonce[0]:          {1, 0},
			crosschaintypes.LatestOracleSetNonce[0]:               {1, 0},
			crosschaintypes.LastSlashedBatchBlock[0]:              {0, 0},
			crosschaintypes.LastObservedBlockHeightKey[0]:         {1, 0},
			crosschaintypes.LastObservedOracleSetKey[0]:           {1, 0},
			crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {2, 0},
			crosschaintypes.LastOracleSlashBlockHeight[0]:         {0, 0},
			crosschaintypes.ProposalOracleKey[0]:                  {1, 0},
			crosschaintypes.LastTotalPowerKey[0]:                  {1, 0},
			crosschaintypes.OracleAttestationKey[0]:               {101, 0},
			crosschaintypes.PastExternalSignatureCheckpointKey[0]: {685, 0},
			crosschaintypes.BatchConfirmKey[0]:                    {1366, 0},
			crosschaintypes.LastProposalBlockHeight[0]:            {1, 0},
		}
	}
	return map[byte][2]int{
		crosschaintypes.OracleKey[0]:                          {2, 0},
		crosschaintypes.OracleAddressByExternalKey[0]:         {2, 0},
		crosschaintypes.OracleAddressByBridgerKey[0]:          {2, 0},
		crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
		crosschaintypes.OracleSetConfirmKey[0]:                {3, 0},
		crosschaintypes.OutgoingTxPoolKey[0]:                  {0, 0},
		crosschaintypes.OutgoingTxBatchKey[0]:                 {0, 0},
		crosschaintypes.OutgoingTxBatchBlockKey[0]:            {0, 0},
		crosschaintypes.LastEventNonceByOracleKey[0]:          {2, 0},
		crosschaintypes.LastObservedEventNonceKey[0]:          {1, 0},
		crosschaintypes.SequenceKeyPrefix[0]:                  {2, 0},
		crosschaintypes.DenomToTokenKey[0]:                    {1, 0},
		crosschaintypes.TokenToDenomKey[0]:                    {1, 0},
		crosschaintypes.LastSlashedOracleSetNonce[0]:          {1, 0},
		crosschaintypes.LatestOracleSetNonce[0]:               {1, 0},
		crosschaintypes.LastSlashedBatchBlock[0]:              {0, 0},
		crosschaintypes.LastObservedBlockHeightKey[0]:         {1, 0},
		crosschaintypes.LastObservedOracleSetKey[0]:           {1, 0},
		crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {2, 0},
		crosschaintypes.LastOracleSlashBlockHeight[0]:         {0, 0},
		crosschaintypes.ProposalOracleKey[0]:                  {1, 0},
		crosschaintypes.LastTotalPowerKey[0]:                  {1, 0},
		crosschaintypes.OracleAttestationKey[0]:               {101, 0},
		crosschaintypes.PastExternalSignatureCheckpointKey[0]: {505, 0},
		crosschaintypes.BatchConfirmKey[0]:                    {1006, 0},
		crosschaintypes.LastProposalBlockHeight[0]:            {1, 0},
	}
}

func getPolygonKey(chainId string) map[byte][2]int {
	if chainId == fxtypes.TestnetChainId {
		return map[byte][2]int{
			crosschaintypes.OracleKey[0]:                          {5, 0},
			crosschaintypes.OracleAddressByExternalKey[0]:         {5, 0},
			crosschaintypes.OracleAddressByBridgerKey[0]:          {5, 0},
			crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
			crosschaintypes.OracleSetConfirmKey[0]:                {15, 0},
			crosschaintypes.OutgoingTxPoolKey[0]:                  {2, 0},
			crosschaintypes.OutgoingTxBatchKey[0]:                 {0, 0},
			crosschaintypes.OutgoingTxBatchBlockKey[0]:            {0, 0},
			crosschaintypes.LastEventNonceByOracleKey[0]:          {5, 0},
			crosschaintypes.LastObservedEventNonceKey[0]:          {1, 0},
			crosschaintypes.SequenceKeyPrefix[0]:                  {2, 0},
			crosschaintypes.DenomToTokenKey[0]:                    {2, 0},
			crosschaintypes.TokenToDenomKey[0]:                    {2, 0},
			crosschaintypes.LastSlashedOracleSetNonce[0]:          {1, 0},
			crosschaintypes.LatestOracleSetNonce[0]:               {1, 0},
			crosschaintypes.LastSlashedBatchBlock[0]:              {1, 0},
			crosschaintypes.LastObservedBlockHeightKey[0]:         {1, 0},
			crosschaintypes.LastObservedOracleSetKey[0]:           {1, 0},
			crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {5, 0},
			crosschaintypes.LastOracleSlashBlockHeight[0]:         {0, 0},
			crosschaintypes.ProposalOracleKey[0]:                  {1, 0},
			crosschaintypes.LastTotalPowerKey[0]:                  {1, 0},
			crosschaintypes.OracleAttestationKey[0]:               {101, 0},
			crosschaintypes.PastExternalSignatureCheckpointKey[0]: {608, 0},
			crosschaintypes.BatchConfirmKey[0]:                    {3015, 0},
			crosschaintypes.LastProposalBlockHeight[0]:            {1, 0},
		}
	}
	return map[byte][2]int{
		crosschaintypes.OracleKey[0]:                          {10, 0},
		crosschaintypes.OracleAddressByExternalKey[0]:         {10, 0},
		crosschaintypes.OracleAddressByBridgerKey[0]:          {10, 0},
		crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
		crosschaintypes.OracleSetConfirmKey[0]:                {15, 0},
		crosschaintypes.OutgoingTxPoolKey[0]:                  {0, 0},
		crosschaintypes.OutgoingTxBatchKey[0]:                 {0, 0},
		crosschaintypes.OutgoingTxBatchBlockKey[0]:            {0, 0},
		crosschaintypes.LastEventNonceByOracleKey[0]:          {10, 0},
		crosschaintypes.LastObservedEventNonceKey[0]:          {1, 0},
		crosschaintypes.SequenceKeyPrefix[0]:                  {2, 0},
		crosschaintypes.DenomToTokenKey[0]:                    {1, 0},
		crosschaintypes.TokenToDenomKey[0]:                    {1, 0},
		crosschaintypes.LastSlashedOracleSetNonce[0]:          {1, 0},
		crosschaintypes.LatestOracleSetNonce[0]:               {1, 0},
		crosschaintypes.LastSlashedBatchBlock[0]:              {0, 0},
		crosschaintypes.LastObservedBlockHeightKey[0]:         {1, 0},
		crosschaintypes.LastObservedOracleSetKey[0]:           {1, 0},
		crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {10, 0},
		crosschaintypes.LastOracleSlashBlockHeight[0]:         {0, 0},
		crosschaintypes.ProposalOracleKey[0]:                  {1, 0},
		crosschaintypes.LastTotalPowerKey[0]:                  {1, 0},
		crosschaintypes.OracleAttestationKey[0]:               {103, 0},
		crosschaintypes.PastExternalSignatureCheckpointKey[0]: {80, 0},
		crosschaintypes.BatchConfirmKey[0]:                    {700, 0},
		crosschaintypes.LastProposalBlockHeight[0]:            {1, 0},
	}
}

func getTronKey(chainId string) map[byte][2]int {
	if chainId == fxtypes.TestnetChainId {
		return map[byte][2]int{
			crosschaintypes.OracleKey[0]:                          {5, 0},
			crosschaintypes.OracleAddressByExternalKey[0]:         {5, 0},
			crosschaintypes.OracleAddressByBridgerKey[0]:          {5, 0},
			crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
			crosschaintypes.OracleSetConfirmKey[0]:                {15, 0},
			crosschaintypes.OutgoingTxPoolKey[0]:                  {1, 0},
			crosschaintypes.OutgoingTxBatchKey[0]:                 {0, 0},
			crosschaintypes.OutgoingTxBatchBlockKey[0]:            {0, 0},
			crosschaintypes.LastEventNonceByOracleKey[0]:          {5, 0},
			crosschaintypes.LastObservedEventNonceKey[0]:          {1, 0},
			crosschaintypes.SequenceKeyPrefix[0]:                  {2, 0},
			crosschaintypes.DenomToTokenKey[0]:                    {4, 0},
			crosschaintypes.TokenToDenomKey[0]:                    {4, 0},
			crosschaintypes.LastSlashedOracleSetNonce[0]:          {1, 0},
			crosschaintypes.LatestOracleSetNonce[0]:               {1, 0},
			crosschaintypes.LastSlashedBatchBlock[0]:              {1, 0},
			crosschaintypes.LastObservedBlockHeightKey[0]:         {1, 0},
			crosschaintypes.LastObservedOracleSetKey[0]:           {1, 0},
			crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {5, 0},
			crosschaintypes.LastOracleSlashBlockHeight[0]:         {0, 0},
			crosschaintypes.ProposalOracleKey[0]:                  {1, 0},
			crosschaintypes.LastTotalPowerKey[0]:                  {1, 0},
			crosschaintypes.OracleAttestationKey[0]:               {101, 0},
			crosschaintypes.PastExternalSignatureCheckpointKey[0]: {625, 0},
			crosschaintypes.BatchConfirmKey[0]:                    {3100, 0},
			crosschaintypes.LastProposalBlockHeight[0]:            {1, 0},
		}
	}
	return map[byte][2]int{
		crosschaintypes.OracleKey[0]:                          {10, 0},
		crosschaintypes.OracleAddressByExternalKey[0]:         {10, 0},
		crosschaintypes.OracleAddressByBridgerKey[0]:          {10, 0},
		crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
		crosschaintypes.OracleSetConfirmKey[0]:                {55, 0},
		crosschaintypes.OutgoingTxPoolKey[0]:                  {0, 0},
		crosschaintypes.OutgoingTxBatchKey[0]:                 {0, 0},
		crosschaintypes.OutgoingTxBatchBlockKey[0]:            {0, 0},
		crosschaintypes.LastEventNonceByOracleKey[0]:          {10, 0},
		crosschaintypes.LastObservedEventNonceKey[0]:          {1, 0},
		crosschaintypes.SequenceKeyPrefix[0]:                  {2, 0},
		crosschaintypes.DenomToTokenKey[0]:                    {1, 0},
		crosschaintypes.TokenToDenomKey[0]:                    {1, 0},
		crosschaintypes.LastSlashedOracleSetNonce[0]:          {1, 0},
		crosschaintypes.LatestOracleSetNonce[0]:               {1, 0},
		crosschaintypes.LastSlashedBatchBlock[0]:              {1, 0},
		crosschaintypes.LastObservedBlockHeightKey[0]:         {1, 0},
		crosschaintypes.LastObservedOracleSetKey[0]:           {1, 0},
		crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {10, 0},
		crosschaintypes.LastOracleSlashBlockHeight[0]:         {0, 0},
		crosschaintypes.ProposalOracleKey[0]:                  {1, 0},
		crosschaintypes.LastTotalPowerKey[0]:                  {1, 0},
		crosschaintypes.OracleAttestationKey[0]:               {103, 0},
		crosschaintypes.PastExternalSignatureCheckpointKey[0]: {55, 0},
		crosschaintypes.BatchConfirmKey[0]:                    {450, 0},
		crosschaintypes.LastProposalBlockHeight[0]:            {1, 0},
	}
}

func getErc20Key(chainId string) map[byte][2]int {
	if chainId == fxtypes.TestnetChainId {
		return map[byte][2]int{
			erc20types.KeyPrefixTokenPair[0]:        {15, 0},
			erc20types.KeyPrefixTokenPairByERC20[0]: {15, 0},
			erc20types.KeyPrefixTokenPairByDenom[0]: {15, 0},
			erc20types.KeyPrefixIBCTransfer[0]:      {513, 0},
			erc20types.KeyPrefixAliasDenom[0]:       {8, 0},
		}
	}
	return map[byte][2]int{
		erc20types.KeyPrefixTokenPair[0]:        {12, 0},
		erc20types.KeyPrefixTokenPairByERC20[0]: {12, 0},
		erc20types.KeyPrefixTokenPairByDenom[0]: {12, 0},
		erc20types.KeyPrefixIBCTransfer[0]:      {182, 0},
		erc20types.KeyPrefixAliasDenom[0]:       {11, 0},
	}
}

func getMigrateKey(chainId string) map[byte][2]int {
	if chainId == fxtypes.TestnetChainId {
		return map[byte][2]int{
			migratetypes.KeyPrefixMigratedRecord[0]:        {332, 0},
			migratetypes.KeyPrefixMigratedDirectionFrom[0]: {166, 0},
			migratetypes.KeyPrefixMigratedDirectionTo[0]:   {166, 0},
		}
	}
	return map[byte][2]int{
		migratetypes.KeyPrefixMigratedRecord[0]:        {1242, 0},
		migratetypes.KeyPrefixMigratedDirectionFrom[0]: {621, 0},
		migratetypes.KeyPrefixMigratedDirectionTo[0]:   {621, 0},
	}
}
