import { ethers } from "ethers";
import { txFee, getSigner, waitForTransaction } from "./common";

async function main() {
  const signer = await getSigner();

  const gasServiceContract = new ethers.Contract(
    "0x2d5d7d31F671F86C782533cc367F14109a082712",
    '[{"inputs":[{"internalType":"bytes32","name":"txHash","type":"bytes32"},{"internalType":"uint256","name":"logIndex","type":"uint256"},{"internalType":"address","name":"refundAddress","type":"address"}],"name":"addNativeGas","outputs":[],"stateMutability":"payable","type":"function"}]',
    signer
  );

  const txHash =
    "0x2b62e96d2047fe855d60d6e018be1de92a80b24e4924c9a1d24d0363c203fd06";
  const logIndex = 223;

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
