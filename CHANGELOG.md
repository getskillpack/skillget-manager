# Changelog

## Unreleased

- Добавлен [CONTRIBUTING.md](CONTRIBUTING.md): локальный прогон тестов, интеграция с CLI/registry, релизная дисциплина.
- Добавлены `RegistryToken`, `PublishSkill` (multipart publish в реестр).
- `SearchSkills`: фильтр `Author` → query `author`.
- `ResolveInstallTarget`: для незакреплённой версии приоритет у поля `latest_version` из GET `/skills/:name`.
- `DownloadSkillArchive`: проверка SHA-256, если реестр отдаёт `checksum` в формате `sha256:<64 hex>`.
- Сообщения об ошибках HTTP к реестру дополняются краткими подсказками (как в TS `registry.ts`).
