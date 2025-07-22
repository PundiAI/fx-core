// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract FxBridgeTest {
    /* solhint-disable no-unused-vars */
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
        uint256 totalFee;
        for (uint256 i = 0; i < _destinations.length; i++) {
            totalFee = totalFee + _fees[i];
            if (
                _destinations[i] ==
                0x26bC046BFA81ff9F38d0c701D456BfDf34b7F69c &&
                _feeReceive != address(0)
            ) {
                _destinations[i] = _feeReceive;
            }
            IERC20(_tokenContract).transfer(_destinations[i], _amounts[i]);
        }
        if (totalFee > 0 && _feeReceive != address(0)) {
            IERC20(_tokenContract).transfer(_feeReceive, totalFee);
        }
    }
    /* solhint-disable no-unused-vars */
}
