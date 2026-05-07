#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

mkdir -p tmp
npm run build
CGO_ENABLED=0 go build -o tmp/accident-server ./cmd/server

rm -rf tmp/pages-preview
mkdir -p tmp/pages-preview/accident-reconstructor
cp -R docs/. tmp/pages-preview/accident-reconstructor/

PAGE_PORT=41739
ADDR=:18080 STORAGE_DIR=./tmp/smoke-cases PAGES_ORIGIN=http://127.0.0.1:${PAGE_PORT} ./tmp/accident-server &
API_PID=$!
npx http-server tmp/pages-preview -a 127.0.0.1 -p ${PAGE_PORT} -c-1 >/tmp/accident-pages.log 2>&1 &
WEB_PID=$!
trap 'kill "$API_PID" "$WEB_PID" >/dev/null 2>&1 || true' EXIT

for _ in {1..40}; do
  if curl -fsS http://127.0.0.1:18080/healthz >/dev/null; then
    break
  fi
  sleep 0.25
done

curl -fsS http://127.0.0.1:18080/readyz >/dev/null
curl -fsS http://127.0.0.1:18080/metrics >/dev/null

printf 'synthetic smoke video placeholder\n' >tmp/smoke.mp4
CREATE_RESPONSE="$(curl -fsS -F case_name=Smoke -F scale_meters=8 -F videos=@tmp/smoke.mp4 http://127.0.0.1:18080/api/v1/cases)"
CASE_ID="$(node -e "const fs=require('fs'); const data=JSON.parse(fs.readFileSync(0,'utf8')); console.log(data.case.id)" <<<"$CREATE_RESPONSE")"

STATUS=processing
for _ in {1..40}; do
  CASE_RESPONSE="$(curl -fsS "http://127.0.0.1:18080/api/v1/cases/$CASE_ID")"
  STATUS="$(node -e "const fs=require('fs'); const data=JSON.parse(fs.readFileSync(0,'utf8')); console.log(data.status)" <<<"$CASE_RESPONSE")"
  if [[ "$STATUS" == "complete" ]]; then
    break
  fi
  sleep 0.25
done

test "$STATUS" = "complete"
curl -fsS "http://127.0.0.1:18080/api/v1/cases/$CASE_ID/artifact" |
  node -e "const fs=require('fs'); const data=JSON.parse(fs.readFileSync(0,'utf8')); if (!data.vehicleTrack?.length) process.exit(1)"

for _ in {1..40}; do
  if curl -fsS http://127.0.0.1:${PAGE_PORT}/accident-reconstructor/ >/dev/null; then
    break
  fi
  sleep 0.25
done

PLAYWRIGHT_BASE_URL=http://127.0.0.1:${PAGE_PORT}/accident-reconstructor/ npm run smoke:web
