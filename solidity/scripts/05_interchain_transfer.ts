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
  sourceChainTokenAddress,
  interchainTokenABI,
  requireSourceChainTokenAddress,
} from "./common";

async function interchainTransfer(crosschainTokenAmount: string) {
  requireSourceChainTokenAddress();
  if (destinationChainName === "") {
    throw new Error("DESTINATION_CHAIN_NAME environment variable is required");
  }

  const tokenId = process.env.TOKEN_ID || "";
  if (tokenId === "") {
    throw new Error("TOKEN_ID environment variable is required");
  }

  const signer = await getSigner();
  const signerAddr = await signer.getAddress();

  const interchainToken = new ethers.Contract(
    sourceChainTokenAddress,
    interchainTokenABI,
    signer
  );
  const allowance = await interchainToken.allowance(
    signerAddr,
    interchainTokenServiceContractAddress
  );
  if (allowance < ethers.parseEther(crosschainTokenAmount)) {
    console.log("Approving allowance for interchain token service contract", {
      address: interchainTokenServiceContractAddress,
      amount: ethers.parseEther(crosschainTokenAmount).toString(),
    });
    const approveTx = await interchainToken.approve(
      interchainTokenServiceContractAddress,
      ethers.parseEther(crosschainTokenAmount)
    );
    console.log("approve tx:", approveTx.hash);
    await waitForTransaction(approveTx);
  } else {
    console.log("Allowance already sufficient:", allowance.toString());
  }

  const interchainTokenService = new ethers.Contract(
    interchainTokenServiceContractAddress,
    interchainTokenServiceContractABI,
    signer
  );
  console.log("interchainTransfer params:", {
    tokenId,
    destinationChain: destinationChainName,
    destinationAddress: signerAddr,
    amount: ethers.parseEther(crosschainTokenAmount),
    metadata: "0x",
    gasValue: ethers.parseEther(txFee),
  });
  const interchainTransferTx = await interchainTokenService.interchainTransfer(
    tokenId,
    destinationChainName,
    signerAddr,
    ethers.parseEther(crosschainTokenAmount),
    "0x",
    ethers.parseEther(txFee),
    {
      value: ethers.parseEther(txFee),
    }
  );
  console.log("interchainTransfer tx:", interchainTransferTx.hash);

  await waitForTransaction(interchainTransferTx);
}

async function main() {
  const crosschainTokenAmount: string =
    process.env.CROSSCHAIN_TOKEN_AMOUNT || "";
  if (crosschainTokenAmount === "") {
    throw new Error("CROSSCHAIN_TOKEN_AMOUNT environment variable is required");
  }

  if (process.env.TOKEN_ID !== undefined) {
    await interchainTransfer(crosschainTokenAmount);
    return;
  }

  const signer = await getSigner();
  const signerAddr = await signer.getAddress();

  const interchainToken = new ethers.Contract(
    sourceChainTokenAddress,
    interchainTokenABI,
    signer
  );
  const allowance = await interchainToken.allowance(
    signerAddr,
    interchainTokenServiceContractAddress
  );
  if (allowance < ethers.parseEther(crosschainTokenAmount)) {
    console.log("Approving allowance for interchain token service contract", {
      address: interchainTokenServiceContractAddress,
      amount: ethers.parseEther(crosschainTokenAmount).toString(),
    });
    const approveTx = await interchainToken.approve(
      interchainTokenServiceContractAddress,
      ethers.parseEther(crosschainTokenAmount)
    );
    console.log("approve tx:", approveTx.hash);
    await waitForTransaction(approveTx);
  } else {
    console.log("Allowance already sufficient:", allowance.toString());
  }
  console.log("interchainTransfer params:", {
    destinationChain: destinationChainName,
    destinationAddress: signerAddr,
    amount: ethers.parseEther(crosschainTokenAmount),
    metadata: "0x",
  });
  const interchainTransferTx = await interchainToken.interchainTransfer(
    destinationChainName,
    signerAddr,
    ethers.parseEther(crosschainTokenAmount),
    "0x"
  );
  console.log("interchainTransfer tx:", interchainTransferTx.hash);
  await waitForTransaction(interchainTransferTx);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
