// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

/* solhint-disable custom-errors */

contract DataCallbackTest {
    address public admin;
    uint256 public id;

    constructor(address _admin) {
        admin = _admin;
        id = 0;
    }

    function setID(uint256 _id) public onlyAdmin {
        id = _id;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "only Admin");
        _;
    }
}
