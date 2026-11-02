# dictionary-service

Reference data and localisation for the Markerloop platform: UI translations, localised error messages, supported
languages, countries and cities. The admin UI, customer sites and the other services read their texts from here,
in Russian, Kazakh and English.

Markerloop is a shared backend (users, auth, files, collections, notifications, i18n) plus a page builder that an AI
agent drives through MCP. dictionary-service is the i18n backbone of that product.

## Responsibilities

- **Translations**: batch lookup of translation keys for a language and application (`POST /dictionary/translations`).
- **Localised errors**: resolve an error code to a message in the requested language; list all known errors.
- **Languages**: list of supported languages with their names in each UI language (`ru-RU`, `kk-KZ`, `en-US`).
- **Geography**: countries and cities with pagination, filtering and sorting.
- **Seeding**: translation seeds are embedded into the binary (`go:embed`) and applied at startup with
  `ON CONFLICT DO UPDATE`, so a deploy is also a content update and re-running is safe.

## API surface

Base path `/api/v1`. Full spec: `docs/swagger.yaml` / `docs/swagger.json`, served at `/swagger/index.html`.
All endpoints are public (read-only) and published through the gateway without authentication.

| Method | Path | Purpose |
|---|---|---|
| POST | `/dictionary/translations` | Translations for a set of keys |
| POST | `/dictionary/error` | One localised error by code |
| GET | `/dictionary/errors` | All errors |
| GET | `/dictionary/languages` | Supported languages |
| GET | `/dictionary/countries` | Countries (paging, filter, sort) |
| GET | `/dictionary/cities` | Cities (paging, filter, sort) |
| GET | `/health` | Health (root path, used by the gateway and k8s probes) |

## Architecture notes

- **Layers**: `cmd/server` → `internal/routes` → `internal/handlers` → `internal/services` → `internal/repository`
  (GORM on PostgreSQL: translations, countries, cities). `pkg/dictionary` holds domain constants — error codes,
  language names, messages — importable by other services.
- **Caching**: read-through Redis cache in the service layer (keys per code/language/app, `CACHE_TTL`, default 1 h).
- **Migrations**: idempotent Go migrations in `internal/migrations` (performance indexes, cleanup of retired keys)
  plus the embedded SQL seeds in `internal/migrations/seeds`. The root-level `seed_*.sql` files are copies of those
  seeds for manual use.

## Running locally

```bash
cp .env.example .env
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=<local-db-password> -e POSTGRES_DB=dictionary postgres:16-alpine
docker run -d -p 6379:6379 redis:7
go run ./cmd/server
curl localhost:8084/health
```

Live reload: `air`. Docker:

```bash
docker build --build-arg DEPLOY_TOKEN_USERNAME=<user> --build-arg DEPLOY_TOKEN_PASSWORD=<token> -t dictionary-service .
docker run -p 8084:8084 --env-file .env dictionary-service
```

## Tests

```bash
go test ./...
```

Covers handlers, the route table contract, seed application and key cleanup, and the error
catalogue in `pkg/dictionary`.

## Notes for reviewers

- Code comments and Swagger descriptions are largely in Russian.
- The shared library is imported as `gitlab.com/markerloop/go-sdk` (private GitLab) and is published as
  [Markerloop-Team/go-sdk](https://github.com/Markerloop-Team/go-sdk); to build outside our CI, clone it and add
  `replace gitlab.com/markerloop/go-sdk => ../go-sdk` to `go.mod`.
- Snapshot without git history. Ticket ids in comments refer to our internal tracker.

## Part of Markerloop

Markerloop is a shared backend and site builder with an AI agent working through MCP.
Website: [markerloop.com](https://markerloop.com) · Code: [github.com/Markerloop-Team](https://github.com/Markerloop-Team)

---

Source published for review. © Markerloop. All rights reserved.
