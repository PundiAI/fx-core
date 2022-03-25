// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

contract FIP20 {
    string public name;
    string public symbol;
    uint8 public decimals;
    uint256 public totalSupply;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    address public owner;

    constructor(string memory name_, string memory symbol_, uint8 decimals_) {
        name = name_;
        symbol = symbol_;
        decimals = decimals_;
        owner = msg.sender;
    }

    function approve(address spender, uint256 amount) public returns (bool) {
        _approve(msg.sender, spender, amount);
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transfer(address recipient, uint256 amount) public returns (bool) {
        _transfer(msg.sender, recipient, amount);
        return true;
    }

    function transferFrom(address sender, address recipient, uint256 amount) public returns (bool) {
        uint256 currentAllowance = allowance[sender][msg.sender];
        require(currentAllowance >= amount, "transfer amount exceeds allowance");
        _approve(sender, msg.sender, currentAllowance - amount);
        _transfer(sender, recipient, amount);
        return true;
    }

    function mint(address account, uint256 amount) public onlyOwner {
        _mint(account, amount);
    }

    function burn(address account, uint256 amount) public onlyOwner {
        _burn(account, amount);
    }

    function transferCross(string memory to, uint256 amount, uint256 fee, string memory target) public notContract returns (bool) {
        _transferCross(msg.sender, to, amount, fee, target);
        return true;
    }


    function module() public view returns (address){
        return owner;
    }

    function isContract(address _addr) public view returns (bool){
        uint32 size;
        assembly {
            size := extcodesize(_addr)
        }
        return (size > 0);
    }

    modifier onlyOwner() {
        require(owner == msg.sender, "caller is not the owner");
        _;
    }

    modifier notContract(){
        require(!isContract(msg.sender), "caller cannot be contract");
        _;
    }


    function _transfer(address sender, address recipient, uint256 amount) internal {
        require(sender != address(0), "transfer from the zero address");
        require(recipient != address(0), "transfer to the zero address");
        uint256 senderBalance = balanceOf[sender];
        require(senderBalance >= amount, "transfer amount exceeds balance");
        balanceOf[sender] = senderBalance - amount;
        balanceOf[recipient] += amount;

        emit Transfer(sender, recipient, amount);
    }

    function _mint(address account, uint256 amount) internal {
        require(account != address(0), "mint to the zero address");
        totalSupply += amount;
        balanceOf[account] += amount;

        emit Transfer(address(0), account, amount);
    }

    function _burn(address account, uint256 amount) internal {
        require(account != address(0), "burn from the zero address");
        uint256 accountBalance = balanceOf[account];
        require(accountBalance >= amount, "burn amount exceeds balance");
        balanceOf[account] = accountBalance - amount;
        totalSupply -= amount;

        emit Transfer(account, address(0), amount);
    }

    function _approve(address sender, address spender, uint256 amount) internal {
        require(sender != address(0), "approve from the zero address");
        allowance[sender][spender] = amount;
    }

    function _transferCross(address from, string memory to, uint256 amount, uint256 fee, string memory target) internal {
        require(from != address(0), "transfer from zero address");
        require(bytes(to).length > 0, "transfer to the empty");
        require(bytes(target).length > 0, "empty target");

        _transfer(from, module(), amount + fee);
        emit TransferCross(from, to, amount, fee, target);
    }

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    event TransferCross(address indexed from, string to, uint256 value, uint256 fee, string target);
}