// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iota/iscmagic/ISC.sol";

contract Random {
    event Byte(uint8);
    event UInt64(uint64);
    event Bytes(bytes);

    function getByte() public returns (uint8) {
        bytes32 entropy = ISC.sandbox.getEntropy();
        uint8 value = uint8(entropy[0]);
        emit Byte(value);
        return value;
    }

    function getUInt64() public returns (uint64) {
        bytes32 entropy = ISC.sandbox.getEntropy();
        uint64 value = uint64(bytes8(entropy));
        emit UInt64(value);
        return value;
    }

    function getBytes(uint n) public returns (bytes memory) {
        bytes memory value = new bytes(n);
        for (uint i = 0; i < n; i += 32) {
            bytes32 entropy = ISC.sandbox.getEntropy();
            for (uint j = 0; j < 32 && i + j < n; j++) {
                value[i + j] = entropy[j];
            }
        }
        emit Bytes(value);
        return value;
    }
}
