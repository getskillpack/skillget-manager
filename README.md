# skillget-manager

Библиотека ядра менеджера пакетов скиллов для org [getskillpack](https://github.com/getskillpack): клиент реестра, `skills.lock`, поиск и загрузка архивов.

**Стек:** Go 1.22+. Прототип на TypeScript **снят** (март 2026, согласование board); ядро только на Go.

## Зачем отдельный репозиторий

CLI ([cli](https://github.com/getskillpack/cli)) остаётся тонкой оболочкой (`skillget`). Вся логика установки и контракт с реестром живёт здесь.

## Контракт реестра

Черновик HTTP API: [registry/API.md](https://github.com/getskillpack/registry/blob/main/API.md).

## Разработка

```bash
go test ./...
```

## Лицензия

MIT
