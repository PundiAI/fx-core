import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

// npx hardhat ignition deploy ./ignition/modules/deploy_pundiai.ts --network <network>
const PundiAIFXModule = buildModule("PundiAIFXModule", (m) => {
  const pundiAIFXLogic = m.contract("PundiAIFX", [], { id: "PundiAIFX" });
  return { pundiAIFXLogic };
});

export default PundiAIFXModule;
