import { ethers } from "ethers";
import {
  sourceChainTokenAddress,
  interchainTokenFactoryContractABI,
  interchainTokenFactoryContractAddress,
  tokenManagerTypeMintBurn,
  salt,
  txFee,
  getSigner,
  waitForTransaction,
  interchainTokenServiceContractAddress,
  requireSourceChainTokenAddress,
} from "./common";

async function main() {
  requireSourceChainTokenAddress();

  const signer = await getSigner();

  const interchainTokenFactoryContract = new ethers.Contract(
    interchainTokenFactoryContractAddress,
    interchainTokenFactoryContractABI,
    signer
  );

  let sigerAddr = await signer.getAddress();
  console.log("registerCustomToken tx params:", {
    sourceChainTokenAddress,
    tokenManagerTypeMintBurn,
    sigerAddr,
    salt,
    txFee,
    value: ethers.parseEther(txFee),
  });
  const registerCustomTokenTx =
    await interchainTokenFactoryContract.registerCustomToken(
      salt,
      sourceChainTokenAddress,
      tokenManagerTypeMintBurn,
      sigerAddr,
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
