import { rmSync } from "node:fs";

for (const path of ["docs/assets", "docs/404.html"]) {
  rmSync(path, { recursive: true, force: true });
}
