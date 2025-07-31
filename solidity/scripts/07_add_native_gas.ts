import { ethers } from "ethers";
import {
  txFee,
  getSigner,
  waitForTransaction,
  gasServiceContractAddress,
  gasServiceABI,
} from "./common";

async function main() {
  const txHash = process.env.Add_GAS_TX_HASH;
  const logIndex = process.env.Add_GAS_LOG_INDEX;
  if (!txHash || !logIndex) {
    throw new Error(
      "Please set the environment variables Add_GAS_TX_HASH and Add_GAS_LOG_INDEX"
    );
  }

  const signer = await getSigner();

  const gasServiceContract = new ethers.Contract(
    gasServiceContractAddress,
    gasServiceABI,
    signer
  );

  let signerAddr = await signer.getAddress();
  console.log("addNativeGas tx params:", {
    addNativeGas: ethers.parseEther(txFee),
    txHash,
    logIndex,
    refundAddress: signerAddr,
  });
  const addNativeGas = await gasServiceContract.addNativeGas(
    txHash,
    logIndex,
    signerAddr,
    { value: ethers.parseEther(txFee) }
  );
  console.log("addNativeGas tx:", addNativeGas.hash);

  await waitForTransaction(addNativeGas);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
