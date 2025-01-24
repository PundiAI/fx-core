// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.10;

import {IBridgeCall} from "./IBridgeCall.sol";
import {IBridgeOracle} from "./IBridgeOracle.sol";

// NOTE: if using an interface to invoke the precompiled contract
// need to use solidity version 0.8.10 and later.
interface ICrosschain is IBridgeCall, IBridgeOracle {
    // Deprecated: please use `IBridgeCall.bridgeCall`
    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable returns (bool _result);

    function bridgeCoinAmount(
        address _token,
        bytes32 _target
    ) external view returns (uint256 _amount);

    function executeClaim(
        string memory _chain,
        uint256 _eventNonce
    ) external returns (bool _result);

    function getERC20Token(
        bytes32 _denom
    ) external view returns (address _token, bool _enable);

    // Deprecated
    event CrossChain(
        address indexed sender,
        address indexed token,
        string denom,
        string receipt,
        uint256 amount,
        uint256 fee,
        bytes32 target,
        string memo
    );

    event ExecuteClaimEvent(
        address indexed _sender,
        uint256 _eventNonce,
        string _chain,
        string _errReason
    );
}
