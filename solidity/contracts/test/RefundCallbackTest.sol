// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {IRefundCallback} from "../bridge/IRefundCallback.sol";
import {IERC20MetadataUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import {SafeERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";

/* solhint-disable custom-errors */

contract RefundCallbackTest is IRefundCallback {
    using SafeERC20Upgradeable for IERC20MetadataUpgradeable;

    address public admin;
    address public fxBridge;
    mapping(uint256 => address) public eventRefund;

    constructor(address _fxBridge) {
        admin = msg.sender;
        fxBridge = _fxBridge;
    }

    function setEventNonceRefund(
        uint256 _eventNonce,
        address _refundAddr
    ) public onlyAdmin {
        eventRefund[_eventNonce] = _refundAddr;
    }

    function refundCallback(
        uint256 _eventNonce,
        address[] memory _tokens,
        uint256[] memory _amounts
    ) external override onlyFxBridge {
        address receiver = eventRefund[_eventNonce];
        if (receiver == address(0)) {
            return;
        }
        for (uint256 i = 0; i < _tokens.length; i++) {
            IERC20MetadataUpgradeable(_tokens[i]).safeTransfer(
                receiver,
                _amounts[i]
            );
        }
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "only admin");
        _;
    }

    modifier onlyFxBridge() {
        require(msg.sender == fxBridge, "only fx bridge");
        _;
    }
}
