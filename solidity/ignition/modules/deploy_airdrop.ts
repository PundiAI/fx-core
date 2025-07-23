import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const DeployAirdropModule = buildModule("DeployAirdropModule", (m) => {
  const airdrop = m.contract("Airdrop", []);
  return { airdrop };
});

export default DeployAirdropModule;
