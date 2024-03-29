// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
/* solhint-disable custom-errors */

contract TryCatchTestV2 {
    mapping(address => TokenStatus) public tokenStatus;
    bool public initialized;
    uint256 public test;
    bytes32 public placeholder;

    struct TokenStatus {
        bool isOriginated;
        bool isActive;
        bool isExist;
        BridgeTokenType tokenType;
    }

    enum BridgeTokenType {
        ERC20,
        ERC721,
        ERC404
    }

    function initialize(bytes32 _placeholder) public {
        require(!initialized, "Already initialized");
        initialized = true;
        placeholder = _placeholder;
    }

    function test1(uint256 i) internal {
        require(i > 0, "i is not greater than 0");
        test++;
    }

    function test2(uint256 i) internal {
        require(i > 1, "i is not greater than 1");
        test++;
    }

    function tryTest(uint256 i) public {
        test1(i);
        test2(i);
    }

    // Example of try / catch with external call
    // tryCatch(0) => Log("call failed")
    // tryCatch(1) => Log("call failed")
    // tryCatch(2) => Log("my func was called")
    function tryCatch(uint256 _i) public {
        try this.tryTest(_i) {
            emit Log("call success");
        } catch {
            emit Log("call failed");
        }
    }

    function setTokenStatus(
        address token,
        bool isOriginated,
        bool isActive,
        bool isExist,
        BridgeTokenType tokenType
    ) public {
        tokenStatus[token].isOriginated = isOriginated;
        tokenStatus[token].isActive = isActive;
        tokenStatus[token].isExist = isExist;
        tokenStatus[token].tokenType = tokenType;
    }

    function setTokenType(address token, BridgeTokenType typeType) public {
        tokenStatus[token].tokenType = typeType;
    }

    event Log(string message);
}
