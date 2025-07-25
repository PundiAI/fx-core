import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

export default buildModule("upgradeBridgeETH", (m) => {
  const bridgeAddress = m.getParameter("bridgeAddress");

  console.log(
    `🏗️  Upgrading Bridge Contract (ETH): ${bridgeAddress.toString()}`
  );

  const newBridgeLogic = m.contract("FxBridgeLogicETH", [], {
    id: "FxBridgeLogicETH",
  });

  console.log(`📦 New Bridge Logic (ETH) deployed`);

  const proxyContract = m.contractAt(
    "ITransparentUpgradeableProxy",
    bridgeAddress
  );

  m.call(proxyContract, "upgradeTo", [newBridgeLogic]);

  console.log(`🔄 Bridge contract upgraded to new logic (ETH)`);

  const bridgeWrapper = m.contract("FxBridgeWrapper", [bridgeAddress], {
    id: "FxBridgeWrapperETH",
  });

  console.log(`🎁 Bridge Wrapper (ETH) deployed`);

  return {
    newBridgeLogic,
    proxyContract,
    bridgeWrapper,
  };
});
