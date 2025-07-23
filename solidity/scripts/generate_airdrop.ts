import { ethers } from "hardhat";
import * as fs from "fs";
import * as path from "path";

interface AirdropConfig {
  airdropAddress: string;
  tokenAddress: string;
  recipients: string[];
  amounts: string[];
  startIndex: number;
  endIndex: number;
  totalRecipients: number;
  totalAmount: string;
}

interface IgnitionParams {
  AirdropModule: {
    airdropAddress: string;
    tokenAddress: string;
    recipients: string[];
    amounts: string[];
    totalAmount: string;
  };
}

async function main() {
  // startIndex from 0
  const startIndex = process.env.START_INDEX
    ? parseInt(process.env.START_INDEX)
    : 0;
  // endIndex Max 6711
  const endIndex = process.env.END_INDEX
    ? parseInt(process.env.END_INDEX)
    : 6711;
  const tokenAddress =
    process.env.TOKEN_ADDRESS || "0x7a986bA67227acfab86385FF33436a80E2BB4CC5";
  const airdropAddress = process.env.AIRDROP_ADDRESS || "";
  const airdropDataFile =
    process.env.AIRDROP_DATA_FILE || "./scripts/airdrop/addresses.json";

  console.log("🚀 Starting airdrop configuration generation...");
  console.log(`📍 Range: ${startIndex} to ${endIndex}`);
  console.log(`🪙 Token: ${tokenAddress}`);
  console.log(
    `🏗️  Airdrop Contract: ${
      airdropAddress || "Not specified - will need to be provided"
    }`
  );

  try {
    const airdropData = path.resolve(airdropDataFile);
    console.log(`📖 Reading: ${airdropData}`);

    if (!fs.existsSync(airdropData)) {
      throw new Error(`File not found: ${airdropData}`);
    }

    const airdropJSON = JSON.parse(fs.readFileSync(airdropData, "utf8"));
    console.log(`📊 Total address in file: ${airdropJSON.length}`);

    const filteredZeroBalance = airdropJSON.filter(
      (user: { balance_wei: string }) => {
        const balance = BigInt(user.balance_wei);
        return balance > 0n;
      }
    );

    console.log(
      `🔍 Filtered address (balance > 0): ${filteredZeroBalance.length}`
    );
    console.log(
      `❌ Excluded address (balance = 0): ${
        airdropJSON.length - filteredZeroBalance.length
      }`
    );

    if (startIndex < 0 || startIndex >= filteredZeroBalance.length) {
      throw new Error(
        `Invalid start index: ${startIndex}. Must be between 0 and ${
          filteredZeroBalance.length - 1
        }`
      );
    }

    const actualEndIndex = Math.min(endIndex, filteredZeroBalance.length - 1);
    if (startIndex > actualEndIndex) {
      throw new Error(
        `Start index (${startIndex}) cannot be greater than end index (${actualEndIndex})`
      );
    }

    const selectedAddr = filteredZeroBalance.slice(
      startIndex,
      actualEndIndex + 1
    );
    console.log(
      `✂️  Selected ${selectedAddr.length} (index ${startIndex} to ${actualEndIndex})`
    );

    const recipients: string[] = [];
    const amounts: string[] = [];

    console.log(`\n📋 Selected Recipients Details:`);
    selectedAddr.forEach((user: { address: string; balance_wei: string }) => {
      recipients.push(user.address);
      amounts.push(user.balance_wei);
    });

    const totalAmount = amounts.reduce((sum, amount) => {
      return sum + BigInt(amount);
    }, BigInt(0));

    const totalAmountInEth = ethers.formatEther(totalAmount.toString());

    const ignitionParams: IgnitionParams = {
      AirdropModule: {
        airdropAddress,
        tokenAddress,
        recipients,
        amounts,
        totalAmount: totalAmount.toString(),
      },
    };

    const airdropConfig: AirdropConfig = {
      airdropAddress,
      tokenAddress,
      recipients,
      amounts,
      startIndex,
      endIndex: actualEndIndex,
      totalRecipients: recipients.length,
      totalAmount: totalAmount.toString(),
    };

    const ignitionOutputPath = path.resolve(
      "./ignition/params/airdrop_params.json"
    );
    fs.writeFileSync(
      ignitionOutputPath,
      JSON.stringify(ignitionParams, null, 2)
    );

    const configOutputPath = path.resolve(
      "./ignition/params/airdrop_config.json"
    );
    fs.writeFileSync(configOutputPath, JSON.stringify(airdropConfig, null, 2));

    console.log(`\n✅ Airdrop configuration generated successfully!`);
    console.log(`📁 Ignition params: ${ignitionOutputPath}`);
    console.log(`📁 Detailed config: ${configOutputPath}`);
    console.log(
      `🏗️  Airdrop contract: ${
        airdropAddress ||
        "⚠️  REQUIRED: Set AIRDROP_ADDRESS environment variable"
      }`
    );
    console.log(`🪙 Token address: ${tokenAddress}`);
    console.log(`📍 Index range: ${startIndex} to ${actualEndIndex}`);
    console.log(`👥 Total recipients: ${recipients.length}`);
    console.log(`💰 Total amount: ${totalAmount.toString()} wei`);
    console.log(`💰 Total amount: ${totalAmountInEth} ETH`);

    if (!airdropAddress) {
      console.log(`\n⚠️  Before deployment, set the airdrop contract address:`);
      console.log(
        `   Edit ${ignitionOutputPath} and set the airdropAddress field`
      );
      console.log(
        `   OR use: AIRDROP_ADDRESS=0x... npx hardhat run scripts/generate_airdrop.ts`
      );
    }

    console.log(`\n🚀 Ready for deployment with:`);
    console.log(
      `   npx hardhat ignition deploy ./ignition/modules/airdrop.ts --parameters ${ignitionOutputPath}`
    );
  } catch (error) {
    console.error("❌ Error generating airdrop configuration:", error);
    process.exit(1);
  }
}

if (require.main === module) {
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });
}

export { main };
