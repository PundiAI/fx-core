// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.10;

// NOTE: if using an interface to invoke the precompiled contract
// need to use solidity version 0.8.10 and later.
interface IBridgeOracle {
    function hasOracle(
        string memory _chain,
        address _externalAddress
    ) external view returns (bool _result);

    function isOracleOnline(
        string memory _chain,
        address _externalAddress
    ) external view returns (bool _result);
}
