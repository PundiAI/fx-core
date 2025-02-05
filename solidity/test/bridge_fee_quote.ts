import { ethers } from "hardhat";
import { expect } from "chai";
import {
  BridgeFeeOracle,
  BridgeFeeQuote,
  BridgeFeeQuoteTest,
  IBridgeFeeQuote,
} from "../typechain-types";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";

describe("BridgeFeeQuoteUpgradeable", function () {
  let bridgeFeeQuote: BridgeFeeQuote;
  let bridgeFeeQuoteTest: BridgeFeeQuoteTest;
  let bridgeFeeOracle: BridgeFeeOracle;
  let chainNameStr = "TestChain";
  let chainName = ethers.encodeBytes32String(chainNameStr);
  let owner: HardhatEthersSigner;
  let token1: string;
  let token2: string;

  beforeEach(async function () {
    [owner] = await ethers.getSigners();

    token1 = ethers.encodeBytes32String("TEST");
    token2 = ethers.encodeBytes32String("usdt");

    const BridgeFeeQuoteTest = await ethers.getContractFactory(
      "BridgeFeeQuoteTest"
    );
    bridgeFeeQuoteTest = await BridgeFeeQuoteTest.deploy();
    const quoteTest = await bridgeFeeQuoteTest.getAddress();

    const state: BridgeFeeQuoteTest.OracleStateStruct = {
      registered: true,
      online: true,
    };

    await bridgeFeeQuoteTest.setOracle(chainName, owner.getAddress(), state);

    const BridgeFeeOracle = await ethers.getContractFactory("BridgeFeeOracle");

    bridgeFeeOracle = await BridgeFeeOracle.deploy();
    await bridgeFeeOracle.initialize(bridgeFeeQuoteTest.getAddress());

    const BridgeFeeQuote = await ethers.getContractFactory("BridgeFeeQuote");
    bridgeFeeQuote = await BridgeFeeQuote.deploy();

    await bridgeFeeQuote.initialize(bridgeFeeOracle.getAddress(), 3);

    const role = await bridgeFeeOracle.QUOTE_ROLE();
    await bridgeFeeOracle.grantRole(role, bridgeFeeQuote.getAddress());

    await bridgeFeeQuote.registerChain(chainName, []);
    await bridgeFeeQuote.addToken(chainName, [token1, token2]);
  });

  describe("Oracle Management", function () {
    it("should block an oracle correctly", async function () {
      await bridgeFeeOracle.blackOracle(chainName, owner.address);
      const oracleStatus = await bridgeFeeOracle.oracleStatus(
        chainName,
        owner.address
      );
      expect(oracleStatus.isBlack).to.be.true;
    });
  });

  describe("Quote Management", function () {
    it("should create a new quote", async function () {
      const amount = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input: IBridgeFeeQuote.QuoteInputStruct = {
        cap: 0,
        gasLimit: gasLimit,
        expiry: expiry,
        chainName: chainName,
        tokenName: token1,
        amount: amount,
      };

      await expect(bridgeFeeQuote.quote([input]))
        .to.emit(bridgeFeeQuote, "NewQuote")
        .withArgs(
          1,
          chainName,
          token1,
          owner.address,
          amount,
          gasLimit,
          expiry,
          0
        );
      const quote = await bridgeFeeQuote.getQuoteById(1);
      expect(quote.id).to.equal(1);
      expect(quote.chainName).to.equal(chainName);
      expect(quote.tokenName).to.equal(token1);
      expect(quote.oracle).to.equal(owner.address);
      expect(quote.amount).to.equal(amount);
      expect(quote.gasLimit).to.equal(gasLimit);
      expect(quote.expiry).to.equal(expiry);

      const oracleList = await bridgeFeeOracle.getOracleList(chainName);
      const quoteList: IBridgeFeeQuote.QuoteInfoStructOutput[] =
        await bridgeFeeQuote.getQuotesByToken(chainName, token1);
      expect(quoteList.length).to.equal(3);
    });

    it("should revert when trying to create a quote without new oracle", async function () {
      const amount = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input: IBridgeFeeQuote.QuoteInputStruct = {
        cap: 0,
        gasLimit: gasLimit,
        expiry: expiry,
        chainName: chainName,
        tokenName: token1,
        amount: amount,
      };

      const input2: IBridgeFeeQuote.QuoteInputStruct = {
        cap: 1,
        gasLimit: gasLimit,
        expiry: expiry,
        chainName: chainName,
        tokenName: token1,
        amount: amount,
      };

      await bridgeFeeQuote.quote([input, input2]);

      const quoteList = await bridgeFeeQuote.getQuotesByToken(
        chainName,
        token1
      );
      expect(quoteList.length).to.be.equal(3);

      const oracles = await bridgeFeeOracle.getOracleList(chainName);
      expect(oracles.length).to.be.equal(1);
    });

    it("test get quote by index", async function () {
      const amount = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input: IBridgeFeeQuote.QuoteInputStruct = {
        cap: 0,
        gasLimit: gasLimit,
        expiry: expiry,
        chainName: chainName,
        tokenName: token1,
        amount: amount,
      };
      await bridgeFeeQuote.quote([input]);

      const quote = await bridgeFeeQuote.getQuoteByIndex(
        chainName,
        token1,
        owner.address,
        0
      );
      expect(quote.id).to.be.equal(1);
      expect(quote.chainName).to.be.equal(chainName);
      expect(quote.tokenName).to.be.equal(token1);
      expect(quote.oracle).to.be.equal(owner.address);
      expect(quote.amount).to.be.equal(amount);
      expect(quote.gasLimit).to.be.equal(gasLimit);
      expect(quote.expiry).to.be.equal(expiry);
    });

    it("test 1 ~ 5 quote gas limit", async function () {
      const number = 5;
      const amount = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;
      const singers = await ethers.getSigners();
      let tokens: string[] = [];
      for (let i = 0; i < number; i++) {
        tokens.push(ethers.encodeBytes32String("test" + i.toString()));
      }
      await bridgeFeeQuote.addToken(chainName, tokens);
      let quoteList: IBridgeFeeQuote.QuoteInputStruct[] = [];
      for (let i = 0; i < number; i++) {
        const input: IBridgeFeeQuote.QuoteInputStruct = {
          cap: 0,
          gasLimit: gasLimit,
          expiry: expiry,
          chainName: chainName,
          tokenName: tokens[i],
          amount: amount,
        };
        quoteList.push(input);
        await bridgeFeeQuote.quote(quoteList);
      }
    });

    it("first oracle quote", async function () {
      await bridgeFeeOracle.setDefaultOracle(owner);

      const amount = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input: IBridgeFeeQuote.QuoteInputStruct = {
        cap: 0,
        gasLimit: gasLimit,
        expiry: expiry,
        chainName: chainName,
        tokenName: token1,
        amount: amount,
      };

      await bridgeFeeQuote.quote([input]);

      const quotes = await bridgeFeeQuote.getQuotesByToken(chainName, token1);
      expect(quotes.length).to.equal(3);

      for (let i = 0; i < 3; i++) {
        if (i == 0) {
          expect(quotes[i].amount).to.equal(amount);
          expect(quotes[i].gasLimit).to.equal(gasLimit);
          expect(quotes[i].expiry).to.equal(expiry);
        } else {
          expect(quotes[i].amount).to.equal(0);
          expect(quotes[i].gasLimit).to.equal(0);
          expect(quotes[i].expiry).to.equal(0);
        }
      }
    });

    it("oracle quote index", async function () {
      const amount = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;
      let quoteList: IBridgeFeeQuote.QuoteInputStruct[] = [];
      for (let i = 0; i < 3; i++) {
        const input: IBridgeFeeQuote.QuoteInputStruct = {
          cap: i,
          gasLimit: gasLimit,
          expiry: expiry,
          chainName: chainName,
          tokenName: token1,
          amount: amount,
        };
        quoteList.push(input);
        await bridgeFeeQuote.quote(quoteList);

        const quoteL = await bridgeFeeQuote.getQuotesByToken(chainName, token1);
        expect(quoteL.length).to.be.equal(3);

        if (i == 0) {
          expect(quoteL[0].id).to.equal(1);
          expect(quoteL[1].id).to.equal(0);
          expect(quoteL[2].id).to.equal(0);
        }
        if (i == 1) {
          expect(quoteL[0].id).to.equal(2);
          expect(quoteL[1].id).to.equal(3);
          expect(quoteL[2].id).to.equal(0);
        }
        if (i == 2) {
          expect(quoteL[0].id).to.equal(4);
          expect(quoteL[1].id).to.equal(5);
          expect(quoteL[2].id).to.equal(6);
        }
      }
    });

    it("get quote by id", async function () {
      const amount = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input: IBridgeFeeQuote.QuoteInputStruct = {
        cap: 0,
        gasLimit: gasLimit,
        expiry: expiry,
        chainName: chainName,
        tokenName: token1,
        amount: amount,
      };

      await bridgeFeeQuote.quote([input]);
      expect((await bridgeFeeQuote.getQuoteById(1)).id).to.equal(1);
      await expect(bridgeFeeQuote.getQuoteById(2)).to.revertedWithCustomError(
        bridgeFeeQuote,
        "QuoteIdInvalid()"
      );
      await expect(bridgeFeeQuote.getQuoteById(0)).to.revertedWithCustomError(
        bridgeFeeQuote,
        "QuoteIdInvalid()"
      );
    });
  });
});

async function currentTime(): Promise<number> {
  const blockNumber = await ethers.provider.getBlockNumber();
  const block = await ethers.provider.getBlock(blockNumber);
  return block ? block.timestamp : Math.floor(Date.now() / 1000);
}
