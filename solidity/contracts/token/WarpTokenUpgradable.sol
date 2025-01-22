// SPDX-License-Identifier: Apache-2.0
/* solhint-disable one-contract-per-file */

pragma solidity ^0.8.0;

/* solhint-disable no-global-import */
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/interfaces/draft-IERC1822Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/ERC1967/ERC1967UpgradeUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/IERC20Upgradeable.sol";

/* solhint-enable no-global-import */
/* solhint-disable custom-errors */

/**
 * @dev An upgradeability mechanism designed for UUPS proxies. The functions included here can perform an upgrade of an
 * {ERC1967Proxy}, when this contract is set as the implementation behind such a proxy.
 *
 * A security mechanism ensures that an upgrade does not turn off upgradeability accidentally, although this risk is
 * reinstated if the upgrade retains upgradeability but removes the security mechanism, e.g. by replacing
 * `UUPSUpgradeable` with a custom implementation of upgrades.
 *
 * The {_authorizeUpgrade} function must be overridden to include access restriction to the upgrade mechanism.
 *
 * _Available since v4.1._
 */
abstract contract UUPSUpgradeable is
    Initializable,
    IERC1822ProxiableUpgradeable,
    ERC1967UpgradeUpgradeable
{
    // solhint-disable-next-line func-name-mixedcase, no-empty-blocks
    function __UUPSUpgradeable_init() internal onlyInitializing {}

    // solhint-disable-next-line func-name-mixedcase, no-empty-blocks
    function __UUPSUpgradeable_init_unchained() internal onlyInitializing {}

    /// @custom:oz-upgrades-unsafe-allow state-variable-immutable state-variable-assignment
    // solhint-disable-next-line immutable-vars-naming
    address private immutable __self =
        0x0000000000000000000000000000000000001002;

    /**
     * @dev Check that the execution is being performed through a delegatecall call and that the execution context is
     * a proxy contract with an implementation (as defined in ERC1967) pointing to self. This should only be the case
     * for UUPS and transparent proxies that are using the current contract as their implementation. Execution of a
     * function through ERC1167 minimal proxies (clones) would not normally pass this test, but is not guaranteed to
     * fail.
     */
    modifier onlyProxy() {
        require(
            address(this) != __self,
            "Function must be called through delegatecall"
        );
        require(
            _getImplementation() == __self,
            "Function must be called through active proxy"
        );
        _;
    }

    /**
     * @dev Check that the execution is not being performed through a delegate call. This allows a function to be
     * callable on the implementing contract but not through proxies.
     */
    modifier notDelegated() {
        require(
            address(this) == __self,
            "UUPSUpgradeable: must not be called through delegatecall"
        );
        _;
    }

    /**
     * @dev Implementation of the ERC1822 {proxiableUUID} function. This returns the storage slot used by the
     * implementation. It is used to validate the implementation's compatibility when performing an upgrade.
     *
     * IMPORTANT: A proxy pointing at a proxiable contract should not be considered proxiable itself, because this risks
     * bricking a proxy that upgrades to it, by delegating to itself until out of gas. Thus it is critical that this
     * function revert if invoked through a proxy. This is guaranteed by the `notDelegated` modifier.
     */
    function proxiableUUID()
        external
        view
        virtual
        override
        notDelegated
        returns (bytes32)
    {
        return _IMPLEMENTATION_SLOT;
    }

    /**
     * @dev Upgrade the implementation of the proxy to `newImplementation`.
     *
     * Calls {_authorizeUpgrade}.
     *
     * Emits an {Upgraded} event.
     */
    function upgradeTo(address newImplementation) external virtual onlyProxy {
        _authorizeUpgrade(newImplementation);
        _upgradeToAndCallUUPS(newImplementation, new bytes(0), false);
    }

    /**
     * @dev Upgrade the implementation of the proxy to `newImplementation`, and subsequently execute the function call
     * encoded in `data`.
     *
     * Calls {_authorizeUpgrade}.
     *
     * Emits an {Upgraded} event.
     */
    function upgradeToAndCall(
        address newImplementation,
        bytes memory data
    ) external payable virtual onlyProxy {
        _authorizeUpgrade(newImplementation);
        _upgradeToAndCallUUPS(newImplementation, data, true);
    }

    /**
     * @dev Function that should revert when `msg.sender` is not authorized to upgrade the contract. Called by
     * {upgradeTo} and {upgradeToAndCall}.
     *
     * Normally, this function will use an xref:access.adoc[access control] modifier such as {Ownable-onlyOwner}.
     *
     * ```solidity
     * function _authorizeUpgrade(address) internal override onlyOwner {}
     * ```
     */
    function _authorizeUpgrade(address newImplementation) internal virtual;

    /**
     * @dev This empty reserved space is put in place to allow future versions to add new
     * variables without shifting down storage in the inheritance chain.
     * See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
     */
    uint256[50] private __gap;
}

