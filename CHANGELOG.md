# Changelog

## 0.1.3 — 2026-03-28

- `DownloadSkillArchive`: подсказки при сетевых сбоях и при ответах **401 / 403 / 404 / 429 / 503** от хоста архива (URL может отличаться от базы реестра).
- [CONTRIBUTING.md](CONTRIBUTING.md): ссылка на контракт реестра выровнена с README (`API.md`).

## 0.1.2 — 2026-03-28

- Подсказка при ответе реестра **HTTP 429** (rate limit) для `FetchJSON` и `PublishSkill`.

## 0.1.1 — 2026-03-28

- Тесты `DownloadSkillArchive`: `archive_url` в моках указывает на `httptest` по HTTP (раньше был `https://` + `127.0.0.1`, клиент падал с «HTTP response to HTTPS client»).

## 0.1.0 — 2026-03-28

Первый semver-тег библиотеки для пинов в Go-модулях и релизной дисциплины.

- `ResolveInstallTarget`: если `latest_version` совпадает с записью в `versions[]` и она **yanked**, выбирается следующая установимая версия (как при отсутствии `latest_version`).
- Добавлен [CONTRIBUTING.md](CONTRIBUTING.md): локальный прогон тестов, интеграция с CLI/registry, релизная дисциплина.
- Добавлен [RELEASE.md](RELEASE.md): порядок тегирования и синхронизации с CLI.
- `RegistryToken`, `PublishSkill` (multipart publish в реестр).
- `SearchSkills`: фильтр `Author` → query `author`.
- `ResolveInstallTarget`: приоритет поля `latest_version` из GET `/skills/:name` (если не yanked в списке версий).
- `DownloadSkillArchive`: проверка SHA-256 при `checksum` в формате `sha256:<64 hex>`.
- Сообщения об ошибках HTTP к реестру с краткими подсказками (как в TS `registry.ts`).
- Документация: `SECURITY.md`, ссылки на лицензию и процесс безопасности в README.

## Unreleased

