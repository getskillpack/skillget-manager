# Контракт клиента skillget-manager ↔ registry HTTP API

Канон спецификации сервера: [`docs/registry-api.md`](https://github.com/getskillpack/registry/blob/main/docs/registry-api.md) (раздел **Compiled core**). Этот документ фиксирует поведение библиотеки `github.com/getskillpack/skillget-manager` как **скомпилированного клиента**.

## Переменные окружения

| Переменная | Назначение |
|------------|------------|
| `SKILLGET_REGISTRY_URL` | Базовый URL API **без** завершающего `/`: `{origin}/api/v1`. Имеет приоритет. |
| `SKPKG_REGISTRY_URL` | Устаревший алиас для базы API. |
| `SKILLGET_REGISTRY_READ_TOKEN` | Bearer для **GET** к `/api/v1/*` и для **GET** по `archive_url`, если оператор реестра включил `REGISTRY_READ_TOKEN`. |
| `SKILLGET_REGISTRY_TOKEN` | Bearer для **записи** (`POST /skills`, `DELETE …/versions/…`). |
| `SKILLGET_TOKEN` | Устаревший алиас для токена записи. |

Если задан только `SKILLGET_REGISTRY_READ_TOKEN`, анонимные запросы к API не выполняются: чтение идёт с этим Bearer. Если `SKILLGET_REGISTRY_READ_TOKEN` пуст, для чтения используется тот же токен, что и для записи (`SKILLGET_REGISTRY_TOKEN` / `SKILLGET_TOKEN`), если он задан — это упрощает сценарий «один секрет на приватный реестр».

Публичный реестр по умолчанию: `https://registry.skpkg.org/api/v1` (без обязательных токенов для чтения).

## HTTP: маршруты и коды

Пути ниже относительны к базе из `SKILLGET_REGISTRY_URL`.

| Действие в библиотеке | Метод и путь | Успех | Ошибки (ожидания клиента) |
|------------------------|--------------|-------|---------------------------|
| `SearchSkills` | `GET /skills?…` | `200` + JSON | `401` при read-token на сервере |
| `ResolveInstallTarget` (без pin) | `GET /skills/{name}` | `200` + JSON detail | `404`, `401` |
| `ResolveInstallTarget` (все случаи) | `GET /skills/{name}/versions/{version}` | `200` + JSON version | `404`, **`410`** после yank, `401` |
| `PublishSkill` | `POST /skills` (multipart `manifest` + `archive`) | **`201`** пустое тело | `400`, `401`, `409`, `503` без write token |
| (сервер) yank | `DELETE /skills/{name}/versions/{version}` | **`204`** | `401`, `404` |

Тела ошибок reference-сервера — короткий `text/plain` (не JSON). Сообщения об ошибках в Go дополняются подсказками (`registry_errors.go`).

## JSON: формы ответов

- **Список** `GET /skills` — объект с полями `data[]`, `meta` (`total`, `limit`, `offset`). Элементы `data` содержат как минимум `name`, `latest_version`, `created_at`, `description`, `author` (как в спецификации).
- **Деталь скилла** `GET /skills/{name}` — поля `name`, `description`, `author`, `created_at`, **`versions`** — объект, ключи = semver-строки версий, значения: `manifest`, `checksum` (`sha256:` + 64 hex), `archive_url`, `published_at`, `yanked`.
- **Версия** `GET /skills/{name}/versions/{version}` — `name`, `version`, `manifest`, `archive_url`, `checksum`.

Библиотека **не** использует несуществующее в каноне поле `latest_version` на ответе `GET /skills/{name}`: для установки без pin выбирается **наибольшая по semver** неснятая (`yanked: false`) версия из `versions`, в духе логики `latestNonYanked` в reference store.

## Публикация и архив

- `PublishSkill`: `multipart/form-data`, поля `manifest` (JSON-строка) и `archive` (файл `.tar.gz`), заголовок `Authorization: Bearer <write token>`.
- `DownloadSkillArchive`: GET по абсолютному `archive_url` из ответа версии; при непустом `RegistryReadBearer()` добавляется тот же `Authorization`, что и к API (для приватных `/downloads/*`).

## Проверка целостности

Если в метаданных версии присутствует `checksum` в формате `sha256:<64 hex>`, после скачивания выполняется сверка SHA-256 с содержимым файла.

## Соответствие ENGINEERING_REQUIREMENTS_SKPKG.md

Файл `plans/ENGINEERING_REQUIREMENTS_SKPKG.md` в onboarding/CTO workspace не входит в этот репозиторий. Стек менеджера — **Go**, контракт с реестром — **HTTP + JSON** по канону выше; расхождения с инженерным чеклистом фиксируются в рабочих тикетах Paperclip (см. раздел в `registry-api.md`).
