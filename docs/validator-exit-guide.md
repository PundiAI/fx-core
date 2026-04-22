# Validator exit guide

This guide explains how to exit as a validator on the **Pundi AIFX** chain by unbonding your self-delegation. When a validator's self-delegation falls below the minimum self-delegation threshold, the validator will automatically go offline. **Importantly, the validator will not be penalized** in this process.

---

## Method 1: Via command line

### Step 1: Query self-delegation

Before unbonding, check your current self-delegation. Use your **account address** (the same key you used to self-delegate) and your **validator operator address** (`fxvaloper...`):

```bash
fxcored query staking delegation [your-account-address] [your-validator-valoper] --node https://fx-json.functionx.io:26657
```

To list **all** delegations to your validator (including self-delegation):

```bash
fxcored query staking delegations-to [your-validator-valoper] --node https://fx-json.functionx.io:26657
```

### Step 2: Query validator information (optional)

To view general validator details:

```bash
fxcored query staking validator [your-validator-valoper] --node https://fx-json.functionx.io:26657
```

### Step 3: Unbond (cancel self-delegation)

Execute the unbond transaction:

```bash
fxcored tx staking unbond [validator-valoper] [amount] --from mywallet --gas auto --gas-adjustment 1.5 --gas-prices 5000000000apundiai --chain-id fxcore --node https://fx-json.functionx.io:26657
```

Replace:

- `[validator-valoper]` — Your validator operator address (`fxvaloper...`)
- `[amount]` — Amount to unbond (e.g. `1000000000000000000apundiai`)
- `mywallet` — Your wallet key name (list keys with `fxcored keys list`)

---

## Method 2: Via block explorer

1. Open the block explorer: [https://pundiscan.io/pundiaifx/validators](https://pundiscan.io/pundiaifx/validators)
2. Connect your wallet
3. Click the **Delegate** button
4. Select the **Unbonding** tab
5. View your delegation list and select the delegation you wish to unbond
6. Click the **Undelegate** button
7. Enter the amount to undelegate
8. Click **Submit** and wait for the transaction to complete

---

## Shutting down the node

After your validator has left the active set (or you no longer intend to run the binary):

1. Stop process supervision (systemd, cosmovisor, Kubernetes, etc.).
2. Stop `fxcored` cleanly so it does not keep trying to sign.
3. Keep a secure backup of keys and any on-disk data you need before decommissioning hardware.

---

## Important notes

- **No penalty**: When your self-delegation drops below the minimum self-delegation threshold, the validator will automatically go offline. You will **not** be slashed or penalized.
- **Unbonding period**: After unbonding, tokens enter an unbonding period before they become fully transferable. Check the chain parameters for the exact unbonding duration.
- **Minimum threshold**: Ensure you understand the minimum self-delegation requirement for validators on the network before proceeding.
