// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {IBridgeCallback} from "../bridge/IBridgeCallback.sol";

contract BridgeCallbackTest is IBridgeCallback {
    address public fxBridge;
    address public admin;
    mapping(address => bool) public whiteList;

    constructor(address _fxBridge) {
        fxBridge = _fxBridge;
        admin = msg.sender;
    }

    function addWhiteList(address _to) public onlyAdmin {
        whiteList[_to] = true;
    }

    function bridgeCallback(
        address,
        address,
        address[] memory,
        uint256[] memory,
        bytes memory _data,
        bytes memory
    ) external override onlyFxBridge {
        (address to, bytes memory data) = abi.decode(_data, (address, bytes));
        // solhint-disable custom-errors
        require(whiteList[to], "not in white list");
        // solhint-disable avoid-low-level-calls
        (bool success, ) = to.call(data);
        // solhint-disable custom-errors
        require(success, "failed to call data");
    }

    modifier onlyFxBridge() {
        require(msg.sender == fxBridge, "only fx bridge");
        _;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "only admin");
        _;
    }
}
