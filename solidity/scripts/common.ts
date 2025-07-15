// @ts-ignore
import hre from "hardhat";
import { Signer } from "ethers";
import { LedgerSigner } from "@ethersproject/hardware-wallets";
// import interchainTokenFactory from "@axelar-network/interchain-token-service/typescript/contracts/InterchainTokenFactory/InterchainTokenFactory.abi";
// import interchainTokenService from "@axelar-network/interchain-token-service/typescript/contracts/InterchainTokenService/InterchainTokenService.abi";
// import interchainToken from "@axelar-network/interchain-token-service/typescript/contracts/interchain-token/InterchainToken/InterchainToken.abi";

// ref: https://docs.axelar.dev/resources/contract-addresses/testnet/
// ref: https://testnet.interchain.axelar.dev/ethereum-sepolia/0xebCb46E14bCd0F8639A24b32fBC6Db6935F046Fe

const sourceChainNameMainnet = "bsc";

// export const interchainTokenServiceContractABI = interchainTokenService.abi;
// export const interchainTokenFactoryContractABI = interchainTokenFactory.abi;
// export const interchainTokenABI = interchainToken.abi;
export const interchainTokenServiceContractABI =
  '[{"inputs":[{"internalType":"address","name":"tokenAddress","type":"address"},{"internalType":"uint256","name":"gasValue","type":"uint256"}],"name":"registerTokenMetadata","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"tokenId","type":"bytes32"},{"internalType":"string","name":"destinationChain","type":"string"},{"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"bytes","name":"metadata","type":"bytes"},{"internalType":"uint256","name":"gasValue","type":"uint256"}],"name":"interchainTransfer","outputs":[],"stateMutability":"payable","type":"function"}]';
export const interchainTokenFactoryContractABI =
  '[{"inputs":[{"internalType":"bytes32","name":"salt","type":"bytes32"},{"internalType":"address","name":"tokenAddress","type":"address"},{"internalType":"enum ITokenManagerType.TokenManagerType","name":"tokenManagerType","type":"uint8"},{"internalType":"address","name":"operator","type":"address"}],"name":"registerCustomToken","outputs":[{"internalType":"bytes32","name":"tokenId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"salt","type":"bytes32"},{"internalType":"string","name":"destinationChain","type":"string"},{"internalType":"bytes","name":"destinationTokenAddress","type":"bytes"},{"internalType":"enum ITokenManagerType.TokenManagerType","name":"tokenManagerType","type":"uint8"},{"internalType":"bytes","name":"linkParams","type":"bytes"},{"internalType":"uint256","name":"gasValue","type":"uint256"}],"name":"linkToken","outputs":[{"internalType":"bytes32","name":"tokenId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"deployer","type":"address"},{"internalType":"bytes32","name":"salt","type":"bytes32"}],"name":"linkedTokenId","outputs":[{"internalType":"bytes32","name":"tokenId","type":"bytes32"}],"stateMutability":"view","type":"function"}]';
export const interchainTokenABI =
  '[{"inputs":[{"internalType":"bytes32","name":"salt","type":"bytes32"}],"name":"setItsSalt","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"its","type":"address"}],"name":"setInterchainTokenService","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"interchainTokenService","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"interchainTokenId","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"grantRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"hasRole","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]';

export const interchainTokenServiceContractAddress =
  "0xB5FB4BE02232B1bBA4dC8f81dc24C26980dE9e3C"; // mainnet,testnet

export const interchainTokenFactoryContractAddress =
  "0x83a93500d23Fbc3e82B410aD07A6a9F7A0670D66"; // mainnet,testnet

export const destinationChainTokenAddress =
  hre.network.name == sourceChainNameMainnet
    ? "0x075F23b9CdfCE2cC0cA466F4eE6cb4bD29d83bef" // ethereum mainnet pundiaifx
    : "0xebCb46E14bCd0F8639A24b32fBC6Db6935F046Fe"; // ethereum sepolia pundiaifx

export const destinationChainName =
  hre.network.name == sourceChainNameMainnet ? "ethereum" : "ethereum-sepolia";

export const sourceChainTokenAddress = process.env.SOURCE_CHAIN_TOKEN_ADDRESS;

export const interchainTokenManagerAddress = process.env.TOKEN_MANAGER_ADDRESS;

// console.log("nwe salt", "0x" + crypto.randomBytes(32).toString("hex"))
export const salt =
  hre.network.name == sourceChainNameMainnet
    ? "0xfb924245be99a0fa4b9699ac455ef575c1d73fb76d624a125e1783956b575a87"
    : "0x009bcad5630a54a1f72fdf8230c151b787c658126c876afee74c0207f05dd028";

export const tokenManagerTypeMintBurn = 4;
export const tokenManagerTypeLockUnLock = 2;

export const txFee = "0.02";

console.log(
  "interchainTokenServiceContractAddress:",
  interchainTokenServiceContractAddress
);
console.log(
  "interchainTokenFactoryContractAddress:",
  interchainTokenFactoryContractAddress
);
console.log("destinationChainTokenAddress:", destinationChainTokenAddress);
console.log("destinationChainName:", destinationChainName);
console.log("sourceChainTokenAddress:", sourceChainTokenAddress);
console.log("interchainTokenManagerAddress:", interchainTokenManagerAddress);
console.log("salt:", salt);
console.log("tokenManagerTypeMintBurn:", tokenManagerTypeMintBurn);
console.log("tokenManagerTypeLockUnLock:", tokenManagerTypeLockUnLock);
console.log("txFee:", txFee);

export async function getSigner(): Promise<Signer> {
  const nodeUrl = hre.network.config.url;
  const provider = new hre.ethers.JsonRpcProvider(nodeUrl);
  if (process.env.RAW_PRIVATE_KEY !== undefined) {
    const signer = new hre.ethers.Wallet(process.env.RAW_PRIVATE_KEY, provider);
    console.log("Using raw private key signer:", await signer.getAddress());
    return signer;
  }
  const derivationPath = process.env.DERIVATION_PATH || "m/44'/60'/0'/0/0";
  const signer = new LedgerSigner(provider, "default", derivationPath);
  console.log("Using Ledger signer:", await signer.getAddress());
  // @ts-ignore
  return signer;
}

export async function waitForTransaction(tx: any): Promise<any> {
  const receipt = await tx.wait();
  console.log(
    "Transaction mined in block:",
    receipt.blockNumber,
    "status:",
    receipt.status
  );
  return receipt;
}

export function requireSourceChainTokenAddress() {
  if (!sourceChainTokenAddress) {
    throw new Error(
      "SOURCE_CHAIN_TOKEN_ADDRESS environment variable is required"
    );
  }
}

export function requireInterchainTokenManagerAddress() {
  if (!interchainTokenManagerAddress) {
    throw new Error("TOKEN_MANAGER_ADDRESS environment variable is required");
  }
}
