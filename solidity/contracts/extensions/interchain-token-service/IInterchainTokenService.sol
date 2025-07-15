// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

interface IInterchainTokenService {
    /**
     * @notice Returns the custom tokenId associated with the given operator and salt.
     * @param operator_ The operator address.
     * @param salt The salt used for token id calculation.
     * @return tokenId The custom tokenId associated with the operator and salt.
     */
    function interchainTokenId(address operator_, bytes32 salt) external view returns (bytes32 tokenId);
}
