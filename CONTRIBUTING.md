# Contributing

Thanks for helping improve Accident Reconstructor.

## Local Setup

```sh
npm install
make install-hooks
make test
make smoke
```

Use Conventional Commits for every commit, for example `feat: add case upload flow`.

## Pull Request Checklist

- Run `make fmt`, `make lint`, `make test`, and `make smoke`.
- Do not commit secrets, private footage, or case evidence.
- Update ADRs before major architecture changes.
- Keep legal and safety claims conservative.

## Evidence Handling

Never upload real accident evidence to a public issue. Use synthetic or redacted sample data only.
