package v3

import (
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
	migratetypes "github.com/functionx/fx-core/v3/x/migrate/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func GetModuleKey() map[string]map[byte][2]int {
	return map[string]map[byte][2]int{
		gravitytypes.ModuleName:   getGravityKey(),
		bsctypes.ModuleName:       getBscKey(),
		polygontypes.ModuleName:   getPolygonKey(),
		trontypes.ModuleName:      getTronKey(),
		ethtypes.ModuleName:       getEthKey(),
		avalanchetypes.ModuleName: getAvalancheKey(),
		erc20types.ModuleName:     getErc20Key(),
		migratetypes.ModuleName:   getMigrateKey(),
	}
}

func getGravityKey() map[byte][2]int {
	return map[byte][2]int{
		gravitytypes.EthAddressByValidatorKey[0]:              {0, 0},
		gravitytypes.ValidatorByEthAddressKey[0]:              {0, 0},
		gravitytypes.ValidatorAddressByOrchestratorAddress[0]: {0, 0},
		gravitytypes.LastEventBlockHeightByValidatorKey[0]:    {0, 0},
		gravitytypes.LastEventNonceByValidatorKey[0]:          {0, 0},
		gravitytypes.LastObservedEventNonceKey[0]:             {0, 0},
		gravitytypes.SequenceKeyPrefix[0]:                     {0, 0},
		gravitytypes.DenomToERC20Key[0]:                       {0, 0},
		gravitytypes.ERC20ToDenomKey[0]:                       {0, 0},
		gravitytypes.LastSlashedValsetNonce[0]:                {0, 0},
		gravitytypes.LatestValsetNonce[0]:                     {0, 0},
		gravitytypes.LastSlashedBatchBlock[0]:                 {0, 0},
		gravitytypes.LastUnBondingBlockHeight[0]:              {0, 0},
		gravitytypes.LastObservedEthereumBlockHeightKey[0]:    {0, 0},
		gravitytypes.LastObservedValsetKey[0]:                 {0, 0},
		gravitytypes.IbcSequenceHeightKey[0]:                  {0, 0},
		gravitytypes.ValsetRequestKey[0]:                      {0, 0},
		gravitytypes.OracleAttestationKey[0]:                  {0, 0},
		gravitytypes.BatchConfirmKey[0]:                       {0, 0},
		gravitytypes.ValsetConfirmKey[0]:                      {0, 0},
	}
}

func getBscKey() map[byte][2]int {
	return map[byte][2]int{
		crosschaintypes.OracleKey[0]:                          {2, 0},
		crosschaintypes.OracleAddressByExternalKey[0]:         {2, 0},
		crosschaintypes.OracleAddressByBridgerKey[0]:          {2, 0},
		crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
		crosschaintypes.OracleSetConfirmKey[0]:                {0, 0},
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
		crosschaintypes.PastExternalSignatureCheckpointKey[0]: {0, 0},
		crosschaintypes.BatchConfirmKey[0]:                    {0, 0},
		crosschaintypes.LastProposalBlockHeight[0]:            {0, 0},
	}
}

func getPolygonKey() map[byte][2]int {
	return map[byte][2]int{
		crosschaintypes.OracleKey[0]:                          {10, 0},
		crosschaintypes.OracleAddressByExternalKey[0]:         {10, 0},
		crosschaintypes.OracleAddressByBridgerKey[0]:          {10, 0},
		crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
		crosschaintypes.OracleSetConfirmKey[0]:                {0, 0},
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
		crosschaintypes.PastExternalSignatureCheckpointKey[0]: {0, 0},
		crosschaintypes.BatchConfirmKey[0]:                    {0, 0},
		crosschaintypes.LastProposalBlockHeight[0]:            {0, 0},
	}
}

func getTronKey() map[byte][2]int {
	return map[byte][2]int{
		crosschaintypes.OracleKey[0]:                          {10, 0},
		crosschaintypes.OracleAddressByExternalKey[0]:         {10, 0},
		crosschaintypes.OracleAddressByBridgerKey[0]:          {10, 0},
		crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
		crosschaintypes.OracleSetConfirmKey[0]:                {0, 0},
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
		crosschaintypes.PastExternalSignatureCheckpointKey[0]: {0, 0},
		crosschaintypes.BatchConfirmKey[0]:                    {0, 0},
		crosschaintypes.LastProposalBlockHeight[0]:            {0, 0},
	}
}

