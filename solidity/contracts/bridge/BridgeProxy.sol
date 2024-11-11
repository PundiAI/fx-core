// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ERC1967Upgrade} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Upgrade.sol";
import {Proxy} from "@openzeppelin/contracts/proxy/Proxy.sol";

error AlreadyInitialized();

contract BridgeProxy is Proxy, ERC1967Upgrade {
    function init(address _logic) public {
        if (_getImplementation() != address(0)) {
            revert AlreadyInitialized();
        }
        _upgradeTo(_logic);
    }

    /**
     * @dev Returns the current implementation address.
     */
    function _implementation()
        internal
        view
        virtual
        override
        returns (address impl)
    {
        return _getImplementation();
    }
}
