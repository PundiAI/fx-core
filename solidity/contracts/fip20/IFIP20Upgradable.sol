// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

interface IFIP20Upgradable {
    /**
     * @dev Returns the name of the token.
     */
    function name() external view returns (string memory);

    /**
     * @dev Returns the symbol of the token.
     */
    function symbol() external view returns (string memory);

    /**
     * @dev Returns the decimals places of the token.
     */
    function decimals() external view returns (uint8);

    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `to`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address to, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(
        address owner,
        address spender
    ) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `from` to `to` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);

    /**
     * @dev Mint `amount` tokens to `to`, increase the
     * total supply.
     *
     * Emits a {Transfer} event with `from` set to the zero address.
     */
    function mint(address to, uint256 amount) external;

    /**
     * @dev Burn `amount` tokens from `from`, reducing the
     * total supply.
     *
     * Emits a {Transfer} event with `to` set to the zero address.
     */
    function burn(address from, uint256 amount) external;

    /**
     * @dev Cross chain moves `amount + fee` tokens from `sender` to `recipient`
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {TransferCrossChain} event.
     */
    // Deprecated: use pre-compiled contract crossChain
    function transferCrossChain(
        string memory recipient,
        uint256 amount,
        uint256 fee,
        bytes32 target
    ) external returns (bool);

    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );

    /**
     * @dev Emitted when `amount + fee` tokens are cross chain moved from one account (`from`) to
     * another (`to`) through (`target`).
     */
    event TransferCrossChain(
        address indexed from,
        string recipient,
        uint256 amount,
        uint256 fee,
        bytes32 target
    );
}
