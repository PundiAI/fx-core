// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

interface IBank {
    function transferFromModuleToAccount(
        string memory _module,
        address _account,
        address _token,
        uint256 _amount
    ) external returns (bool _result);

    function transferFromAccountToModule(
        address _account,
        string memory _module,
        address _token,
        uint256 _amount
    ) external returns (bool _result);
}
