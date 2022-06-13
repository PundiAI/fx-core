<!--
order: 2
-->

# State

## Keys

* OracleKey                           `0x12  + oracle.Bytes()`                                                            -> `k.cdc.MustMarshal(&Oracle)`
* OracleAddressByExternalKey          `0x13  + []byte(externalAddress)`                                                   -> `oracle.Bytes()`
* OracleAddressByBridgerKey           `0x14  + oracle.Bytes()`                                                            -> `[]byte(externalAddress)`
* OracleSetRequestKey                 `0x15  + sdk.Uint64ToBigEndian(nonce)`                                              -> `k.cdc.MustMarshal(&OracleSet)`
* OracleSetConfirmKey                 `0x16  + sdk.Uint64ToBigEndian(nonce) + validator.Bytes()`                          -> `k.cdc.MustMarshal(&MsgOracleSetConfirm)`
* OracleAttestationKey                `0x17  + sdk.Uint64ToBigEndian(nonce) + claimHash`                                  -> `k.cdc.MustMarshal(&Attestation)`
* OutgoingTxPoolKey                   `0x18  + sdk.Uint64ToBigEndian(outgoingTransferTxId)`                               -> `k.cdc.MustMarshal(&OutgoingTransferTx)`
* ? SecondIndexOutgoingTxFeeKey         `0x19  + []byte(tokenContract) + fee.Amount.BigInt().FillBytes(amount)`             -> `k.cdc.MustMarshal(&idSet)`
* OutgoingTxBatchKey                  `0x20  + []byte(tokenContract) + sdk.Uint64ToBigEndian(nonce)`                      -> `k.cdc.MustMarshal(&OutgoingTxBatch)`
* OutgoingTxBatchBlockKey             `0x21  + sdk.Uint64ToBigEndian(blockHeight)`                                        -> `k.cdc.MustMarshal(&OutgoingTxBatch)`
* BatchConfirmKey                     `0x22  + []byte(tokenContract) + sdk.Uint64ToBigEndian(nonce) + oracle.Bytes()`     -> `k.cdc.MustMarshal(&MsgConfirmBatch)`
* LastEventNonceByValidatorKey        `0x23  + validator.Bytes()`                                                         -> `sdk.Uint64ToBigEndian(nonce)`
* LastObservedEventNonceKey           `0x24`                                                                              -> `sdk.Uint64ToBigEndian(nonce)`
* KeyLastTxPoolID                     `0x25  + []byte("lastTxPoolId")`                                                    -> `sdk.Uint64ToBigEndian(id)`
* KeyLastOutgoingBatchID              `0x25  + []byte("lastBatchId")`                                                     -> `sdk.Uint64ToBigEndian(id)`
* DenomToTokenKey                     `0x26  + []byte(denom)`                                                             -> `[]byte(tokenContract)` 
* TokenToDenomKey                     `0x27 + []byte(tokenContract)`                                                     -> `[]byte(denom)` 
* LastSlashedOracleSetNonce           `0x28 `                                                                            -> `sdk.Uint64ToBigEndian(nonce)` 
* LatestOracleSetNonce                `0x29 `                                                                            -> `sdk.Uint64ToBigEndian(nonce)`
* LastSlashedBatchBlock               `0x30 `                                                                            -> `sdk.Uint64ToBigEndian(blockHeight)`
* LastProposalBlockHeight             `0x31 `                                                                            -> `sdk.Uint64ToBigEndian(blockHeight)`
* LastObservedBlockHeightKey          `0x32 `                                                                            -> `k.cdc.MustMarshal(lastObservedEthereumBlockHeight)`
* LastObservedOracleSetKey            `0x33 `                                                                            -> `k.cdc.MustMarshal(&OracleSet)`
* LastEventBlockHeightByValidatorKey  `0x35 + oracle.Bytes()`                                                            -> `sdk.Uint64ToBigEndian(blockHeight)`                
* PastExternalSignatureCheckpointKey  `0x36 + checkpoint`
* LastOracleSlashBlockHeight          `0x37 `                                                                            -> `sdk.Uint64ToBigEndian(blockHeight)`
* ProposalOraclesKey                  `0x38 `                                                                            -> `k.cdc.MustMarshal(&ChainOracle)`
* LastTotalPowerKey                   `0x39 `                                                                            -> `power`