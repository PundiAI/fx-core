import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const WithdrawAirdropModule = buildModule("WithdrawAirdropModule", (m) => {
  const airdropAddress = m.getParameter("airdropAddress");
  const tokenAddress = m.getParameter(
    "tokenAddress",
    "0x7a986bA67227acfab86385FF33436a80E2BB4CC5"
  );
  const totalAmount = m.getParameter("totalAmount");

  const airdrop = m.contractAt("Airdrop", airdropAddress);

  m.call(airdrop, "withdrawTokens", [tokenAddress, totalAmount], {
    id: "withdraw_tokens",
  });

  return { airdrop };
});

export default WithdrawAirdropModule;
