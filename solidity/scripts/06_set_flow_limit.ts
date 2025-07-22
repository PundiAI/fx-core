import { ethers, keccak256, toUtf8Bytes } from "ethers";
import {
  sourceChainTokenAddress,
  interchainTokenManagerAddress,
  interchainTokenManagerABI,
  getSigner,
  waitForTransaction,
  requireSourceChainTokenAddress,
  requireInterchainTokenManagerAddress,
} from "./common";

async function main() {
  requireSourceChainTokenAddress();
  requireInterchainTokenManagerAddress();

  const signer = await getSigner();

  const tokenManger = new ethers.Contract(
    interchainTokenManagerAddress,
    interchainTokenManagerABI,
    signer
  );

  const [managerType, tokenAddr] =
    await tokenManger.getImplementationTypeAndTokenAddress();
  console.log(`ManagerType :${managerType}, TokenAddr: ${tokenAddr}`);

  if (tokenAddr !== sourceChainTokenAddress) {
    throw new Error(
      `tokenAddr !== sourceChainTokenAddress, tokenAddr: ${tokenAddr}, sourceChainTokenAddress: ${sourceChainTokenAddress}`
    );
  }

  const flowLimitValue = 1;

  const flowLimit = await tokenManger.flowLimit();
  console.log(`current flowLimit: ${flowLimit}`);
  if (flowLimitValue == flowLimit) {
    console.log(
      `flowLimitValue == flowLimit, no need to set, flowLimitValue: ${flowLimitValue}`
    );
    return;
  }
  console.log(`set flowLimit to: ${flowLimitValue}`);
  const setFlowLimitTx = await tokenManger.setFlowLimit(flowLimitValue);
  await waitForTransaction(setFlowLimitTx);

  const newFlowLimit = await tokenManger.flowLimit();
  console.log(`new flowLimit: ${newFlowLimit}`);
  if (newFlowLimit != flowLimitValue) {
    throw new Error(
      `newFlowLimit !== flowLimitValue, newFlowLimit: ${newFlowLimit}, flowLimitValue: ${flowLimitValue}`
    );
  }
  console.log(`set flowLimit success`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
