# skillget-manager

Библиотека ядра менеджера пакетов скиллов для org [getskillpack](https://github.com/getskillpack): клиент реестра, `skills.lock`, поиск и загрузка архивов. Реализация на **Go 1.22+** (модуль `github.com/getskillpack/skillget-manager`).

## Зачем отдельный репозиторий

CLI ([cli](https://github.com/getskillpack/cli)) остаётся тонкой оболочкой (`skillget`). Вся логика установки и контракт с реестром живёт здесь, чтобы board и интеграции смотрели на **менеджер** и **registry**, а не только на CLI.

## API (кратко)

- `RegistryBaseURL` / `RegistryConfigSource` — выбор базы API из env (`SKILLGET_REGISTRY_URL`, legacy `SKPKG_REGISTRY_URL`).
- `RegistryToken` — bearer для записи в реестр из `SKILLGET_REGISTRY_TOKEN` или `SKILLGET_TOKEN`.
- `RegistryReadBearer` — bearer для чтения (`SKILLGET_REGISTRY_READ_TOKEN` или, если пусто, тот же токен, что и для записи): добавляется к `GET` API и к загрузке `archive_url`, когда токен задан.
- `FetchJSON` — JSON GET к реестру (через `context.Context`); типичные ответы реестра сопровождаются короткими подсказками (сеть, токен, 404, 409, 429 и др.).
- `ReadSkillsLock` / `WriteSkillsLock` — файл `skills.lock`.
- `SearchSkills` — список/поиск скиллов (опции `Query`, `Author`, `Limit`, `Offset`).
- `ParseNameVersion` / `ResolveInstallTarget` — разрешение версии и метаданные архива; для «latest» без pin — **наибольшая semver** среди записей `versions` в ответе `GET /skills/{name}`, исключая `yanked: true` (как в reference registry).
- `DownloadSkillArchive` — скачивание tarball, проверка `checksum` вида `sha256:<hex>` при наличии, обновление lockfile; при сбоях HTTP к **archive URL** (часто отдельно от базы реестра) — свои подсказки (401–429, 503, сеть).
- `PublishSkill` — multipart POST `/skills` (manifest + archive), как у TS-клиента в `skpkg-cli`.

Контракт HTTP API реестра: [registry/API.md](https://github.com/getskillpack/registry/blob/main/API.md). Поведение этой библиотеки как клиента: [docs/REGISTRY_CLIENT_CONTRACT.md](docs/REGISTRY_CLIENT_CONTRACT.md).

## Разработка

```bash
go test ./...
```

Подробнее: [CONTRIBUTING.md](CONTRIBUTING.md) (интеграционные тесты, связка с CLI, релизы).

## Лицензия

MIT — см. [LICENSE](LICENSE).

## Безопасность

Сообщения об уязвимостях: [SECURITY.md](SECURITY.md).
