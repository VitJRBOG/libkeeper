function highlintNote() {
    let queryString = window.location.search
    let urlParams = new URLSearchParams(queryString)
    if (urlParams.has('id')) {
        let elementID = urlParams.get('id')
        let tag = document.getElementById(`note_${elementID}`)
        tag.classList.remove('note-announce-regular')
        tag.classList.add('note-announce-highlinted')
    }
}

function highlintVersion() {
    let queryString = window.location.search
    let urlParams = new URLSearchParams(queryString)
    if (urlParams.has('version_id')) {
        let elementID = urlParams.get('version_id')
        let tag = document.getElementById(`version_${elementID}`)
        tag.classList.remove('version-announce-regular')
        tag.classList.add('version-announce-highlinted')
    }
}

function toggleCategoryCreationPromptDisplay() {
    let tag = document.getElementById('category-creation-prompt')

    if (tag.style['visibility'] == 'hidden') {
        tag.style['visibility'] = 'visible'
    } else if (tag.style['visibility'] == 'visible') {
        tag.style['visibility'] = 'hidden'
    }
}

function toggleVersionsListDisplay() {
    let tag = document.getElementById('versions-list')

    if (tag.style['visibility'] == 'hidden') {
        tag.style['visibility'] = 'visible'
    } else if (tag.style['visibility'] == 'visible') {
        tag.style['visibility'] = 'hidden'
    }
}

function toggleCategoriesListDisplay() {
    let tag = document.getElementById('note-categories-list')

    if (tag.style['visibility'] == 'hidden') {
        tag.style['visibility'] = 'visible'
    } else if (tag.style['visibility'] == 'visible') {
        tag.style['visibility'] = 'hidden'
    }
}

function filterByCategory(category) {
    if (category === 'All') {
        window.location.replace('/')
    } else {
        let newLocation = `/?category=${encodeURIComponent(category)}`
        window.location.replace(newLocation)
    }
}

function deleteCategory(category_id, category_is_immutable) {
    if (category_is_immutable == 1) {
        console.error('unable to delete the immutable category')
    } else {
        fetch(`/category?id=${category_id}`, {
            method: 'delete'
        }).then(response => {
            if (response.redirected) {
                window.location.href = response.url
            }
        })
    }
}

function openNewCanvas() {
    window.location.replace('/')
}

function handleTheNote(full_text) {
    let title = _composeTitle(full_text)
    let c_date = _getDate()
    let categories = _getCategories()

    let queryURL = ''
    let queryMethod = ''

    let queryString = window.location.search
    let urlParams = new URLSearchParams(queryString)
    if (urlParams.has('id')) {
        let note_id = urlParams.get('id')

        queryURL = `/note?id=${note_id}`
        queryMethod = 'put'
    } else {
        queryURL = '/note'
        queryMethod = 'post'
    }

    fetch(queryURL, {
        method: queryMethod,
        body: JSON.stringify({
            full_text: full_text,
            title: title,
            c_date: c_date,
            categories: categories,
        }),
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
    }).then(response => {
        if (response.redirected) {
            window.location.href = response.url
        }
    })
}

function _composeTitle(full_text) {
    let n = full_text.indexOf('\n')
    if (n > -1 && n <= 45) {
        return full_text.substring(0, n)
    } else if (full_text.length > 45) {
        return `${full_text.substring(0, 42)}...`
    } else {
        return full_text
    }
}

function _getDate() {
    let now = new Date()

    let tz = _getTimezone(now)
    let c_date = _formatDate(now, tz)

    return c_date
}

function _getTimezone(now) {
    let tzOffset = now.getTimezoneOffset() * -1
    let tz = (tzOffset / 60)

    if (tz > 0) {
        let filler = ''

        if (tz < 10) {
            filler = '0'
        }

        tz = `+${filler}${tz}`
    } else if (tz < 0) {
        let filler = ''

        if (tz > -10) {
            filler = '0'
        }

        tz = `-${filler}${tz * -1}`
    }

    tz = tz.replace('.5', '3')

    let filler = '0'.repeat(5 - tz.length)

    tz = `${tz}${filler}`

    return tz
}

function _formatDate(now, tz) {
    let year = now.getFullYear()
    let month = now.getMonth() + 1
    let day = now.getDate()

    if (day < 10) {
        day = `0${day}`
    }

    if (month < 10) {
        month = `0${month}`
    }

    let sec = now.getSeconds()

    if (sec < 10) {
        sec = `0${sec}`
    }

    let min = now.getMinutes()

    if (min < 10) {
        min = `0${min}`
    }

    let hours = now.getHours()

    if (hours < 10) {
        hours = `0${hours}`
    }

    return `${year}-${month}-${day} ${hours}:${min}:${sec} ${tz}`
}

function _getCategories() {
    let checkboxes = document.getElementById('note-categories-list').getElementsByTagName('ul')[0].getElementsByTagName('li')

    let categories = []

    for (let i = 0; i < checkboxes.length; i++) {
        console.log(checkboxes[i])
        if (checkboxes[i].getElementsByTagName('input')[0].checked) {
            categories.push(checkboxes[i].getElementsByTagName('label')[0].textContent)
        }
    }

    return categories
}

function deleteNote() {
    let queryString = window.location.search
    let urlParams = new URLSearchParams(queryString)
    if (urlParams.has('id')) {
        let note_id = urlParams.get('id')
        fetch(`/note?id=${note_id}`, {
            method: 'delete'
        }).then(response => {
            if (response.redirected) {
                window.location.href = response.url
            }
        })
    }
}