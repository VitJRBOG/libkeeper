const express = require('express')
const ejs = require('ejs')
const syncrequest = require('sync-request')
const crypto = require('crypto')
const fs = require('fs')

const PUBLIC_DIR = String(process.env.PUBLIC_DIR)
const VIEWS_DIR = String(process.env.VIEWS_DIR)
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

    app.set('views', `${VIEWS_DIR}`)
    app.set('view engine', 'ejs')

    app.use(express.static(`${PUBLIC_DIR}`))
    app.use(express.json())
    app.use(express.urlencoded({ extended: true }))

    return app
}

function makeHandlers(app) {
    var data = {
        'categories_list': null,
        'notes_list': null,
        'note_versions': null,
        'current_category': null,
        'current_note': null,
        'current_version': null
    }

    app.get('/', function (req, res) {
        data['categories_list'] = null
        data['current_category'] = null
        data['current_note'] = null
        data['current_version'] = null
        data['note_versions'] = null

        data['categories_list'] = fetchCategoriesList()
        data['categories_list'] = _categoriesIconFinding(data['categories_list'])

        data['notes_list'] = fetchNotesList()

        if (typeof req.query.category !== 'undefined') {
            let selectedCategory = decodeURIComponent(req.query.category)
            for (let i = 0; i < data['categories_list'].length; i++) {
                if (data['categories_list'][i].name === selectedCategory) {
                    data['current_category'] = data['categories_list'][i]
                }
            }
            for (let i = 0; i < data['notes_list'].length; i++) {
                if (!data['notes_list'][i].categories.includes(data['current_category'].name)) {
                    data['notes_list'].splice(i, 1)
                    i--
                }
            }
        }

        res.render('main', data)
    })

    app.post('/category', function (req, res) {
        createCategory(req.body.category_name)

        res.redirect(`/`)
    })

    app.delete('/category', function (req, res) {
        let category_id = req.query.id

        deleteCategory(category_id)

        res.redirect(`/`)
    })

    app.get('/note', function (req, res) {
        data['categories_list'] = fetchCategoriesList()
        data['categories_list'] = _categoriesIconFinding(data['categories_list'])

        data['notes_list'] = fetchNotesList()

        if (typeof req.query.category !== 'undefined') {
            let selectedCategory = decodeURIComponent(req.query.category)
            for (let i = 0; i < data['categories_list'].length; i++) {
                if (data['categories_list'][i].name === selectedCategory) {
                    data['current_category'] = data['categories_list'][i]
                }
            }
            for (let i = 0; i < data['notes_list'].length; i++) {
                if (!data['notes_list'][i].categories.includes(data['current_category'].name)) {
                    data['notes_list'].splice(i, 1)
                    i--
                }
            }
        }

        for (let i = 0; i < data['notes_list'].length; i++) {
            if (data['notes_list'][i].id == req.query.id) {
                data['current_note'] = data['notes_list'][i]
            }
        }

        data['note_versions'] = fetchNoteVersions(req.query.id)

        if (req.query.version_id) {
            for (let i = 0; i < data['note_versions'].length; i++) {
                if (data['note_versions'][i].id == req.query.version_id) {
                    data['current_version'] = data['note_versions'][i]
                }
            }
        } else {
            data['current_version'] = data['note_versions'][0]
        }

        res.render('main', data)
    })

    app.post('/note', function (req, res) {
        let full_text = req.body.full_text
        let title = req.body.title
        let c_date = req.body.c_date
        let checksum = crypto.createHash('md5').update(full_text).digest('hex')
        let categories = ''
        if (req.body.categories.length > 0 && typeof req.body.categories !== 'undefined') {
            categories = req.body.categories
        } else {
            categories = 'Uncategorised'
        }

        let values = {
            full_text: full_text,
            title: title,
            c_date: c_date,
            checksum: checksum,
            categories: categories
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
        let categories = ''
        if (req.body.categories.length > 0 && typeof req.body.categories !== 'undefined') {
            categories = req.body.categories
        } else {
            categories = 'Uncategorised'
        }

        let values = {
            full_text: full_text,
            title: title,
            c_date: c_date,
            checksum: checksum,
            note_id: note_id,
            categories: categories,
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

function fetchCategoriesList() {
    let u = `http://${API_HOST}:${API_PORT}/categories`

    var res = syncrequest('GET', u)

    let result = JSON.parse(res.getBody('utf-8'))

    let categories_list = null
    if (result.hasOwnProperty('response')) {
        categories_list = result['response']
    } else if (result.hasOwnProperty('error')) {
        console.log(`error: ${result['error']}`)
    }

    return categories_list
}

function _categoriesIconFinding(categories_list) {
    for (let i = 0; i < categories_list.length; i++) {
        if (fs.existsSync(`${PUBLIC_DIR}/img/icons/category-${categories_list[i]['name']}.png`)) {
            categories_list[i]['icon'] = `/img/icons/category-${categories_list[i]['name']}.png`
        } else {
            categories_list[i]['icon'] = `/img/icons/category-uncategorised.png`
        }
    }

    return categories_list
}

function createCategory(category_name) {
    let u = `http://${API_HOST}:${API_PORT}/category`

    let str_params = `name=${encodeURIComponent(category_name)}`

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

function deleteCategory(category_id) {
    let u = `http://${API_HOST}:${API_PORT}/category?` + new URLSearchParams({ category_id: category_id }).toString()

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

    for (let i = 0; i < note_versions.length; i++) {
        note_versions[i]['c_date'] = _formatDate(note_versions[i]['c_date'])
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
