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

async function main() {
  requireSourceChainTokenAddress();
  let destChainName = process.env.DESTINATION_CHAIN_NAME;
  if (!destChainName) {
    destChainName = destinationChainName;
  }
  console.log("destChainName:", destChainName);
  let crosschainTokenAmount: string = process.env.CROSSCHAIN_TOKEN_AMOUNT || "";
  if (crosschainTokenAmount === "") {
    throw new Error("CROSSCHAIN_TOKEN_AMOUNT environment variable is required");
  }
  const signer = await getSigner();
  let signerAddr = await signer.getAddress();

  const pundiaifxContract = new ethers.Contract(
    sourceChainTokenAddress,
    interchainTokenABI,
    signer
  );
  const allowance = await pundiaifxContract.allowance(
    signerAddr,
    interchainTokenServiceContractAddress
  );
  if (allowance < ethers.parseEther(crosschainTokenAmount)) {
    console.log(
      "Approving allowance for interchain token service contract...",
        {
          address:interchainTokenServiceContractAddress,
          amount: ethers.parseEther(crosschainTokenAmount).toString(),
        }
    );
    const approveTx = await pundiaifxContract.approve(
      interchainTokenServiceContractAddress,
      ethers.parseEther(crosschainTokenAmount)
    );
    console.log("approve tx:", approveTx.hash);
    await waitForTransaction(approveTx);
  } else {
    console.log("Allowance already sufficient:", allowance.toString());
  }
  let tokenId = process.env.TOKEN_ID;

  if (tokenId === undefined || tokenId === "") {
    const interchainTokenFactoryContract = new ethers.Contract(
      interchainTokenFactoryContractAddress,
      interchainTokenFactoryContractABI,
      signer
    );
    tokenId = await interchainTokenFactoryContract.linkedTokenId(
      signerAddr,
      salt
    );
  }

  const interchainTokenService = new ethers.Contract(
    interchainTokenServiceContractAddress,
    interchainTokenServiceContractABI,
    signer
  );
  console.log("interchainTransfer params:", {
    tokenId,
    destinationChain: destChainName,
    destinationAddress: signerAddr,
    amount: ethers.parseEther(crosschainTokenAmount),
    metadata: "0x",
    gasValue: ethers.parseEther(txFee),
  });
  const interchainTransferTx = await interchainTokenService.interchainTransfer(
    tokenId,
    destChainName,
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

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
