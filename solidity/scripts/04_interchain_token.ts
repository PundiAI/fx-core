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

  /*
  console.log("setItsSalt tx params:", {
    sourceChainTokenAddress,
    interchainTokenServiceContractAddress,
    interchainTokenManagerAddress,
    salt,
  });
  const setItsSaltTx = await pundiaifxContract.setItsSalt(salt);
  console.log("setItsSalt tx:", setItsSaltTx.hash);
  await waitForTransaction(setItsSaltTx);

  console.log("setInterchainTokenService tx params:", {
    sourceChainTokenAddress,
    interchainTokenServiceContractAddress,
    interchainTokenManagerAddress,
  });
  const setItsTx = await pundiaifxContract.setInterchainTokenService(
    interchainTokenServiceContractAddress
  );
  console.log("setInterchainTokenService tx:", setItsTx.hash);
  await waitForTransaction(setItsTx);
  */

  const roleBytes = keccak256(toUtf8Bytes("ADMIN_ROLE"));
  const revokeAddress = process.env.REVOKE_ADDRESS;
  if (revokeAddress) {
    const hasRole = await pundiaifxContract.hasRole(roleBytes, revokeAddress);
    if (hasRole) {
      console.log("revokeRole tx params:", {
        roleBytes,
        revokeAddress,
      });
      const revokeRole = await pundiaifxContract.revokeRole(
        roleBytes,
        revokeAddress
      );
      await waitForTransaction(revokeRole);
    }
    return;
  }

  const hasRole = await pundiaifxContract.hasRole(
    roleBytes,
    interchainTokenManagerAddress
  );
  if (hasRole) {
    console.log("interchainTokenManagerAddress has ADMIN_ROLE");
    return;
  }
  console.log("interchainTokenManagerAddress does not have ADMIN_ROLE");

  console.log("grantRole tx params:", {
    roleBytes,
    interchainTokenManagerAddress,
  });
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
