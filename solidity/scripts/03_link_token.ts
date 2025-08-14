import { ethers } from "ethers";
import {
  destinationChainTokenAddress,
  destinationChainName,
  interchainTokenFactoryContractABI,
  interchainTokenFactoryContractAddress,
  tokenManagerType,
  salt,
  txFee,
  getSigner,
  waitForTransaction,
} from "./common";

async function main() {
  if (destinationChainName === "") {
    throw new Error("DESTINATION_CHAIN_NAME environment variable is required");
  }
  if (destinationChainTokenAddress === "") {
    throw new Error(
      "DESTINATION_CHAIN_TOKEN_ADDRESS environment variable is required"
    );
  }
  if (tokenManagerType === "") {
    throw new Error("TOKEN_MANAGER_TYPE environment variable is required");
  }
  const signer = await getSigner();
  let signerAddr = await signer.getAddress();

  const interchainTokenFactoryContract = new ethers.Contract(
    interchainTokenFactoryContractAddress,
    interchainTokenFactoryContractABI,
    signer
  );

  console.log("linkToken tx params:", {
    destinationChainTokenAddress,
    destinationChainName,
    tokenManagerType,
    signerAddr,
    salt,
    txFee,
    value: ethers.parseEther(txFee),
  });
  const linkToken = await interchainTokenFactoryContract.linkToken(
    salt,
    destinationChainName,
    destinationChainTokenAddress,
    tokenManagerType,
    signerAddr,
    ethers.parseEther(txFee),
    { value: ethers.parseEther(txFee) }
  );
  console.log("linkToken tx:", linkToken.hash);
  console.log("axelascan: https://axelarscan.io/gmp/" + linkToken.hash);

  await waitForTransaction(linkToken);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
