# skillget-manager

Библиотека ядра менеджера пакетов скиллов для org [getskillpack](https://github.com/getskillpack): клиент реестра, `skills.lock`, поиск и загрузка архивов. Реализация на **Go 1.22+** (модуль `github.com/getskillpack/skillget-manager`).

## Зачем отдельный репозиторий

CLI ([cli](https://github.com/getskillpack/cli)) остаётся тонкой оболочкой (`skillget`). Вся логика установки и контракт с реестром живёт здесь, чтобы board и интеграции смотрели на **менеджер** и **registry**, а не только на CLI.

## API (кратко)

- `RegistryBaseURL` / `RegistryConfigSource` — выбор базы API из env (`SKILLGET_REGISTRY_URL`, legacy `SKPKG_REGISTRY_URL`).
- `RegistryToken` — bearer для записи в реестр из `SKILLGET_REGISTRY_TOKEN` или `SKILLGET_TOKEN`.
- `FetchJSON` — JSON GET к реестру (через `context.Context`).
- `ReadSkillsLock` / `WriteSkillsLock` — файл `skills.lock`.
- `SearchSkills` — список/поиск скиллов (опции `Query`, `Author`, `Limit`, `Offset`).
- `ParseNameVersion` / `ResolveInstallTarget` — разрешение версии и метаданные архива (для «latest» сначала поле `latest_version` в ответе реестра, иначе первая неснятая версия в списке).
- `DownloadSkillArchive` — скачивание tarball, проверка `checksum` вида `sha256:<hex>` при наличии, обновление lockfile.
- `PublishSkill` — multipart POST `/skills` (manifest + archive), как у TS-клиента в `skpkg-cli`.

Контракт HTTP API реестра: [registry/API.md](https://github.com/getskillpack/registry/blob/main/API.md).

## Разработка

```bash
go test ./...
```

Подробнее: [CONTRIBUTING.md](CONTRIBUTING.md) (интеграционные тесты, связка с CLI, релизы).

## Лицензия

MIT — см. [LICENSE](LICENSE).

## Безопасность

Сообщения об уязвимостях: [SECURITY.md](SECURITY.md).
