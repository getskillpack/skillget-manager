# skillget-manager

Библиотека ядра менеджера пакетов скиллов для org [getskillpack](https://github.com/getskillpack): клиент реестра, `skills.lock`, поиск и загрузка архивов.

## Зачем отдельный репозиторий

CLI ([cli](https://github.com/getskillpack/cli)) остаётся тонкой оболочкой (`skillget`). Вся логика установки и контракт с реестром живёт здесь, чтобы board и интеграции смотрели на **менеджер** и **registry**, а не только на CLI.

## API (кратко)

- `registryBaseUrl` / `registryConfigSource` — выбор базы API из env (`SKILLGET_REGISTRY_URL`, legacy `SKPKG_REGISTRY_URL`).
- `fetchJson` — JSON-запросы к реестру.
- `readSkillsLock` / `writeSkillsLock` — файл `skills.lock`.
- `searchSkills` — список/поиск скиллов.
- `resolveInstallTarget` — разрешение версии и метаданные архива.
- `downloadSkillArchive` — скачивание tarball и обновление lockfile.

Контракт HTTP API реестра: [registry/API.md](https://github.com/getskillpack/registry/blob/main/API.md).

## Разработка

```bash
npm install
npm run build
```

Node 18+.

## Лицензия

MIT
