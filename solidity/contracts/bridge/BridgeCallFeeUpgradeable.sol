// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IERC20MetadataUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import {SafeERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";

import {IBridgeCallback} from "./IBridgeCallback.sol";

/* solhint-disable custom-errors */

contract BridgeCallFeeUpgradeable is
    IBridgeCallback,
    Initializable,
    ContextUpgradeable,
    UUPSUpgradeable,
    OwnableUpgradeable
{
    using SafeERC20Upgradeable for IERC20MetadataUpgradeable;

    address public bridge;
    mapping(address => mapping(address => uint256)) public fee;

    function updateBridge(address _bridge) public onlyOwner {
        bridge = _bridge;
    }

    function bridgeCallback(
        address _sender,
        address _refund,
        address[] calldata _tokens,
        uint256[] calldata _amounts,
        bytes memory _data,
        bytes memory _memo
    ) external override onlyBridge {
        uint256 tokenLen = _tokens.length;

        // deduction of fee
        if (tokenLen > 0) {
            address feeToken = _tokens[tokenLen - 1];
            uint256 feeAmount = _amounts[tokenLen - 1];
            // solhint-disable-next-line avoid-tx-origin
            fee[tx.origin][feeToken] += feeAmount;
            emit BridgeCallFeeEvent(
                // solhint-disable-next-line avoid-tx-origin
                tx.origin,
                feeToken,
                feeAmount
            );
            tokenLen = tokenLen - 1;
        }

        (address target, bytes memory targetMemo) = abi.decode(
            _memo,
            (address, bytes)
        );
        for (uint256 i = 0; i < tokenLen; i++) {
            IERC20MetadataUpgradeable(_tokens[i]).transfer(target, _amounts[i]);
        }
        IBridgeCallback(target).bridgeCallback(
            _sender,
            _refund,
            _tokens[0:tokenLen],
            _amounts[0:tokenLen],
            _data,
            targetMemo
        );
    }

    function withdraw(address[] memory _tokens) public {
        require(_tokens.length > 0, "empty tokens");
        for (uint256 i = 0; i < _tokens.length; i++) {
            uint256 feeAmount = fee[_msgSender()][_tokens[i]];
            if (feeAmount == 0) {
                continue;
            }
            fee[_msgSender()][_tokens[i]] = 0;
            IERC20MetadataUpgradeable(_tokens[i]).transfer(
                _msgSender(),
                feeAmount
            );
            emit WithdrawFeeEvent(_msgSender(), _tokens[i], feeAmount);
        }
    }

    modifier onlyBridge() {
        require(_msgSender() == bridge, "only bridge can call this function");
        _;
    }

    event BridgeCallFeeEvent(
        address indexed _receiver,
        address indexed _token,
        uint256 _amount
    );

    event WithdrawFeeEvent(
        address indexed _receiver,
        address indexed _token,
        uint256 _amount
    );

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function initialize(address _bridge) public virtual initializer {
        bridge = _bridge;

        __Ownable_init();
        __UUPSUpgradeable_init();
    }
}
