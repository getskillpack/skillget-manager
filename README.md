# skillget-manager

Библиотека ядра менеджера пакетов скиллов для org [getskillpack](https://github.com/getskillpack): клиент реестра, `skills.lock`, поиск и загрузка архивов. Реализация на **Go 1.22+** (модуль `github.com/getskillpack/skillget-manager`).

## Зачем отдельный репозиторий

CLI ([cli](https://github.com/getskillpack/cli)) остаётся тонкой оболочкой (`skillget`). Вся логика установки и контракт с реестром живёт здесь, чтобы board и интеграции смотрели на **менеджер** и **registry**, а не только на CLI.

## API (кратко)

- `RegistryBaseURL` / `RegistryConfigSource` — выбор базы API из env (`SKILLGET_REGISTRY_URL`, legacy `SKPKG_REGISTRY_URL`).
- `FetchJSON` — JSON GET к реестру (через `context.Context`).
- `ReadSkillsLock` / `WriteSkillsLock` — файл `skills.lock`.
- `SearchSkills` — список/поиск скиллов.
- `ParseNameVersion` / `ResolveInstallTarget` — разрешение версии и метаданные архива.
- `DownloadSkillArchive` — скачивание tarball и обновление lockfile.

Контракт HTTP API реестра: [registry/API.md](https://github.com/getskillpack/registry/blob/main/API.md).

## Разработка

```bash
go test ./...
```

## Лицензия

MIT