func getEthKey() map[byte][2]int {
	return map[byte][2]int{
		crosschaintypes.OracleKey[0]:                          {20, 0},
		crosschaintypes.OracleAddressByExternalKey[0]:         {20, 0},
		crosschaintypes.OracleAddressByBridgerKey[0]:          {20, 0},
		crosschaintypes.OracleSetRequestKey[0]:                {1, 0},
		crosschaintypes.OracleSetConfirmKey[0]:                {0, 0},
		crosschaintypes.OutgoingTxPoolKey[0]:                  {0, 0},
		crosschaintypes.OutgoingTxBatchKey[0]:                 {0, 0},
		crosschaintypes.OutgoingTxBatchBlockKey[0]:            {0, 0},
		crosschaintypes.LastEventNonceByOracleKey[0]:          {20, 0},
		crosschaintypes.LastObservedEventNonceKey[0]:          {1, 0},
		crosschaintypes.SequenceKeyPrefix[0]:                  {2, 0},
		crosschaintypes.DenomToTokenKey[0]:                    {10, 0},
		crosschaintypes.TokenToDenomKey[0]:                    {10, 0},
		crosschaintypes.LastSlashedOracleSetNonce[0]:          {1, 0},
		crosschaintypes.LatestOracleSetNonce[0]:               {1, 0},
		crosschaintypes.LastSlashedBatchBlock[0]:              {1, 0},
		crosschaintypes.LastObservedBlockHeightKey[0]:         {1, 0},
		crosschaintypes.LastObservedOracleSetKey[0]:           {1, 0},
		crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {20, 0},
		crosschaintypes.LastOracleSlashBlockHeight[0]:         {0, 0},
		crosschaintypes.ProposalOracleKey[0]:                  {1, 0},
		crosschaintypes.LastTotalPowerKey[0]:                  {1, 0},
		crosschaintypes.OracleAttestationKey[0]:               {1, 0},
		crosschaintypes.PastExternalSignatureCheckpointKey[0]: {0, 0},
		crosschaintypes.BatchConfirmKey[0]:                    {0, 0},
		crosschaintypes.LastProposalBlockHeight[0]:            {0, 0},
	}
}

func getAvalancheKey() map[byte][2]int {
	return map[byte][2]int{
		crosschaintypes.OracleKey[0]:                          {0, 0},
		crosschaintypes.OracleAddressByExternalKey[0]:         {0, 0},
		crosschaintypes.OracleAddressByBridgerKey[0]:          {0, 0},
		crosschaintypes.OracleSetRequestKey[0]:                {0, 0},
		crosschaintypes.OracleSetConfirmKey[0]:                {0, 0},
		crosschaintypes.OutgoingTxPoolKey[0]:                  {0, 0},
		crosschaintypes.OutgoingTxBatchKey[0]:                 {0, 0},
		crosschaintypes.OutgoingTxBatchBlockKey[0]:            {0, 0},
		crosschaintypes.LastEventNonceByOracleKey[0]:          {0, 0},
		crosschaintypes.LastObservedEventNonceKey[0]:          {1, 0},
		crosschaintypes.SequenceKeyPrefix[0]:                  {0, 0},
		crosschaintypes.DenomToTokenKey[0]:                    {0, 0},
		crosschaintypes.TokenToDenomKey[0]:                    {0, 0},
		crosschaintypes.LastSlashedOracleSetNonce[0]:          {1, 0},
		crosschaintypes.LatestOracleSetNonce[0]:               {1, 0},
		crosschaintypes.LastSlashedBatchBlock[0]:              {1, 0},
		crosschaintypes.LastObservedBlockHeightKey[0]:         {1, 0},
		crosschaintypes.LastObservedOracleSetKey[0]:           {1, 0},
		crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {0, 0},
		crosschaintypes.LastOracleSlashBlockHeight[0]:         {0, 0},
		crosschaintypes.ProposalOracleKey[0]:                  {1, 0},
		crosschaintypes.LastTotalPowerKey[0]:                  {1, 0},
		crosschaintypes.OracleAttestationKey[0]:               {0, 0},
		crosschaintypes.PastExternalSignatureCheckpointKey[0]: {0, 0},
		crosschaintypes.BatchConfirmKey[0]:                    {0, 0},
		crosschaintypes.LastProposalBlockHeight[0]:            {0, 0},
	}
}

func getErc20Key() map[byte][2]int {
	return map[byte][2]int{
		erc20types.KeyPrefixTokenPair[0]:        {12 + 5, 0},
		erc20types.KeyPrefixTokenPairByERC20[0]: {12 + 5, 0},
		erc20types.KeyPrefixTokenPairByDenom[0]: {12 + 5, 0},
		erc20types.KeyPrefixIBCTransfer[0]:      {0, 0},
		erc20types.KeyPrefixAliasDenom[0]:       {11 + 5, 0},
	}
}

func getMigrateKey() map[byte][2]int {
	return map[byte][2]int{
		migratetypes.KeyPrefixMigratedRecord[0]:        {1242, 0},
		migratetypes.KeyPrefixMigratedDirectionFrom[0]: {621, 0},
		migratetypes.KeyPrefixMigratedDirectionTo[0]:   {621, 0},
	}
}
