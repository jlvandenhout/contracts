import { expect } from "chai";
import { ethers } from "hardhat";

describe("Random", function () {
  describe("Events", function () {
    let value = new Uint8Array(32);
    it("Should emit random bytes32", async function () {
      const Random = await ethers.getContractFactory("Random");
      let owner = await ethers.getSigner(
        "0xb62ea087c36eBd77Fb58174Ae162395722dE9Cb0"
      );

      const random = await Random.deploy({ from: owner.address });
      await expect(random.Bytes32())
        .to.emit(random, "Bytes32Event")
        .withArgs(value);
    });
  });
});
