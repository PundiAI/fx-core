// @ts-ignore
import { ethers, keccak256, toUtf8Bytes } from "ethers";
import {
  sourceChainTokenAddress,
  interchainTokenABI,
  interchainTokenManagerAddress,
  interchainTokenServiceContractAddress,
  salt,
  getSigner,
  waitForTransaction,
  requireSourceChainTokenAddress,
  requireInterchainTokenManagerAddress,
} from "./common";

async function main() {
  requireSourceChainTokenAddress();
  requireInterchainTokenManagerAddress();

  const signer = await getSigner();

  const pundiaifxContract = new ethers.Contract(
    sourceChainTokenAddress,
    interchainTokenABI,
    signer
  );

  const setItsSaltTx = await pundiaifxContract.setItsSalt(salt);
  console.log("setItsSalt tx:", setItsSaltTx.hash);
  await waitForTransaction(setItsSaltTx);

  const setItsTx = await pundiaifxContract.setInterchainTokenService(
    interchainTokenServiceContractAddress
  );
  console.log("setInterchainTokenService tx:", setItsTx.hash);
  await waitForTransaction(setItsTx);

  const roleBytes = keccak256(toUtf8Bytes("ADMIN_ROLE"));
  const hasRole = await pundiaifxContract.hasRole(
    roleBytes,
    interchainTokenManagerAddress
  );
  if (hasRole) {
    console.log("interchainTokenManagerAddress has ADMIN_ROLE");
    return;
  }
  console.log("interchainTokenManagerAddress does not have ADMIN_ROLE");

  const grantRoleTx = await pundiaifxContract.grantRole(
    roleBytes,
    interchainTokenManagerAddress
  );
  console.log("grantRole tx:", grantRoleTx.hash);

  await waitForTransaction(grantRoleTx);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
