// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.10;

import {IBridgeCall} from "./IBridgeCall.sol";

// NOTE: if using an interface to invoke the precompiled contract
// need to use solidity version 0.8.10 and later.
interface ICrossChain is IBridgeCall {
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

    function hasOracle(
        string memory _chain,
        address _externalAddress
    ) external view returns (bool _result);

    function isOracleOnline(
        string memory _chain,
        address _externalAddress
    ) external view returns (bool _result);

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

    event BridgeCallEvent(
        address indexed _sender,
        address indexed _receiver,
        address indexed _to,
        address _txOrigin,
        uint256 _value,
        uint256 _eventNonce,
        string _dstChain,
        address[] _tokens,
        uint256[] _amounts,
        bytes _data,
        bytes _memo
    );

    event ExecuteClaimEvent(
        address indexed _sender,
        uint256 _eventNonce,
        string _chain
    );
}
