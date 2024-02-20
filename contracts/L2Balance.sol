// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iota/iscmagic/ISC.sol";

contract L2Balance {
    function getBalance(
        NativeTokenID memory nativeTokenID,
        ISCAgentID memory agentID
    ) public view returns (uint256) {
        return ISC.accounts.getL2BalanceNativeTokens(nativeTokenID, agentID);
    }
}
