import { ethers } from "ethers";
import {
  destinationChainName,
  getSigner,
  interchainTokenFactoryContractABI,
  interchainTokenFactoryContractAddress,
  interchainTokenServiceContractABI,
  interchainTokenServiceContractAddress,
  salt,
  txFee,
  waitForTransaction,
} from "./common";

async function main() {
  let destChainName = process.env.DESTINATION_CHAIN_NAME;
  if (!destChainName) {
    destChainName = destinationChainName;
  }
  console.log("destChainName:", destChainName);
  const signer = await getSigner();

  const interchainTokenFactoryContract = new ethers.Contract(
    interchainTokenFactoryContractAddress,
    interchainTokenFactoryContractABI,
    signer
  );
  let signerAddr = await signer.getAddress();
  const tokenId = await interchainTokenFactoryContract.linkedTokenId(
    signerAddr,
    salt
  );

  const interchainTokenService = new ethers.Contract(
    interchainTokenServiceContractAddress,
    interchainTokenServiceContractABI,
    signer
  );
  console.log("tokenId:", tokenId);
  const interchainTransferTx = await interchainTokenService.interchainTransfer(
    tokenId,
    destChainName,
    signerAddr,
    ethers.parseEther("0.12"),
    "0x",
    ethers.parseEther(txFee),
    {
      value: ethers.parseEther(txFee),
    }
  );
  console.log("interchainTransfer tx:", interchainTransferTx.hash);

  await waitForTransaction(interchainTransferTx);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
