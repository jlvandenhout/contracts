import { ethers } from "hardhat";

async function main() {
  const random = await ethers.deployContract("Random");
  await random.waitForDeployment();

  console.log(`Random deployed to ${random.target}`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
