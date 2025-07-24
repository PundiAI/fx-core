// @ts-ignore
import hre from "hardhat";
import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const DeployPundiAISwapModule = buildModule("DeployPundiAISwapModule", (m) => {
  if (hre.network.name === "ethereum") {
    const swap = m.contract("FXtoPUNDIAISwap", [
      "0x8c15ef5b4b21951d50e53e4fbda8298ffad25057",
      "0x075f23b9cdfce2cc0ca466f4ee6cb4bd29d83bef",
    ]);
    return { swap };
  } else if (hre.network.name === "sepolia") {
    const swap = m.contract("FXtoPUNDIAISwap", [
      "0xBb19939D96ca5cd34ec2919eE9Da3a1b70D7A77C",
      "0xebCb46E14bCd0F8639A24b32fBC6Db6935F046Fe",
    ]);
    return { swap };
  }
});

export default DeployPundiAISwapModule;
