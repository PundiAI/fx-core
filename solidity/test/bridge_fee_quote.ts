import { ethers } from "hardhat";
import { expect } from "chai";
import { HDNodeWallet, Wallet } from "ethers";
import {
  BridgeFeeOracle,
  BridgeFeeQuote,
  BridgeFeeQuoteTest,
  IBridgeFeeQuote,
} from "../typechain-types";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";

const messagePrefix = "\x19Ethereum Signed Message:\n32";

describe("BridgeFeeQuoteUpgradeable", function () {
  let bridgeFeeQuote: BridgeFeeQuote;
  let bridgeFeeQuoteTest: BridgeFeeQuoteTest;
  let bridgeFeeOracle: BridgeFeeOracle;
  let chainName = "TestChain";
  let oracle: HDNodeWallet;
  let owner: HardhatEthersSigner;
  let token1: string;
  let token2: string;

  beforeEach(async function () {
    oracle = Wallet.createRandom();

    [owner] = await ethers.getSigners();

    token1 = "FX";
    token2 = "usdt";

    const BridgeFeeQuoteTest = await ethers.getContractFactory(
      "BridgeFeeQuoteTest"
    );
    bridgeFeeQuoteTest = await BridgeFeeQuoteTest.deploy();
    const quoteTest = await bridgeFeeQuoteTest.getAddress();

    const state: BridgeFeeQuoteTest.OracleStateStruct = {
      registered: true,
      online: true,
    };

    await bridgeFeeQuoteTest.setOracle(oracle.getAddress(), state);

    const BridgeFeeOracle = await ethers.getContractFactory("BridgeFeeOracle");

    bridgeFeeOracle = await BridgeFeeOracle.deploy();
    await bridgeFeeOracle.initialize(bridgeFeeQuoteTest.getAddress());

    const BridgeFeeQuote = await ethers.getContractFactory("BridgeFeeQuote");
    bridgeFeeQuote = await BridgeFeeQuote.deploy();

    await bridgeFeeQuote.initialize(bridgeFeeOracle.getAddress(), 3);

    const role = await bridgeFeeOracle.QUOTE_ROLE();
    await bridgeFeeOracle.grantRole(role, bridgeFeeQuote.getAddress());

    await bridgeFeeQuote.registerChain(chainName, []);
    await bridgeFeeQuote.registerTokenName(chainName, [token1, token2]);
  });

  describe("Oracle Management", function () {
    it("should block an oracle correctly", async function () {
      await bridgeFeeOracle.blackOracle(oracle.address);
      const oracleStatus = await bridgeFeeOracle.oracleStatus(oracle.address);
      expect(oracleStatus.isBlacklisted).to.be.true;
    });
  });

  describe("Quote Management", function () {
    it("should create a new quote", async function () {
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input = await newBridgeFeeQuote(
        chainName,
        token1,
        fee,
        gasLimit,
        expiry,
        oracle,
        0
      );

      await expect(bridgeFeeQuote.quote([input]))
        .to.emit(bridgeFeeQuote, "NewQuote")
        .withArgs(
          0,
          input.oracle,
          chainName,
          input.tokenName,
          input.fee,
          input.gasLimit,
          input.expiry
        );

      const quoteList = await bridgeFeeQuote.getQuoteList(chainName);
      expect(quoteList.length).to.be.equal(1);
    });

    it("should revert when trying to get quotes for an inactive chain", async function () {
      const chainName = "InactiveChain";
      await expect(
        bridgeFeeQuote.getQuoteList(chainName)
      ).to.be.revertedWithCustomError(bridgeFeeQuote, "ChainNameInvalid");
    });

    it("should revert when trying to create a quote with an expired expiry", async function () {
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) - 3600;

      const signature = await generateSignature(
        chainName,
        token1,
        fee,
        gasLimit,
        expiry,
        oracle
      );

      const quoteInput: IBridgeFeeQuote.QuoteInputStruct = {
        chainName: chainName,
        tokenName: token1,
        oracle: oracle.address,
        quoteIndex: 0,
        fee: fee,
        gasLimit: gasLimit,
        expiry: expiry,
        signature: signature,
      };
      await expect(
        bridgeFeeQuote.quote([quoteInput])
      ).to.be.revertedWithCustomError(bridgeFeeQuote, "QuoteExpired");
    });

    it("should revert when trying to create a quote without new oracle", async function () {
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input = await newBridgeFeeQuote(
        chainName,
        token1,
        fee,
        gasLimit,
        expiry,
        oracle,
        0
      );
      const input2 = await newBridgeFeeQuote(
        chainName,
        token2,
        fee,
        gasLimit,
        expiry,
        oracle,
        1
      );
      await bridgeFeeQuote.quote([input, input2]);

      const quoteList = await bridgeFeeQuote.getQuoteList(chainName);
      expect(quoteList.length).to.be.equal(2);

      const oracles = await bridgeFeeOracle.getOracleList();
      expect(oracles.length).to.be.equal(1);
    });

    it("test 1 ~ 5 quote gas limit", async function () {
      const number = 5;
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;
      const singers = await ethers.getSigners();
      let tokens: string[] = [];
      for (let i = 0; i < number; i++) {
        tokens.push("test" + i.toString());
      }
      await bridgeFeeQuote.registerTokenName(chainName, tokens);
      let quoteList: IBridgeFeeQuote.QuoteInputStruct[] = [];
      for (let i = 0; i < number; i++) {
        const input = await newBridgeFeeQuote(
          chainName,
          "test" + i.toString(),
          fee,
          gasLimit,
          expiry,
          oracle,
          0
        );
        quoteList.push(input);
        await bridgeFeeQuote.quote(quoteList);

        const quoteL = await bridgeFeeQuote.getQuoteList(chainName);
        expect(quoteL.length).to.be.equal(i + 1);
      }
    });

    it("first oracle quote", async function () {
      await bridgeFeeOracle.setDefaultOracle(oracle);

      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      await bridgeFeeQuote.quote([
        {
          chainName: chainName,
          tokenName: token1,
          oracle: oracle.address,
          quoteIndex: 0,
          fee: fee,
          gasLimit: gasLimit,
          expiry: expiry,
          signature: await generateSignature(
            chainName,
            token1,
            fee,
            gasLimit,
            expiry,
            oracle
          ),
        },
      ]);

      const quotes = await bridgeFeeQuote.getQuoteByToken(chainName, token1);
      expect(quotes.length).to.equal(3);

      for (let i = 0; i < 3; i++) {
        if (i == 0) {
          expect(quotes[i].fee).to.equal(fee);
          expect(quotes[i].gasLimit).to.equal(gasLimit);
          expect(quotes[i].expiry).to.equal(expiry);
        } else {
          expect(quotes[i].fee).to.equal(0);
          expect(quotes[i].gasLimit).to.equal(0);
          expect(quotes[i].expiry).to.equal(0);
        }
      }
    });

    it("oracle quote index", async function () {
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;
      let quoteList: IBridgeFeeQuote.QuoteInputStruct[] = [];
      for (let i = 0; i < 3; i++) {
        const input = await newBridgeFeeQuote(
          chainName,
          token1,
          fee,
          gasLimit,
          expiry,
          oracle,
          i
        );
        quoteList.push(input);
        await bridgeFeeQuote.quote(quoteList);
        const quoteL = await bridgeFeeQuote.getQuoteList(chainName);
        expect(quoteL.length).to.be.equal(i + 1);
      }
    });
  });
});

