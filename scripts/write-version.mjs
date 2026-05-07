import { execSync } from 'node:child_process'
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs'

const packageJson = JSON.parse(readFileSync('package.json', 'utf8'))

function run(command, fallback) {
  try {
    return execSync(command, { stdio: ['ignore', 'pipe', 'ignore'] }).toString().trim()
  } catch {
    return fallback
  }
}

const version = packageJson.version
const commit = run('git rev-parse --short HEAD', 'dev')
const builtAt = new Date().toISOString()

mkdirSync('src/generated', { recursive: true })
writeFileSync(
  'src/generated/version.ts',
  `export const APP_VERSION = ${JSON.stringify(version)}\n` +
    `export const APP_COMMIT = ${JSON.stringify(commit)}\n` +
    `export const APP_BUILT_AT = ${JSON.stringify(builtAt)}\n`,
)
