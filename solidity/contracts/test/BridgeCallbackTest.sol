// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {IBridgeCallback} from "../bridge/IBridgeCallback.sol";

/* solhint-disable custom-errors */

contract BridgeCallbackTest is IBridgeCallback {
    address public fxBridge;
    bool public callFlag;

    constructor(address _fxBridge) {
        fxBridge = _fxBridge;
        callFlag = false;
    }

    function bridgeCallback(
        address,
        address,
        address[] memory,
        uint256[] memory,
        bytes memory,
        bytes memory
    ) external override onlyFxBridge {
        callFlag = !callFlag;
    }

    modifier onlyFxBridge() {
        require(msg.sender == fxBridge, "only fx bridge");
        _;
    }
}
