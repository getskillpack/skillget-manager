# Участие в разработке `getskillpack/skillget-manager`

Спасибо за интерес к библиотеке ядра менеджера пакетов скиллов — её использует официальный CLI [`getskillpack/cli`](https://github.com/getskillpack/cli). Ниже — минимум, чтобы собрать проект и открыть осмысленный PR; общий контекст экосистемы — в ссылках, без длинных дубликатов.

## Что почитать в первую очередь

- **Сборка CLI, переменные окружения, пользовательские сценарии:** [docs/QUICKSTART.md](https://github.com/getskillpack/cli/blob/main/docs/QUICKSTART.md)
- **Матрица совместимости** (CLI, skillget-manager, registry): [docs/COMPATIBILITY_MATRIX_RU.md](https://github.com/getskillpack/cli/blob/main/docs/COMPATIBILITY_MATRIX_RU.md)
- **Индекс документации CLI:** [docs/README.md](https://github.com/getskillpack/cli/blob/main/docs/README.md)
- **Обзор репозитория и ссылки на смежные репо:** [README.md](README.md)
- **Релизы и semver:** [RELEASE.md](RELEASE.md) и [CHANGELOG.md](CHANGELOG.md)
- **Поведение этой библиотеки как клиента реестра:** [docs/REGISTRY_CLIENT_CONTRACT.md](docs/REGISTRY_CLIENT_CONTRACT.md)

## Требования к окружению

- **Go 1.22+** (см. `go.mod` / `toolchain` в корне).

## Сборка и проверка из исходников

Из корня репозитория:

```bash
go test ./... -count=1
```

Это совпадает с шагом **Test** в [.github/workflows/go.yml](.github/workflows/go.yml). Тесты поднимают локальные `httptest`-серверы там, где нужен ответ реестра; отдельно поднимать registry для `go test` обычно не требуется.

### Приватный модуль org getskillpack

Для `go test` при необходимости доступа к приватным модулям GitHub org настройте `GOPRIVATE`, `GONOSUMDB` и учётные данные git — чеклист: [BOARD_PAT_QUICK_RU.md](https://github.com/getskillpack/cli/blob/main/docs/BOARD_PAT_QUICK_RU.md) (в т.ч. § про Go modules). **Не коммитьте** токены и PAT в репозиторий.

### Связка с CLI при локальной разработке

`getskillpack/cli` пинит этот модуль по **semver** в `go.mod`. Для локального цикла можно временно добавить директиву `replace` на соседний checkout менеджера. Когда модуль публично доступен через Go proxy, `GOPRIVATE` для этого пути не нужен.

## Issues и pull requests

1. **Issue** — опишите воспроизведение, ожидаемое и фактическое поведение, задействованный тег модуля (`v…`) и при необходимости `SKILLGET_REGISTRY_URL` / фрагмент лога.
2. **PR** — нацеливайте на ветку **`main`**, одна логическая тема PR; для пользовательски заметных изменений обновляйте [CHANGELOG.md](CHANGELOG.md) (см. [RELEASE.md](RELEASE.md)).
3. Агентам Paperclip и автоматизации в org **getskillpack**: гигиена PAT, ветки и push — [AGENT_GITHUB_REPO_WORKFLOW_RU.md](https://github.com/getskillpack/cli/blob/main/docs/AGENT_GITHUB_REPO_WORKFLOW_RU.md).

Контракт HTTP API реестра — в [`getskillpack/registry`](https://github.com/getskillpack/registry); пользовательский CLI — в [`getskillpack/cli`](https://github.com/getskillpack/cli).

## Безопасность

Сообщения об уязвимостях — по [SECURITY.md](SECURITY.md).
