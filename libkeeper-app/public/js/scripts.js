function highlintCategory() {
    let queryString = window.location.search
    let urlParams = new URLSearchParams(queryString)
    if (urlParams.has('category_id')) {
        let categoryID = urlParams.get('category_id')
        let tag = document.getElementById(`category-${categoryID}`)
        tag.classList.remove('category-button-regular')
        tag.classList.add('category-button-highlinted')
    } else {
        let tag = document.getElementById('category-all')
        tag.classList.remove('category-button-regular')
        tag.classList.add('category-button-highlinted')
    }
}

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

function hideErrorBlock(div_id) {
    let tag = document.getElementById(div_id)

    if (tag.style['visibility'] == 'visible') {
        tag.style['visibility'] = 'hidden'
    }
}

function toggleCategoryCreationPromptDisplay() {
    let tag = document.getElementById('category-creation-prompt')
    let tagIconsList = document.getElementById('category-icons-list-prompt')

    if (tag.style['visibility'] == 'hidden') {
        tag.style['visibility'] = 'visible'
    } else if (tag.style['visibility'] == 'visible') {
        tag.style['visibility'] = 'hidden'
    }

    if (tagIconsList.style['visibility'] == 'visible') {
        tagIconsList.style['visibility'] = 'hidden'
    }
}

function toggleCaregoryIconsListPromptDisplay() {
    let tag = document.getElementById('category-icons-list-prompt')

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

function filterByCategory(category_id) {
    if (category_id === '-1') {
        window.location.replace('/')
    } else {
        let newLocation = `/?category_id=${category_id}`
        window.location.replace(newLocation)
    }
}

function createCategory() {
    let name = document.getElementById('category-name-textfield').value
    let icon_id = ''

    let radioBtns = document.getElementsByName('category-icon-buttons')

    for (let i = 0; i < radioBtns.length; i++) {
        if (radioBtns[i].checked) {
            icon_id = radioBtns[i].value
        }
    }

    fetch('/category', {
        method: 'POST',
        body: JSON.stringify({
            name: name,
            icon_id: icon_id,
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
    let categories_list = document.getElementById('note-categories-list').getElementsByTagName('ul')
    console.log(categories_list)
    if (categories_list.length == 0) {
        return []
    }

    let checkboxes = categories_list[0].getElementsByTagName('li')

    let categories = []
    for (let i = 0; i < checkboxes.length; i++) {
        console.log(checkboxes[i])
        if (checkboxes[i].getElementsByTagName('input')[0].checked) {
            categories.push(checkboxes[i].getElementsByTagName('label')[0].textContent)
        }
    }

    return categories
}

function deleteNote(categories) {
    let queryString = window.location.search
    let urlParams = new URLSearchParams(queryString)
    if (urlParams.has('id')) {
        let note_id = urlParams.get('id')
        let queryURL = `/note?id=${note_id}`
        if (categories == 'Trashed') {
            fetch(queryURL, {
                method: 'delete'
            }).then(response => {
                if (response.redirected) {
                    window.location.href = response.url
                }
            })
        } else {
            fetch(queryURL, {
                method: 'put',
                body: JSON.stringify({
                    categories: 'Trashed',
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
    }
}