// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iota/iscmagic/ISC.sol";

contract Random {
    event Bytes32Event(bytes32 value);

    function Bytes32() public returns (bytes32) {
        bytes32 value = ISC.sandbox.getEntropy();
        emit Bytes32Event(value);
        return value;
    }
}