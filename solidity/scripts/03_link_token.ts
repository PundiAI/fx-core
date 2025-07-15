// @ts-ignore
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

  const linkToken = await interchainTokenFactoryContract.linkToken(
    salt,
    destinationChainName,
    destinationChainTokenAddress,
    tokenManagerTypeLockUnLock,
    await signer.getAddress(),
    ethers.parseEther(txFee),
    { value: ethers.parseEther(txFee) }
  );
  console.log("linkToken tx:", linkToken.hash);

  await waitForTransaction(linkToken);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
