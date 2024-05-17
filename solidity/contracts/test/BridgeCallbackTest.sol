// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {IBridgeCallback} from "../bridge/IBridgeCallback.sol";
import {IERC20MetadataUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import {SafeERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";

contract BridgeCallbackTest is IBridgeCallback {
    using SafeERC20Upgradeable for IERC20MetadataUpgradeable;

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
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts,
        bytes memory _data,
        bytes memory
    ) external override onlyFxBridge {
        for (uint256 i = 0; i < _tokens.length; i++) {
            IERC20MetadataUpgradeable(_tokens[i]).transferFrom(
                _receiver,
                address(this),
                _amounts[i]
            );
        }

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
