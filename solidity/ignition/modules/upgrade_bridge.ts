import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

export default buildModule("upgradeBridge", (m) => {
  const bridgeAddress = m.getParameter("bridgeAddress");

  console.log(`🏗️  Upgrading Bridge Contract: ${bridgeAddress.toString()}`);

  const newBridgeLogic = m.contract("FxBridgeLogicBSC", [], {
    id: "FxBridgeLogicBSC",
  });

  console.log(`📦 New Bridge Logic deployed`);

  const proxyContract = m.contractAt(
    "ITransparentUpgradeableProxy",
    bridgeAddress
  );

  m.call(proxyContract, "upgradeTo", [newBridgeLogic]);

  console.log(`🔄 Bridge contract upgraded to new logic`);

  const bridgeWrapper = m.contract("FxBridgeWrapper", [bridgeAddress], {
    id: "FxBridgeWrapper",
  });

  console.log(`🎁 Bridge Wrapper deployed`);

  return {
    newBridgeLogic,
    proxyContract,
    bridgeWrapper,
  };
});
