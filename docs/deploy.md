# GitHub Pages Deploy

Live site:

https://baditaflorin.github.io/accident-reconstructor/

GitHub Pages source:

main branch, `/docs` folder

Repository settings:

https://github.com/baditaflorin/accident-reconstructor/settings/pages

## Publish

```sh
make build
git add docs
git commit -m "ops: publish pages build"
git push
```

## Rollback

Revert the publishing commit and push:

```sh
git revert COMMIT_SHA
git push
```

## Pages Gotchas

- Vite base path is `/accident-reconstructor/`.
- `docs/404.html` is copied from `docs/index.html` for SPA fallback.
- GitHub Pages does not support `_headers` or `_redirects`.
- Service worker scope is `/accident-reconstructor/`.
- The backend is not served by Pages.
