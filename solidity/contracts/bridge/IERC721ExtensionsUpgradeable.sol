// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

interface IERC721ExtensionsUpgradeable {
    /**
     * @dev Destroys `id` tokens from the caller.
     */
    function burn(uint256 id) external;

    /** @dev Creates `id` tokens and assigns them to `account`, increasing
     * the total supply.
     */
    function mint(address account, uint256 id) external;
}
