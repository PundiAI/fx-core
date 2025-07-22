# Register Custom Token to Axelar ITS

## Install Dependencies

```bash
yarn install
```

## Compile Contracts

```bash
yarn compile
```

## Deploy ERC20 Token

```bash
export RAW_PRIVATE_KEY="<YOUR PRIVATE KEY>"
npx hardhat ignition deploy ignition/modules/interchain_token.ts --network bsc
```

## Register Token Metadata with the ITS Contract

```bash
export SOURCE_CHAIN_TOKEN_ADDRESS="<BSC PUNDIAI TOKEN ADDRESS>"
```

```bash
npx hardhat run scripts/01_register_token_metadata.ts --network bsc
```

## Register Custom Token with the Interchain Token Factory

```bash
npx hardhat run scripts/02_register_custom_token.ts --network bsc
```

## Link Custom Token with the Interchain Token Factory

```bash
export TOKEN_MANAGER_ADDRESS="<BSC PUNDIAI TOKEN MANAGER ADDRESS>"
npx hardhat run scripts/03_link_token.ts --network bsc
```

## Assign the Minter Role to TokenManger
```bash
npx hardhat run scripts/04_interchain_token.ts --network bsc
```

## Set Flow Limit
```bash
npx hardhat run scripts/06_set_flow_limit.ts --network bsc
```
