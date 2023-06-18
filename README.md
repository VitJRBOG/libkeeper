# Installation

1. `git clone https://github.com/VitJRBOG/libkeeper-api`
2. `cd libkeeper-api`
3. `docker-compose up`

# API

### `[GET] /notes`

Returns a list of all the notes.

**Required parameters:**

No parameters required.

**Returns values:**

-   `id` - type: integer. Note ID.
-   `title` - type: string. Title of the note.
-   `c_date` - type: string. Note creation date in unix timestamp format.

**Server response examples:**

-   `STATUS 200 OK`

```
{
    response: [
        {
            "id": 1,
            "title": "Title of the first note",
            "c_date": "1681241952"
        },
        {
            "id": 2,
            "title": "Title of the second note",
            "c_date": "1682927889"
        },
        {
            "id": 3,
            "title": "Title of the third note",
            "c_date": "1686574910"
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

-   `note_id` - type: integer.

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

Creates a new version of the existing note.

**Parameters required:**

-   `note_id` - type: integer. **Required**
-   `c_date` - type: string. Date in the format `yyyy-mm-dd hh:mm:ss -0000`. **Required**.
-   `checksum` - type: string. **Required**.
-   `title` - type: string.
-   `full_text` - type: string.

**Server response examples:**

-   `STATUS 200 OK`

_Nothing._

-   `STATUS [http_code]`

```
{
    "error": "text of the error"
}
```