contract WarpTokenUpgradable is
    Initializable,
    ContextUpgradeable,
    UUPSUpgradeable,
    OwnableUpgradeable,
    IERC20Upgradeable
{
    string private _name;
    string private _symbol;
    uint8 private _decimals;
    uint256 private _totalSupply;
    mapping(address => uint256) private _balanceOf;
    mapping(address => mapping(address => uint256)) private _allowance;

    address private _module;

    function name() external view returns (string memory) {
        return _name;
    }

    function symbol() external view returns (string memory) {
        return _symbol;
    }

    function decimals() external view returns (uint8) {
        return _decimals;
    }

    function totalSupply() external view returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) external view returns (uint256) {
        return _balanceOf[account];
    }

    function allowance(
        address owner,
        address spender
    ) external view returns (uint256) {
        return _allowance[owner][spender];
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        _approve(_msgSender(), spender, amount);
        emit Approval(_msgSender(), spender, amount);
        return true;
    }

    function transfer(
        address recipient,
        uint256 amount
    ) external returns (bool) {
        _transfer(_msgSender(), recipient, amount);
        return true;
    }

    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) external returns (bool) {
        address spender = _msgSender();
        _spendAllowance(sender, spender, amount);
        _transfer(sender, recipient, amount);
        return true;
    }

    function mint(address account, uint256 amount) external virtual onlyOwner {
        _mint(account, amount);
    }

    /*
     * @notice Only governance can call
     */
    function burn(address account, uint256 amount) external virtual onlyOwner {
        _burn(account, amount);
    }

    function burn(uint256 amount) external virtual {
        _burn(_msgSender(), amount);
    }

    function burnFrom(address account, uint256 amount) external virtual {
        _spendAllowance(account, _msgSender(), amount);
        _burn(account, amount);
    }

    function deposit() public payable {
        _mint(_msgSender(), msg.value);
        emit Deposit(_msgSender(), msg.value);
    }

    function withdraw(address payable to, uint256 value) public {
        _burn(_msgSender(), value);
        to.transfer(value);
        emit Withdraw(_msgSender(), to, value);
    }

    function withdraw(uint256 value) public {
        _burn(_msgSender(), value);
        payable(_msgSender()).transfer(value);
        emit Withdraw(_msgSender(), value);
    }

    event Deposit(address indexed from, uint256 value);
    event Withdraw(address indexed from, address indexed to, uint256 value);
    event Withdraw(address indexed from, uint256 value);

    function module() external view returns (address) {
        return _module;
    }

    modifier notContract() {
        require(!_isContract(_msgSender()), "caller cannot be contract");
        _;
    }

    function _isContract(address _addr) internal view returns (bool) {
        uint32 size;
        // solhint-disable-next-line no-inline-assembly
        assembly {
            size := extcodesize(_addr)
        }
        return (size > 0);
    }

    function _transfer(
        address sender,
        address recipient,
        uint256 amount
    ) internal {
        require(sender != address(0), "transfer from the zero address");
        require(recipient != address(0), "transfer to the zero address");
        uint256 senderBalance = _balanceOf[sender];
        require(senderBalance >= amount, "transfer amount exceeds balance");
        _balanceOf[sender] = senderBalance - amount;
        _balanceOf[recipient] += amount;

        emit Transfer(sender, recipient, amount);
    }

    function _mint(address account, uint256 amount) internal {
        require(account != address(0), "mint to the zero address");
        _totalSupply += amount;
        _balanceOf[account] += amount;

        emit Transfer(address(0), account, amount);
    }

    function _burn(address account, uint256 amount) internal {
        require(account != address(0), "burn from the zero address");
        uint256 accountBalance = _balanceOf[account];
        require(accountBalance >= amount, "burn amount exceeds balance");
        _balanceOf[account] = accountBalance - amount;
        _totalSupply -= amount;

        emit Transfer(account, address(0), amount);
    }

    function _approve(
        address sender,
        address spender,
        uint256 amount
    ) internal {
        require(sender != address(0), "approve from the zero address");
        _allowance[sender][spender] = amount;
    }

    function _spendAllowance(
        address owner,
        address spender,
        uint256 amount
    ) internal {
        uint256 currentAllowance = _allowance[owner][spender];
        if (currentAllowance != type(uint256).max) {
            require(
                currentAllowance >= amount,
                "ERC20: insufficient allowance"
            );
            unchecked {
                _approve(owner, spender, currentAllowance - amount);
            }
        }
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function initialize(
        string memory name_,
        string memory symbol_,
        uint8 decimals_,
        address module_
    ) public virtual initializer {
        _name = name_;
        _symbol = symbol_;
        _decimals = decimals_;
        _module = module_;

        __Ownable_init();
        __UUPSUpgradeable_init();
    }

    receive() external payable {
        deposit();
    }

    fallback() external payable {
        deposit();
    }
}
