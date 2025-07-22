import { ethers } from "ethers";
import {
  destinationChainTokenAddress,
  destinationChainName,
  interchainTokenFactoryContractABI,
  interchainTokenFactoryContractAddress,
  tokenManagerTypeLockUnLock,
  salt,
  txFee,
  getSigner,
  waitForTransaction,
} from "./common";

async function main() {
  const signer = await getSigner();

  const interchainTokenFactoryContract = new ethers.Contract(
    interchainTokenFactoryContractAddress,
    interchainTokenFactoryContractABI,
    signer
  );

  let signerAddr = await signer.getAddress();
  console.log("linkToken tx params:", {
    destinationChainTokenAddress,
    destinationChainName,
    tokenManagerTypeLockUnLock,
    signerAddr,
    salt,
    txFee,
    value: ethers.parseEther(txFee),
  });
  const linkToken = await interchainTokenFactoryContract.linkToken(
    salt,
    destinationChainName,
    destinationChainTokenAddress,
    tokenManagerTypeLockUnLock,
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
