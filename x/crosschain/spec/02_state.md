# State

## Keys

* OracleKey                           `0x12  + oracleAddr.Bytes()`                                                            -> `MustMarshal(&Oracle)`
* OracleAddressByExternalKey          `0x13  + []byte(externalAddress)`                                                       -> `oracleAddr.Bytes()`
* OracleAddressByBridgerKey           `0x14  + oracleAddr.Bytes()`                                                            -> `[]byte(externalAddress)`
* OracleSetRequestKey                 `0x15  + sdk.Uint64ToBigEndian(nonce)`                                                  -> `MustMarshal(&OracleSet)`
* OracleSetConfirmKey                 `0x16  + sdk.Uint64ToBigEndian(nonce) + oracleAddr.Bytes()`                             -> `MustMarshal(&MsgOracleSetConfirm)`
* OracleAttestationKey                `0x17  + sdk.Uint64ToBigEndian(eventNonce) + claimHash`                                 -> `MustMarshal(&Attestation)`
* OutgoingTxPoolKey                   `0x18  + []byte(tokenContract) + fee.amount + sdk.Uint64ToBigEndian(id)`                -> `MustMarshal(&OutgoingTransferTx)`
* OutgoingTxBatchKey                  `0x20  + []byte(tokenContract) + sdk.Uint64ToBigEndian(batchNonce)`                     -> `MustMarshal(&OutgoingTxBatch)`
* OutgoingTxBatchBlockKey             `0x21  + sdk.Uint64ToBigEndian(blockHeight)`                                            -> `MustMarshal(&OutgoingTxBatch)`
* BatchConfirmKey                     `0x22  + []byte(tokenContract) + sdk.Uint64ToBigEndian(batchNonce) + oracleAddr.Bytes()`-> `MustMarshal(&MsgConfirmBatch)`
* LastEventNonceByOracleKey           `0x23  + oracleAddr.Bytes()`                                                            -> `sdk.Uint64ToBigEndian(eventNonce)`
* LastObservedEventNonceKey           `0x24`                                                                                  -> `sdk.Uint64ToBigEndian(eventNonce)`
* KeyLastTxPoolID                     `0x25  + []byte("lastTxPoolId")`                                                        -> `sdk.Uint64ToBigEndian(id)`
* KeyLastOutgoingBatchID              `0x25  + []byte("lastBatchId")`                                                         -> `sdk.Uint64ToBigEndian(id)`
* DenomToTokenKey                     `0x26  + []byte(denom)`                                                                 -> `[]byte(tokenContract)` 
* TokenToDenomKey                     `0x27  + []byte(tokenContract)`                                                         -> `[]byte(denom)` 
* LastSlashedOracleSetNonce           `0x28 `                                                                                 -> `sdk.Uint64ToBigEndian(nonce)` 
* LatestOracleSetNonce                `0x29 `                                                                                 -> `sdk.Uint64ToBigEndian(nonce)`
* LastSlashedBatchBlock               `0x30 `                                                                                 -> `sdk.Uint64ToBigEndian(blockHeight)`
* LastObservedBlockHeightKey          `0x32 `                                                                                 -> `MustMarshal(lastObservedEthereumBlockHeight)`
* LastObservedOracleSetKey            `0x33 `                                                                                 -> `MustMarshal(&OracleSet)`
* LastEventBlockHeightByOracleKey     `0x35 + oracleAddr.Bytes()`                                                             -> `sdk.Uint64ToBigEndian(blockHeight)`                
* PastExternalSignatureCheckpointKey  `0x36 + checkpoint`
* LastOracleSlashBlockHeight          `0x37 `                                                                                 -> `sdk.Uint64ToBigEndian(blockHeight)`
* ProposalOracleKey                   `0x38 `                                                                                 -> `MustMarshal(&ChainOracle)`
* LastTotalPowerKey                   `0x39 `                                                                                 -> `power`