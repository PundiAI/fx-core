// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {IBridgeCallContext} from "../interfaces/IBridgeCallContext.sol";

/* solhint-disable custom-errors */

contract BridgeCallContextTest is IBridgeCallContext {
    address public fxBridge;
    bool public callFlag;
    bool public revertFlag;

    constructor(address _fxBridge) {
        fxBridge = _fxBridge;
        callFlag = false;
        revertFlag = false;
    }

    function onBridgeCall(
        address,
        address,
        address[] memory,
        uint256[] memory,
        bytes memory,
        bytes memory
    ) external override onlyFxBridge {
        callFlag = !callFlag;
    }

    function onRevert(uint256, bytes memory) external override onlyFxBridge {
        revertFlag = !revertFlag;
    }

    modifier onlyFxBridge() {
        require(msg.sender == fxBridge, "only fx bridge");
        _;
    }
}
