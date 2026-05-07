import { execSync } from "node:child_process";
import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";

const target = "src/generated/version.ts";

if (process.env.UPDATE_VERSION !== "1" && existsSync(target)) {
  process.exit(0);
}

const packageJson = JSON.parse(readFileSync("package.json", "utf8"));

function run(command, fallback) {
  try {
    return execSync(command, { stdio: ["ignore", "pipe", "ignore"] })
      .toString()
      .trim();
  } catch {
    return fallback;
  }
}

const version = packageJson.version;
const commit = run("git rev-parse --short HEAD", "dev");
const builtAt = run("git show -s --format=%cI HEAD", new Date().toISOString());

mkdirSync("src/generated", { recursive: true });
writeFileSync(
  target,
  `export const APP_VERSION = ${JSON.stringify(version)}\n` +
    `export const APP_COMMIT = ${JSON.stringify(commit)}\n` +
    `export const APP_BUILT_AT = ${JSON.stringify(builtAt)}\n`,
);
