// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {IBridgeCallContext} from "../interfaces/IBridgeCallContext.sol";
import {IFxBridgeLogic} from "../interfaces/IFxBridgeLogic.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/* solhint-disable custom-errors */

contract BridgeCallContextTest is IBridgeCallContext {
    address public fxBridge;
    address public receiver;

    bool public retryBridgeCall;

    string public dstChain;
    address public refund;
    address[] public tokens;
    uint256[] public amounts;
    address public to;
    bytes public data;
    uint256 public quoteId;
    uint256 public gasLimit;
    string public memo;

    constructor(address _fxBridge, address _receiver) {
        fxBridge = _fxBridge;
        receiver = _receiver;
    }

    function onBridgeCall(
        address,
        address,
        address[] memory _tokens,
        uint256[] memory _amounts,
        bytes memory,
        bytes memory
    ) external override onlyFxBridge {
        for (uint256 i = 0; i < _tokens.length; i++) {
            IERC20(_tokens[i]).transfer(receiver, _amounts[i]);
        }
    }

    function onRevert(uint256, bytes memory) external override onlyFxBridge {
        for (uint256 i = 0; i < tokens.length; i++) {
            IERC20(tokens[i]).approve(fxBridge, amounts[i]);
        }

        if (!retryBridgeCall) {
            return;
        }

        IFxBridgeLogic(fxBridge).bridgeCall(
            dstChain,
            refund,
            tokens,
            amounts,
            to,
            data,
            quoteId,
            gasLimit,
            bytes(memo)
        );
    }

    function setRetryBridgeCall(bool _retry) external {
        retryBridgeCall = _retry;
    }

    function setBridgeCallParams(
        string memory _dstChain,
        address _refund,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        uint256 _quoteId,
        uint256 _gasLimit,
        string memory _memo
    ) external {
        require(_tokens.length == _amounts.length, "length mismatch");
        dstChain = _dstChain;
        refund = _refund;
        tokens = _tokens;
        amounts = _amounts;
        to = _to;
        data = _data;
        quoteId = _quoteId;
        gasLimit = _gasLimit;
        memo = _memo;
    }

    modifier onlyFxBridge() {
        require(msg.sender == fxBridge, "only fx bridge");
        _;
    }
}
