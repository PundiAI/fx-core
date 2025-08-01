// @ts-ignore
import hre from "hardhat";
import { Signer } from "ethers";

// ref: https://docs.axelar.dev/resources/contract-addresses/testnet/
// ref: https://testnet.interchain.axelar.dev/ethereum-sepolia/0xebCb46E14bCd0F8639A24b32fBC6Db6935F046Fe

const sourceChainName = process.env.SOURCE_CHAIN_NAME || "bsc";

export const interchainTokenServiceContractABI =
  '[{"inputs":[{"internalType":"address","name":"tokenAddress","type":"address"},{"internalType":"uint256","name":"gasValue","type":"uint256"}],"name":"registerTokenMetadata","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"tokenId","type":"bytes32"},{"internalType":"string","name":"destinationChain","type":"string"},{"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"bytes","name":"metadata","type":"bytes"},{"internalType":"uint256","name":"gasValue","type":"uint256"}],"name":"interchainTransfer","outputs":[],"stateMutability":"payable","type":"function"}]';
export const interchainTokenFactoryContractABI =
  '[{"inputs":[{"internalType":"bytes32","name":"salt","type":"bytes32"},{"internalType":"address","name":"tokenAddress","type":"address"},{"internalType":"enum ITokenManagerType.TokenManagerType","name":"tokenManagerType","type":"uint8"},{"internalType":"address","name":"operator","type":"address"}],"name":"registerCustomToken","outputs":[{"internalType":"bytes32","name":"tokenId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"salt","type":"bytes32"},{"internalType":"string","name":"destinationChain","type":"string"},{"internalType":"bytes","name":"destinationTokenAddress","type":"bytes"},{"internalType":"enum ITokenManagerType.TokenManagerType","name":"tokenManagerType","type":"uint8"},{"internalType":"bytes","name":"linkParams","type":"bytes"},{"internalType":"uint256","name":"gasValue","type":"uint256"}],"name":"linkToken","outputs":[{"internalType":"bytes32","name":"tokenId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"deployer","type":"address"},{"internalType":"bytes32","name":"salt","type":"bytes32"}],"name":"linkedTokenId","outputs":[{"internalType":"bytes32","name":"tokenId","type":"bytes32"}],"stateMutability":"view","type":"function"}]';
export const interchainTokenABI =
  '[{"inputs":[{"internalType":"bytes32","name":"salt","type":"bytes32"}],"name":"setItsSalt","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"its","type":"address"}],"name":"setInterchainTokenService","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"interchainTokenService","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"interchainTokenId","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"grantRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"hasRole","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"grantRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"revokeRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]';
export const gasServiceABI =
  '[{"inputs":[{"internalType":"bytes32","name":"txHash","type":"bytes32"},{"internalType":"uint256","name":"logIndex","type":"uint256"},{"internalType":"address","name":"refundAddress","type":"address"}],"name":"addNativeGas","outputs":[],"stateMutability":"payable","type":"function"}]';

export const interchainTokenManagerABI =
  '[{"inputs":[],"name":"getImplementationTypeAndTokenAddress","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"flowLimit","outputs":[{"internalType":"uint256","name":"flowLimit_","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"flowLimit_","type":"uint256"}],"name":"setFlowLimit","outputs":[],"stateMutability":"nonpayable","type":"function"}]';

export const interchainTokenServiceContractAddress =
  "0xB5FB4BE02232B1bBA4dC8f81dc24C26980dE9e3C"; // mainnet,testnet

export const interchainTokenFactoryContractAddress =
  "0x83a93500d23Fbc3e82B410aD07A6a9F7A0670D66"; // mainnet,testnet

export const gasServiceContractAddress =
  "0x2d5d7d31F671F86C782533cc367F14109a082712";

export const destinationChainTokenAddress =
  hre.network.name == sourceChainName
    ? "0x075F23b9CdfCE2cC0cA466F4eE6cb4bD29d83bef" // ethereum mainnet pundiaifx
    : "0xebCb46E14bCd0F8639A24b32fBC6Db6935F046Fe"; // ethereum sepolia pundiaifx

export let destinationChainName = "";

if (sourceChainName != "ethereum") {
  destinationChainName =
    hre.network.name == sourceChainName ? "Ethereum" : "ethereum-sepolia";
}

export const sourceChainTokenAddress: string =
  process.env.SOURCE_CHAIN_TOKEN_ADDRESS || "";

export const interchainTokenManagerAddress: string =
  process.env.TOKEN_MANAGER_ADDRESS || "";

// import crypto from "crypto";
// console.log("nwe salt", "0x" + crypto.randomBytes(32).toString("hex"))
export const salt =
  hre.network.name == sourceChainName
    ? "0xc144c0dcfd41bcf51f2ee6d6cea553d8bf07f31c330cf5b698b0dee6bdb41308"
    : "0x2b1f7645a5ea54f9c9281d2a71038e31b50422570195c08ce3124929c1a709ef";

export const tokenManagerTypeMintBurn = 4;
export const tokenManagerTypeLockUnLock = 2;

export const txFee = process.env.TX_FEE ? process.env.TX_FEE : "0.002";

console.log(
  "interchainTokenServiceContractAddress:",
  interchainTokenServiceContractAddress
);
console.log(
  "interchainTokenFactoryContractAddress:",
  interchainTokenFactoryContractAddress
);

console.log({
  sourceChainName,
  network: hre.network.name,
  destinationChainTokenAddress,
  destinationChainName,
  sourceChainTokenAddress,
  interchainTokenManagerAddress,
  gasServiceContractAddress,
  salt,
  tokenManagerTypeMintBurn,
  tokenManagerTypeLockUnLock,
  txFee,
});

export async function getSigner(): Promise<Signer> {
  if (hre.network.config.accounts !== undefined) {
    const signer = (await hre.ethers.getSigners())[0];
    console.log("Using raw private key signer:", await signer.getAddress());
    return signer;
  }
  if (
    !hre.network.config.ledgerAccounts ||
    hre.network.config.ledgerAccounts.length === 0
  ) {
    throw new Error("No ledger accounts configured in hardhat config");
  }
  console.log("Using Ledger signer:", hre.network.config.ledgerAccounts[0]);
  return await hre.ethers.getSigner(hre.network.config.ledgerAccounts[0]);
}

export async function waitForTransaction(tx: any): Promise<any> {
  const receipt = await tx.wait();
  console.log(
    "Transaction mined in block:",
    receipt.blockNumber,
    "status:",
    receipt.status
  );
  if (receipt.status !== 1) {
    throw new Error(`Transaction failed with status: ${receipt.status}`);
  }
  return receipt;
}

export function requireSourceChainTokenAddress() {
  if (sourceChainTokenAddress === "") {
    throw new Error(
      "SOURCE_CHAIN_TOKEN_ADDRESS environment variable is required"
    );
  }
}

export function requireInterchainTokenManagerAddress() {
  if (interchainTokenManagerAddress === "") {
    throw new Error("TOKEN_MANAGER_ADDRESS environment variable is required");
  }
}
