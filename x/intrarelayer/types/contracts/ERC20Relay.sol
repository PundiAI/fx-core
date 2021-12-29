// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

contract ERC20RelayExternal {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;


    string private _name;
    string private _symbol;
    uint8 private _decimals;
    uint256 private _totalSupply;

    address private _owner;
    address private _intrarelayer;

    constructor(string memory name_, string memory symbol_, uint8 decimals_, address intrarelayer_) {
        _name = name_;
        _symbol = symbol_;
        _decimals = decimals_;
        _intrarelayer = intrarelayer_;
        _owner = msg.sender;
    }

    function name() public view returns (string memory) {
        return _name;
    }

    function symbol() public view returns (string memory) {
        return _symbol;
    }

    function decimals() public view returns (uint8) {
        return _decimals;
    }

    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    function transfer(address recipient, uint256 amount) public returns (bool) {
        _transfer(msg.sender, recipient, amount);
        return true;
    }

    function allowance(address account, address spender) public view returns (uint256) {
        return _allowances[account][spender];
    }

    function approve(address spender, uint256 amount) public returns (bool) {
        _approve(msg.sender, spender, amount);
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(address sender, address recipient, uint256 amount) public returns (bool) {
        uint256 currentAllowance = _allowances[sender][msg.sender];
        require(currentAllowance >= amount, "transfer amount exceeds allowance");
        _approve(sender, msg.sender, currentAllowance - amount);
        _transfer(sender, recipient, amount);
        return true;
    }

    function mint(address account, uint256 amount) public onlyOwner {
        _mint(account, amount);
    }

    function burn(uint256 amount) public {
        _burn(msg.sender, amount);
    }

    function relay(address recipient, uint256 amount) public returns (bool){
        _relay(msg.sender, recipient, amount);
        return true;
    }

    function relayFrom(address sender, address recipient, uint256 amount) public returns (bool){
        uint256 currentAllowance = _allowances[sender][msg.sender];
        require(currentAllowance >= amount, "relay amount exceeds allowance");
        _approve(sender, msg.sender, currentAllowance - amount);
        _relay(sender, recipient, amount);
        return true;
    }

    function intrarelayer() public view returns (address){
        return _intrarelayer;
    }

    function owner() public view virtual returns (address) {
        return _owner;
    }
    modifier onlyOwner() {
        require(owner() == msg.sender, "caller is not the owner");
        _;
    }

    function _transfer(address sender, address recipient, uint256 amount) internal {
        require(sender != address(0), "transfer from the zero address");
        require(recipient != address(0), "transfer to the zero address");
        uint256 senderBalance = _balances[sender];
        require(senderBalance >= amount, "transfer amount exceeds balance");
        _balances[sender] = senderBalance - amount;
        _balances[recipient] += amount;

        emit Transfer(sender, recipient, amount);
    }

    function _mint(address account, uint256 amount) internal {
        require(account != address(0), "mint to the zero address");
        _totalSupply += amount;
        _balances[account] += amount;

        emit Transfer(address(0), account, amount);
    }

    function _burn(address account, uint256 amount) internal {
        require(account != address(0), "burn from the zero address");
        uint256 accountBalance = _balances[account];
        require(accountBalance >= amount, "burn amount exceeds balance");
        _balances[account] = accountBalance - amount;
        _totalSupply -= amount;

        emit Transfer(account, address(0), amount);
    }

    function _approve(address sender, address spender, uint256 amount) internal {
        require(sender != address(0), "approve from the zero address");
        _allowances[sender][spender] = amount;
    }

    function _relay(address sender, address recipient, uint256 amount) internal {
        require(sender != address(0), "relay from the zero address");
        require(recipient != address(0), "relay to the zero address");

        _transfer(sender, _intrarelayer, amount);

        emit Relay(sender, recipient, amount);
    }

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Relay(address indexed from, address indexed to, uint256 value);
}