import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const AirdropModule = buildModule("AirdropModule", (m) => {
  console.log("🚀 Starting Airdrop deployment...");

  const airdropAddress = m.getParameter("airdropAddress");
  const tokenAddress = m.getParameter("tokenAddress", "");
  const recipients = m.getParameter("recipients");
  const amounts = m.getParameter("amounts");
  const totalAmount = m.getParameter("totalAmount");

  console.log(`📋 Airdrop Configuration:`);
  console.log(`🏗️  Airdrop Contract: ${airdropAddress.toString()}`);
  console.log(`🪙 Token Address: ${tokenAddress.toString()}`);
  console.log(`👥 Recipients and amounts configured`);
  console.log(`💰 Airdrop distribution configured`);

  const token = m.contractAt("ERC20Upgradable", tokenAddress);
  const transferTx = m.call(token, "transfer", [airdropAddress, totalAmount], {
    id: "transfer_tokens_to_airdrop",
  });

  const airdrop = m.contractAt("Airdrop", airdropAddress);

  const distributeTokens = m.call(
    airdrop,
    "distributeTokens",
    [
      tokenAddress, // IERC20 token
      recipients, // address[] recipients
      amounts, // uint256[] amounts
    ],
    {
      id: "distribute_tokens",
      after: [transferTx],
    }
  );

  console.log(`✅ Airdrop distribution configured successfully!`);
  console.log(`📦 Using existing contract: ${airdropAddress.toString()}`);
  console.log(`🎯 Distribution call: ${distributeTokens}`);

  return { airdrop };
});

export default AirdropModule;
