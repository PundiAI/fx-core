import { ethers, keccak256, toUtf8Bytes } from "ethers";
import {
  sourceChainTokenAddress,
  interchainTokenABI,
  interchainTokenManagerAddress,
  getSigner,
  waitForTransaction,
  requireSourceChainTokenAddress,
  interchainTokenServiceContractAddress,
} from "./common";

async function main() {
  requireSourceChainTokenAddress();

  if (interchainTokenManagerAddress === "") {
    throw new Error("TOKEN_MANAGER_ADDRESS environment variable is required");
  }

  let tokenId = process.env.TOKEN_ID || "";
  if (tokenId === "") {
    throw new Error("TOKEN_ID environment variable is required");
  }

  const signer = await getSigner();
  const signerAddr = await signer.getAddress();

  const interchainTokenContract = new ethers.Contract(
    sourceChainTokenAddress,
    interchainTokenABI,
    signer
  );

  let roleBytes = keccak256(toUtf8Bytes("ADMIN_ROLE"));
  let hasRole = await interchainTokenContract.hasRole(roleBytes, signerAddr);
  if (!hasRole) {
    console.log("grantRole tx params:", {
      roleBytes,
      signerAddr,
    });
    const grantRoleTx = await interchainTokenContract.grantRole(
      roleBytes,
      signerAddr
    );
    console.log("grantRole tx:", grantRoleTx.hash);
  }

  const interchainTokenId = await interchainTokenContract.interchainTokenId();
  if (interchainTokenId !== tokenId) {
    console.log("setTokenId tx params:", {
      tokenId,
    });
    const setTokenIdTx = await interchainTokenContract.setTokenId(tokenId);
    console.log("setTokenId tx:", setTokenIdTx.hash);
    await waitForTransaction(setTokenIdTx);
  }

  const interchainTokenService =
    await interchainTokenContract.interchainTokenService();
  if (interchainTokenService !== interchainTokenServiceContractAddress) {
    console.log("setInterchainTokenService tx params:", {
      interchainTokenServiceContractAddress,
    });
    const setItsTx = await interchainTokenContract.setInterchainTokenService(
      interchainTokenServiceContractAddress
    );
    console.log("setInterchainTokenService tx:", setItsTx.hash);
    await waitForTransaction(setItsTx);
  }

  roleBytes = keccak256(toUtf8Bytes("ADMIN_ROLE"));
  hasRole = await interchainTokenContract.hasRole(
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
  const grantRoleTx = await interchainTokenContract.grantRole(
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
