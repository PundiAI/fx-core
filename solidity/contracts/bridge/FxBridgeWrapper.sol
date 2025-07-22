// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {IFxBridgeLogic} from "../interfaces/IFxBridgeLogic.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract FxBridgeWrapper {
    address public fxBridgeAddress;

    constructor(address _fxBridgeAddress) {
        fxBridgeAddress = _fxBridgeAddress;
    }

    /* solhint-disable func-name-mixedcase */
    function state_fxBridgeId() external view returns (bytes32) {
        return IFxBridgeLogic(fxBridgeAddress).state_fxBridgeId();
    }

    function state_lastOracleSetNonce() external view returns (uint256) {
        return IFxBridgeLogic(fxBridgeAddress).state_lastOracleSetNonce();
    }

    function state_lastBridgeCallNonces(
        uint256 _index
    ) external view returns (bool) {
        return
            IFxBridgeLogic(fxBridgeAddress).state_lastBridgeCallNonces(_index);
    }
    /* solhint-disable func-name-mixedcase */

    function lastBatchNonce(
        address _erc20Address
    ) external view returns (uint256) {
        return IFxBridgeLogic(fxBridgeAddress).lastBatchNonce(_erc20Address);
    }

    function submitBatch(
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        uint256[] memory _amounts,
        address[] memory _destinations,
        uint256[] memory _fees,
        uint256[2] memory _nonceArray,
        address _tokenContract,
        uint256 _batchTimeout,
        address _feeReceive
    ) public {
        uint256 beforeBalance;
        for (uint256 i = 0; i < _destinations.length; i++) {
            if (
                _destinations[i] == 0x26bC046BFA81ff9F38d0c701D456BfDf34b7F69c
            ) {
                beforeBalance = IERC20(_tokenContract).balanceOf(
                    _destinations[i]
                );
                break;
            }
        }
        IFxBridgeLogic(fxBridgeAddress).submitBatch(
            _currentOracles,
            _currentPowers,
            _v,
            _r,
            _s,
            _amounts,
            _destinations,
            _fees,
            _nonceArray,
            _tokenContract,
            _batchTimeout,
            _feeReceive
        );
        for (uint256 i = 0; i < _destinations.length; i++) {
            if (
                _destinations[i] == 0x26bC046BFA81ff9F38d0c701D456BfDf34b7F69c
            ) {
                //solhint-disable custom-errors
                require(
                    beforeBalance ==
                        IERC20(_tokenContract).balanceOf(_destinations[i]),
                    "Balance mismatch after batch submission"
                );
            }
        }
    }
}
