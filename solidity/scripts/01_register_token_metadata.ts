import { ethers } from "ethers";
import {
  sourceChainTokenAddress,
  interchainTokenServiceContractABI,
  interchainTokenServiceContractAddress,
  txFee,
  getSigner,
  waitForTransaction,
  requireSourceChainTokenAddress,
} from "./common";

async function main() {
  requireSourceChainTokenAddress();

  const signer = await getSigner();

  const interchainTokenService = new ethers.Contract(
    interchainTokenServiceContractAddress,
    interchainTokenServiceContractABI,
    signer
  );

  console.log("registerTokenMetadata tx params:", {
    sourceChainTokenAddress,
    txFee,
    value: ethers.parseEther(txFee),
  });
  const registerTokenMetadata =
    await interchainTokenService.registerTokenMetadata(
      sourceChainTokenAddress,
      ethers.parseEther(txFee),
      { value: ethers.parseEther(txFee) }
    );
  console.log("registerTokenMetadata tx:", registerTokenMetadata.hash);
  console.log(
    "axelascan: https://axelarscan.io/gmp/" + registerTokenMetadata.hash
  );

  await waitForTransaction(registerTokenMetadata);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
