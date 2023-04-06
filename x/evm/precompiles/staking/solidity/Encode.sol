// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

library Encode {
    function delegate(
        string memory _validator
    ) internal pure returns (bytes memory) {
        return abi.encodeWithSignature("delegate(string)", _validator);
    }

    function undelegate(
        string memory _validator,
        uint256 _shares
    ) internal pure returns (bytes memory) {
        return
        abi.encodeWithSignature(
            "undelegate(string,uint256)",
            _validator,
            _shares
        );
    }

    function withdraw(
        string memory _validator
    ) internal pure returns (bytes memory) {
        return abi.encodeWithSignature("withdraw(string)", _validator);
    }

    function approve(
        string memory _validator,
        address _spender,
        uint256 _shares
    ) internal pure returns (bytes memory) {
        return
        abi.encodeWithSignature(
            "approve(string,address,uint256)",
            _validator,
            _spender,
            _shares
        );
    }

    function transfer(
        string memory _validator,
        address _to,
        uint256 _shares
    ) internal pure returns (bytes memory) {
        return
        abi.encodeWithSignature(
            "transfer(string,address,uint256)",
            _validator,
            _to,
            _shares
        );
    }

    function transferFrom(
        string memory _validator,
        address _from,
        address _to,
        uint256 _shares
    ) internal pure returns (bytes memory) {
        return
        abi.encodeWithSignature(
            "transferFrom(string,address,address,uint256)",
            _validator,
            _from,
            _to,
            _shares
        );
    }

    function delegation(
        string memory _validator,
        address _delegate
    ) internal pure returns (bytes memory) {
        return
        abi.encodeWithSignature(
            "delegation(string,address)",
            _validator,
            _delegate
        );
    }

    function delegationRewards(
        string memory _validator,
        address _delegate
    ) internal pure returns (bytes memory) {
        return
        abi.encodeWithSignature(
            "delegationRewards(string,address)",
            _validator,
            _delegate
        );
    }

    function allowance(
        string memory _validator,
        address _owner,
        address _spender
    ) internal pure returns (bytes memory) {
        return
        abi.encodeWithSignature(
            "allowance(string,address,address)",
            _validator,
            _owner,
            _spender
        );
    }
}
