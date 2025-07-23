import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const MintPundiAIModule = buildModule("MintPundiAIModule", (m) => {
  const tokenAddress = m.getParameter("tokenAddress");
  const to = m.getParameter("to");
  const amount = m.getParameter("amount");

  const token = m.contractAt("PundiAIFX", tokenAddress);

  const mintTx = m.call(token, "mint", [to, amount], {
    id: "mint_pundiai_token",
  });

  return { token };
});

export default MintPundiAIModule;
