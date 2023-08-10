# API

### `[GET] /icons`

Returns a list of all the icons.

**Required parameters:**

No parameters required.

**Returns values:**

-   `id` - type: integer. Icon ID.
-   `type` - type: string. Icon assignment.
-   `path` - type: string. Path to the icon file.

**Server response examples:**

-   `STATUS 200 OK`

```
{
    response: [
        {
            "id": 1,
            "type": "category",
            "path": "uncategorised.png"
        },
        {
            "id": 2,
            "type": "category",
            "path": "trashed.png"
        },
        {
            "id": 3,
            "type": "category",
            "path": "noicon.png"
        }
    ]
}
```

### `[GET] /categories`

Returns a list of all the categories.

**Required parameters:**

No parameters required.

**Returns values:**

-   `id` - type: integer. Category ID.
-   `name` - type: string. Name of the category.
-   `immutable` - type: integer. Flag for the system categories.
-   `icon_id` - type: integer. Icon ID.

**Server response examples:**

-   `STATUS 200 OK`

```
{
    response: [
        {
            "id": 1,
            "name": "Uncategorised",
            "immutable": 1,
            "icon_id": 1
        },
        {
            "id": 2,
            "name": "Trashed",
            "immutable": 1,
            "icon_id": 2
        },
        {
            "id": 3,
            "name": "Some category",
            "immutable": 0,
            "icon_id": 3
        }
    ]
}
```

### `[POST] /category`

Creates a new category.

**Required parameters:**

-   `name` - type: string. Name of the category. **Required**.
-   `icon_id` - type: integer. Icon ID. **Required**.

**Server response examples:**

-   `STATUS 200 OK`

_Nothing._

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```

### `[PUT] /category`

Updates an existing category.

**Required parameters:**

-   `category_id` - type: integer. Category ID. **Required**.
-   `name` - type: string. Name of the category. **Required**.
-   `icon_id` - type: integer. Icon ID. **Required**.

**Server response examples:**

-   `STATUS 200 OK`

_Nothing._

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```

### `[DElETE] /category`

Deletes an existing category.

**Required parameters:**

-   `category_id` - type: integer. Category ID. **Required**.

**Server response examples:**

-   `STATUS 204 No Content`

_Nothing._

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```

### `[GET] /notes`

Returns a list of all the notes.

**Required parameters:**

No parameters required.

**Returns values:**

-   `id` - type: integer. Note ID.
-   `title` - type: string. Title of the note.
-   `c_date` - type: string. Note creation date in unix timestamp format.
-   `categories` - type: list of strings. A list of the category names of this note.

**Server response examples:**

-   `STATUS 200 OK`

```
{
    response: [
        {
            "id": 1,
            "title": "Title of the first note",
            "c_date": "1681241952",
            "categories": [
                "Uncategorised"
            ]
        },
        {
            "id": 2,
            "title": "Title of the second note",
            "c_date": "1682927889",
            "categories": [
                "Some category",
                "Another category"
            ]
        },
        {
            "id": 3,
            "title": "Title of the third note",
            "c_date": "1686574910",
            "categories": [
                "Trashed"
            ]
        }
    ]
}
```

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```

### `[GET] /note`

Returns a list of all versions of the note with specified `note_id`.

**Parameters required:**

-   `note_id` - type: integer. **Required**.

**Returns values:**

-   `id` - type: integer. Version ID.
-   `full_text` - type: string. Full text of the note.
-   `c_date` - type: string. Version creation date in unix timestamp format.
-   `checksum` - type: string. The check sum of the full text of the note.
-   `note_id` - type: integer. Note ID.

**Server response examples:**

-   `STATUS 200 OK`

```
{
    response: [
        {
            "id": 3,
            "full_text": "Title of the first note\nThird full text of the first note.",
            "c_date", "1681268383",
            "checksum": "7c776685ce6b0594738e570896d18ce1",
            "note_id": 1
        },
        {
            "id": 2,
            "full_text": "Title of the first note\nSecond full text of the first note.",
            "c_date", "1681251672",
            "checksum": "06763c7582351198187af3a2ae1901a2",
            "note_id": 1
        },
        {
            "id": 1,
            "full_text": "Title of the first note\nFirst full text of the first note.",
            "c_date", "1681241952",
            "checksum": "b92675dd01365fce9d4a724d5e1614e5",
            "note_id": 1
        }
    ]
}
```

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```

### `[POST] /note`

Creates a new note.

**Parameters required:**

-   `c_date` - type: string. Date in the format `yyyy-mm-dd hh:mm:ss -0000`. **Required**.
-   `checksum` - type: string. **Required**.
-   `title` - type: string.
-   `full_text` - type: string.
-   `categories` - type: string. Stringified list of the category names separated by comma. **Required**

**Server response examples:**

-   `STATUS 200 OK`

_Nothing._

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```

### `[PUT] /note`

Creates a new version of the existing note and changes relations to the categories.

**Parameters required (for full note update):**

-   `note_id` - type: integer. **Required**
-   `c_date` - type: string. Date in the format `yyyy-mm-dd hh:mm:ss -0000`. **Required**.
-   `checksum` - type: string. **Required**.
-   `title` - type: string.
-   `full_text` - type: string.
-   `categories` - type: string. Stringified list of the category names separated by comma.

**Parameters required (for categories update):**

-   `note_id` - type: integer. **Required**
-   `categories` - type: string. Stringified list of the category names separated by comma. **Required**

**Server response examples:**

-   `STATUS 200 OK`

_Nothing._

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```

### `[DELETE] /note`

Deletes an existing note and all associated versions.

**Required parameters:**

-   `note_id` - type: integer. Note ID. **Required**.

**Server response examples:**

-   `STATUS 204 No Content`

_Nothing._

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```
