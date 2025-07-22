import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

// npx hardhat ignition deploy ./ignition/modules/interchain_token.ts --network <network>
const InterchainTokenModule = buildModule("InterchainTokenModule", (m) => {
  const pundiAIFXInterchainTokenLogic = m.contract(
    "PundiAIFXInterchainToken",
    [],
    { id: "PundiAIFXInterchainToken" }
  );
  const initializeData = m.encodeFunctionCall(
    pundiAIFXInterchainTokenLogic,
    "initialize",
    []
  );
  const pundiAIFXInterchainTokenProxy = m.contract(
    "ERC1967Proxy",
    [pundiAIFXInterchainTokenLogic, initializeData],
    { id: "PundiAIFXInterchainTokenProxy" }
  );
  return { pundiAIFXInterchainTokenProxy };
});

export default InterchainTokenModule;
