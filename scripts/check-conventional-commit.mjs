import { readFileSync } from "node:fs";

const file = process.argv[2];
const message = readFileSync(file, "utf8").split("\n")[0];
const allowed =
  /^(feat|fix|docs|chore|refactor|test|ops|data)(\([^)]+\))?: .{1,120}$/;

if (!allowed.test(message)) {
  console.error(`Invalid Conventional Commit message: ${message}`);
  process.exit(1);
}
