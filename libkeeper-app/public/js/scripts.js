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

function openNewCanvas() {
    window.location.replace('/')
}

function createNote(full_text) {
    let title = _composeTitle(full_text)

    let c_date = _getDate()

    fetch("/note", {
        method: 'post',
        body: JSON.stringify({
            full_text: full_text,
            title: title,
            c_date: c_date
        }),
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
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
    let tzOffset = now.getTimezoneOffset()
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