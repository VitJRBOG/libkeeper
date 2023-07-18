const express = require('express')
const ejs = require('ejs')
const syncrequest = require('sync-request')
const crypto = require('crypto');

const APP_PORT = Number(process.env.APP_PORT)
const APP_HOST = String(process.env.APP_HOST)
const API_PORT = Number(process.env.API_PORT)
const API_HOST = String(process.env.API_HOST)

function LaunchServer() {
    osSignalsReception()

    let app = makeServerApp()
    app = makeHandlers(app)

    app.listen(APP_PORT, APP_HOST, function () {
        console.log(`server started at http://${APP_HOST}:${APP_PORT}`)
    })
}

function makeServerApp() {
    const app = express()

    app.set('views', './views')
    app.set('view engine', 'ejs')

    app.use(express.static('./public'))
    app.use(express.json())
    app.use(express.urlencoded({ extended: true }))

    return app
}

function makeHandlers(app) {
    var data = {
        'notes_list': null,
        'note_versions': null
    }

    app.get('/', function (req, res) {
        data['note_versions'] = null

        data['notes_list'] = fetchNotesList()

        res.render('main', data)
    })

    app.get('/note', function (req, res) {
        data['notes_list'] = fetchNotesList()
        data['note_versions'] = fetchNoteVersions(req.query.id)

        res.render('main', data)
    })

    app.post('/note', function (req, res) {
        let full_text = req.body.full_text
        let title = req.body.title
        let c_date = req.body.c_date
        let checksum = crypto.createHash('md5').update(full_text).digest('hex')

        let values = {
            full_text: full_text,
            title: title,
            c_date: c_date,
            checksum: checksum
        }

        createNewNote(values)

        res.redirect(`/`)
    })

    app.put('/note', function (req, res) {
        let note_id = req.query.id

        let full_text = req.body.full_text
        let title = req.body.title
        let c_date = req.body.c_date
        let checksum = crypto.createHash('md5').update(full_text).digest('hex')

        let values = {
            full_text: full_text,
            title: title,
            c_date: c_date,
            checksum: checksum,
            note_id: note_id
        }

        updateNote(values)

        res.redirect(303, `/note?id=${note_id}`)
    })

    app.delete('/note', function (req, res) {
        let note_id = req.query.id

        deleteNote(note_id)

        res.redirect('/')
    })

    return app
}

function fetchNotesList() {
    let u = `http://${API_HOST}:${API_PORT}/notes`

    var res = syncrequest('GET', u)

    let result = JSON.parse(res.getBody('utf-8'))

    let notes_list = null
    if (result.hasOwnProperty('response')) {
        notes_list = result['response']
    } else if (result.hasOwnProperty('error')) {
        console.log(`error: ${result['error']}`)
    }

    for (let i = 0; i < notes_list.length; i++) {
        notes_list[i]['c_date'] = _formatDate(notes_list[i]['c_date'])
    }

    return notes_list
}

function _formatDate(strUnixTS) {
    let unixTS = parseInt(strUnixTS)

    let date = new Date(unixTS * 1000)

    let full_date = ''

    if (date.getDate() < 10) {
        full_date += `0${date.getDate()}`
    } else {
        full_date += `${date.getDate()}`
    }

    if (date.getMonth() + 1 < 10) {
        full_date += `.0${date.getMonth() + 1}`
    } else {
        full_date += `.${date.getMonth() + 1}`
    }

    full_date += `.${date.getFullYear()}`

    let full_time = ''

    if (date.getHours() < 10) {
        full_time += `0${date.getHours()}`
    } else {
        full_time += `${date.getHours()}`
    }

    if (date.getMinutes() < 10) {
        full_time += `:0${date.getMinutes()}`
    } else {
        full_time += `:${date.getMinutes()}`
    }

    if (date.getSeconds() < 10) {
        full_time += `:0${date.getSeconds()}`
    } else {
        full_time += `:${date.getSeconds()}`
    }

    return `${full_date} ${full_time}`
}

function fetchNoteVersions(note_id) {
    let u = `http://${API_HOST}:${API_PORT}/note?` + new URLSearchParams({ note_id: note_id }).toString()

    var res = syncrequest('GET', u)

    let result = JSON.parse(res.getBody('utf-8'))

    let note_versions = null

    if (result.hasOwnProperty('response')) {
        note_versions = result['response']
    } else if (result.hasOwnProperty('error')) {
        console.log(`error: ${result['error']}`)
    }

    return note_versions
}

function createNewNote(values) {
    let u = `http://${API_HOST}:${API_PORT}/note`

    let str_params = ''

    let keys = Object.keys(values)

    for (let i = 0; i < keys.length; i++) {
        str_params += `${keys[i]}=${encodeURIComponent(values[keys[i]])}`
        if (i < keys.length - 1) {
            str_params += '&'
        }
    }

    var res = syncrequest('POST', u, {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: str_params
    })

    let result = {}
    let body = res.getBody('utf-8')
    if (body.length > 0) {
        result = JSON.parse(body)
    }

    if (result.hasOwnProperty('error')) {
        console.log(`error: ${result['error']}`)
    }
}

function updateNote(values) {
    let u = `http://${API_HOST}:${API_PORT}/note`

    let str_params = ''

    let keys = Object.keys(values)

    for (let i = 0; i < keys.length; i++) {
        str_params += `${keys[i]}=${encodeURIComponent(values[keys[i]])}`
        if (i < keys.length - 1) {
            str_params += '&'
        }
    }

    var res = syncrequest('PUT', u, {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: str_params
    })

    let result = {}
    let body = res.getBody('utf-8')
    if (body.length > 0) {
        result = JSON.parse(body)
    }

    if (result.hasOwnProperty('error')) {
        console.log(`error: ${result['error']}`)
    }
}

function deleteNote(note_id) {
    let u = `http://${API_HOST}:${API_PORT}/note?` + new URLSearchParams({ note_id: note_id }).toString()

    var res = syncrequest('DELETE', u)

    let result = {}
    let body = res.getBody('utf-8')
    if (body.length > 0) {
        result = JSON.parse(body)
    }

    if (result.hasOwnProperty('error')) {
        console.log(`error: ${result['error']}`)
    }
}

function osSignalsReception() {
    process.on('SIGINT' || 'SIGTERM', function () {
        console.log('program exited successfully')
        process.exit(0)
    })
}

LaunchServer()
