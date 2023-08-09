const express = require('express')
const ejs = require('ejs')
const syncrequest = require('sync-request')
const crypto = require('crypto')
const fs = require('fs')
const { error } = require('console')

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
        'error': null,
        'categories_list': null,
        'notes_list': null,
        'note_versions': null,
        'current_category': null,
        'current_note': null,
        'current_version': null
    }

    app.get('/', function (req, res) {
        data['error'] = []
        data['categories_list'] = null
        data['icons_list'] = null
        data['current_category'] = null
        data['current_note'] = null
        data['current_version'] = null
        data['note_versions'] = null

        let result = fetchCategoriesList()
        if (typeof result == 'string') {
            data['error'].push(result)
        } else {
            data['categories_list'] = result
        }

        result = fetchIconsList()
        if (typeof result == 'string') {
            data['error'].push(result)
        } else {
            data['icons_list'] = result
        }

        if (data['categories_list'] !== null && data['icons_list'] !== null) {
            data['categories_list'] = _categoriesIconFinding(data['categories_list'], data['icons_list'])
        }

        result = fetchNotesList()
        if (typeof result == 'string') {
            data['error'].push(result)
        } else {
            data['notes_list'] = result
        }

        if (typeof req.query.category_id !== 'undefined') {
            let selectedCategoryID = decodeURIComponent(req.query.category_id)
            if (data['categories_list'] !== null) {
                for (let i = 0; i < data['categories_list'].length; i++) {
                    if (data['categories_list'][i].id == selectedCategoryID) {
                        data['current_category'] = data['categories_list'][i]
                    }
                }
            }
            if (data['notes_list'] !== null) {
                for (let i = 0; i < data['notes_list'].length; i++) {
                    if (!data['notes_list'][i].categories.includes(data['current_category'].name)) {
                        data['notes_list'].splice(i, 1)
                        i--
                    }
                }
            }
        }

        res.render('main', data)
    })

    app.post('/category', function (req, res) {
        let values = {
            name: req.body.name,
            icon_id: req.body.icon_id
        }

        createCategory(values)

        res.redirect(`/`)
    })

    app.delete('/category', function (req, res) {
        let category_id = req.query.id

        deleteCategory(category_id)

        res.redirect(`/`)
    })

    app.get('/note', function (req, res) {
        data['error'] = null

        let result = fetchCategoriesList()
        if (typeof result == 'string') {
            data['error'].push(result)
        } else {
            data['categories_list'] = result
        }

        result = fetchIconsList()
        if (typeof result == 'string') {
            data['error'].push(result)
        } else {
            data['icons_list'] = result
        }

        if (data['categories_list'] !== null && data['icons_list'] !== null) {
            data['categories_list'] = _categoriesIconFinding(data['categories_list'], data['icons_list'])
        }

        result = fetchNotesList()
        if (typeof result == 'string') {
            data['error'].push(result)
        } else {
            data['notes_list'] = result
        }

        if (typeof req.query.category_id !== 'undefined') {
            let selectedCategoryID = decodeURIComponent(req.query.category_id)
            if (data['categories_list'] !== null) {
                for (let i = 0; i < data['categories_list'].length; i++) {
                    if (data['categories_list'][i].id == selectedCategoryID) {
                        data['current_category'] = data['categories_list'][i]
                    }
                }
            }
            if (data['notes_list'] !== null) {
                for (let i = 0; i < data['notes_list'].length; i++) {
                    if (!data['notes_list'][i].categories.includes(data['current_category'].name)) {
                        data['notes_list'].splice(i, 1)
                        i--
                    }
                }
            }
        }

        if (data['notes_list'] !== null) {
            for (let i = 0; i < data['notes_list'].length; i++) {
                if (data['notes_list'][i].id == req.query.id) {
                    data['current_note'] = data['notes_list'][i]
                }
            }
        }

        result = fetchNoteVersions(req.query.id)
        if (typeof result == 'string') {
            data['error'].push(result)
        } else {
            data['note_versions'] = result
        }

        if (req.query.version_id) {
            if (data['note_versions'] !== null) {
                for (let i = 0; i < data['note_versions'].length; i++) {
                    if (data['note_versions'][i].id == req.query.version_id) {
                        data['current_version'] = data['note_versions'][i]
                    }
                }
            }
        } else {
            if (data['note_versions'] !== null) {
                data['current_version'] = data['note_versions'][0]
            }
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
    try {
        let u = `http://${API_HOST}:${API_PORT}/categories`

        var res = syncrequest('GET', u)

        let values = JSON.parse(res.getBody('utf-8'))

        let result = null
        if (values.hasOwnProperty('response')) {
            result = values['response']
        } else if (values.hasOwnProperty('error')) {
            result = values['error']
        }

        return result

    } catch (err) {
        error(err)
        return 'Unable to fetch categories list from server.'
    }
}

function fetchIconsList() {
    try {
        let u = `http://${API_HOST}:${API_PORT}/icons`

        var res = syncrequest('GET', u)

        let values = JSON.parse(res.getBody('utf-8'))

        let result = null
        if (values.hasOwnProperty('response')) {
            result = values['response']
        } else if (values.hasOwnProperty('error')) {
            result = values['error']
        }

        return result

    } catch (err) {
        error(err)
        return 'Unable to fetch icons list from server.'
    }
}

function _categoriesIconFinding(categories_list, icons_list) {
    for (let i = 0; i < categories_list.length; i++) {
        for (let j = 0; j < icons_list.length; j++) {
            if (categories_list[i].icon_id == icons_list[j].id) {
                categories_list[i]['icon'] = icons_list[j].path
            }
        }
    }

    return categories_list
}

function createCategory(values) {
    let u = `http://${API_HOST}:${API_PORT}/category`

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
    try {
        let u = `http://${API_HOST}:${API_PORT}/notes`

        var res = syncrequest('GET', u)

        let values = JSON.parse(res.getBody('utf-8'))

        let result = null
        if (values.hasOwnProperty('response')) {
            let notes_list = values['response']
            for (let i = 0; i < notes_list.length; i++) {
                notes_list[i]['c_date'] = _formatDate(notes_list[i]['c_date'])
            }
            result = notes_list
        } else if (values.hasOwnProperty('error')) {
            result = values['error']
        }

        return result

    } catch (err) {
        error(err)
        return 'Unable to fetch notes list from server.'
    }
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
    try {
        let u = `http://${API_HOST}:${API_PORT}/note?` + new URLSearchParams({ note_id: note_id }).toString()

        var res = syncrequest('GET', u)

        let values = JSON.parse(res.getBody('utf-8'))

        let result = null
        if (values.hasOwnProperty('response')) {
            note_versions = values['response']
            for (let i = 0; i < note_versions.length; i++) {
                note_versions[i]['c_date'] = _formatDate(note_versions[i]['c_date'])
            }
            result = note_versions
        } else if (result.hasOwnProperty('error')) {
            result = values['error']
        }

        return result
    }
    catch (err) {
        error(err)
        return 'Unable to fetch versions list from server.'
    }
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