async function generateSignature(
  chainName: string,
  tokenName: string,
  fee: number,
  gasLimit: number,
  expiry: number,
  wallet: HDNodeWallet
): Promise<string> {
  const abiCoder = new ethers.AbiCoder();
  const coderHash = abiCoder.encode(
    ["string", "string", "uint256", "uint256", "uint256"],
    [chainName, tokenName, fee, gasLimit, expiry]
  );
  const hash = ethers.keccak256(coderHash);
  const messageHash = ethers.solidityPackedKeccak256(
    ["string", "bytes32"],
    [messagePrefix, hash]
  );

  const signatureW = wallet.signingKey.sign(messageHash);
  let v = "0x1b";
  if (signatureW.v === 28) {
    v = "0x1c";
  }
  return ethers.concat([signatureW.r, signatureW.s, v]);
}

async function newBridgeFeeQuote(
  chainName: string,
  tokenName: string,
  fee: number,
  gasLimit: number,
  expiry: number,
  oracle: HDNodeWallet,
  index: number
): Promise<IBridgeFeeQuote.QuoteInputStruct> {
  const signature = await generateSignature(
    chainName,
    tokenName,
    fee,
    gasLimit,
    expiry,
    oracle
  );
  return {
    chainName: chainName,
    tokenName: tokenName,
    oracle: oracle.address,
    quoteIndex: index,
    fee: fee,
    gasLimit: gasLimit,
    expiry: expiry,
    signature: signature,
  };
}

async function currentTime(): Promise<number> {
  const blockNumber = await ethers.provider.getBlockNumber();
  const block = await ethers.provider.getBlock(blockNumber);
  return block ? block.timestamp : Math.floor(Date.now() / 1000);
}
