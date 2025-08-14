import { ethers } from "ethers";
import {
  getSigner,
  interchainTokenFactoryContractABI,
  interchainTokenFactoryContractAddress,
  interchainTokenServiceContractAddress,
  requireSourceChainTokenAddress,
  salt,
  sourceChainTokenAddress,
  tokenManagerType,
  txFee,
  waitForTransaction,
} from "./common";

async function main() {
  requireSourceChainTokenAddress();
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

  console.log("registerCustomToken tx params:", {
    sourceChainTokenAddress,
    tokenManagerType,
    signerAddr,
    salt,
    txFee,
    value: ethers.parseEther(txFee),
  });
  const registerCustomTokenTx =
    await interchainTokenFactoryContract.registerCustomToken(
      salt,
      sourceChainTokenAddress,
      tokenManagerType,
      signerAddr,
      { value: ethers.parseEther(txFee) }
    );
  console.log("registerCustomToken tx:", registerCustomTokenTx.hash);

  const receipt = await waitForTransaction(registerCustomTokenTx);
  receipt.logs.forEach(
    (log: { address: string; topics: any[] }, index: number) => {
      if (
        index === 0 &&
        log.address === interchainTokenServiceContractAddress
      ) {
        console.log(
          "source chain new tokenId:",
          log.topics[1],
          "salt:",
          log.topics[3]
        );
      }
      if (index === 1) {
        console.log("source chain new tokenManagerAddress:", log.address);
      }
    }
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
