import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import {
  AbiCoder,
  encodeBytes32String,
  Interface,
  keccak256,
  Signature,
} from "ethers";
import { arrayify } from "@ethersproject/bytes";
import { FxBridgeLogic } from "../typechain-types";

export function examplePowers(): number[] {
  return [
    3000, 2000, 900, 800, 700, 600, 500, 400, 300, 200, 200, 200, 200, 200, 200,
    100, 100, 100, 100, 100,
  ];
}

export async function getSignerAddresses(signers: HardhatEthersSigner[]) {
  return await Promise.all(signers.map((signer) => signer.getAddress()));
}

export function makeSubmitBridgeCallHash(
  gravityId: string,
  sender: string,
  receiver: string,
  tokens: string[],
  amounts: string[],
  to: string,
  data: string,
  memo: string,
  nonce: number | string,
  timeout: number | string,
  eventNonce: number | string
) {
  let methodName = encodeBytes32String("bridgeCall");
  let abiCoder = new AbiCoder();
  return keccak256(
    abiCoder.encode(
      [
        "bytes32",
        "bytes32",
        "address",
        "address",
        "address[]",
        "uint256[]",
        "address",
        "bytes",
        "bytes",
        "uint256",
        "uint256",
        "uint256",
      ],
      [
        gravityId,
        methodName,
        sender,
        receiver,
        tokens,
        amounts,
        to,
        data,
        memo,
        nonce,
        timeout,
        eventNonce,
      ]
    )
  );
}

export async function signHash(signers: HardhatEthersSigner[], hash: string) {
  let v: number[] = [];
  let r: string[] = [];
  let s: string[] = [];

  const signMessage = arrayify(hash);
  for (let i = 0; i < signers.length; i = i + 1) {
    const sig = await signers[i].signMessage(signMessage);
    const signature = Signature.from(sig);

    v.push(signature.v);
    r.push(signature.r);
    s.push(signature.s);
  }
  return { v, r, s };
}

export function encodeFunctionData(abi: string, funcName: string, args: any[]) {
  const iface = new Interface(abi);
  return iface.encodeFunctionData(funcName, args);
}

export async function submitBridgeCall(
  gravityId: string,
  sender: string,
  refund: string,
  to: string,
  data: string,
  memo: string,
  tokens: string[],
  amounts: string[],
  nonce: number | string,
  timeout: number,
  eventNonce: number,
  validators: HardhatEthersSigner[],
  powers: number[],
  fxBridge: FxBridgeLogic
) {
  const digest = makeSubmitBridgeCallHash(
    gravityId,
    sender,
    refund,
    tokens,
    amounts,
    to,
    data,
    memo,
    nonce,
    timeout,
    eventNonce
  );

  const { v, r, s } = await signHash(validators, digest);

  const valAddresses = await getSignerAddresses(validators);
  const bridgeCallData: FxBridgeLogic.BridgeCallDataStruct = {
    sender: sender,
    refund: refund,
    tokens: tokens,
    amounts: amounts,
    to: to,
    data: data,
    memo: memo,
    timeout: timeout,
    eventNonce: eventNonce,
  };
  return await fxBridge.submitBridgeCall(
    valAddresses,
    powers,
    v,
    r,
    s,
    [0, nonce],
    bridgeCallData
  );
}
