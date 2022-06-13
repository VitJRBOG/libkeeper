import datetime
import time
import hashlib

from tools import logging
from django.http import HttpRequest, JsonResponse
from . import models


def add_note(request: HttpRequest) -> JsonResponse:
    if request.method != 'POST':
        return JsonResponse({
            'status_code': 400,
            'error': f'{request.method} is not allowed.',
        })

    if request.POST.get('text') is None or \
       request.POST.get('text') == '' or \
       request.POST.get('text') == ' ':
        return JsonResponse({
            'status_code': 400,
            'error': 'attribute "text" is required',
        })
    text = str(request.POST.get('text'))
    err, title = __compose_title(text)

    if err != '':
        return JsonResponse({
            'status_code': 500,
            'error': err,
        })

    err, date = __get_date()

    if err != '':
        return JsonResponse({
            'status_code': 500,
            'error': err,
        })

    try:
        note = models.Note()
        note.create(title, date)
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        return JsonResponse({
            'status_code': 500,
            'error': 'note creation error',
        })

    last_note_object = models.Note.objects.last()
    note_id: int
    if last_note_object is not None:
        note_id = last_note_object.id
    else:
        e = 'last_note_object is None'
        logging.Logger('critical').critical(e, exc_info=True)
        return JsonResponse({
            'status_code': 500,
            'error': 'note creation error',
        })

    err, checksum = __gen_checksum(text)

    if err != '':
        return JsonResponse({
            'status_code': 500,
            'error': err,
        })

    try:
        version = models.Version()
        version.create(text, date, checksum, note_id)
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        note.delete(note_id)
        return JsonResponse({
            'status_code': 500,
            'error': 'note creation error',
        })

    last_version_object = models.Version.objects.last()
    version_id: int
    if last_version_object is not None:
        version_id = last_version_object.id
    else:
        e = 'last_version_object is None'
        logging.Logger('critical').critical(e, exc_info=True)
        return JsonResponse({
            'status_code': 500,
            'error': 'note creation error',
        })

    return JsonResponse({
        'status_code': 200,
        'count': 1,  # because item always single
        'items': [{
            'note_id': note_id,
            'title': title,
            'date_creation': date,
            'version_id': version_id,
            'text': text,
            'date_modified': date,
            'checksum': checksum,
        }]
    })


def update_note(request: HttpRequest) -> JsonResponse:
    if request.method != 'POST':
        return JsonResponse({
            'status_code': 400,
            'error': f'{request.method} is not allowed.',
        })

    id_ = request.POST.get('id')
    if id_ is None:
        return JsonResponse({
            'status_code': 400,
            'error': '"id" attribute is required',
        })

    note_id: int
    try:
        note_id = int(id_)
    except ValueError:
        return JsonResponse({
            'status_code': 400,
            'error': '"id" attribute must be integer',
        })

    try:
        note = models.Note.objects.get(id=note_id)
    except models.Note.DoesNotExist:
        return JsonResponse({
            'status_code': 404,
            'error': 'no objects with the specified ID were found',
        })

    text = request.POST.get('text')
    if text is None or text == '' or text == ' ':
        return JsonResponse({
            'status_code': 400,
            'error': 'attribute "text" is required',
        })

    err, title = __compose_title(text)
    if err != '':
        return JsonResponse({
            'status_code': 500,
            'error': err,
        })

    err, date = __get_date()
    if err != '':
        return JsonResponse({
            'status_code': 500,
            'error': err,
        })

    err, checksum = __gen_checksum(text)
    if err != '':
        return JsonResponse({
            'status_code': 500,
            'error': err,
        })

    try:
        version = models.Version()
        version.create(text, date, checksum, note_id)
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        return JsonResponse({
            'status_code': 500,
            'error': 'note updation error',
        })

    last_version_object = models.Version.objects.last()
    version_id: int
    if last_version_object is not None:
        version_id = last_version_object.id
    else:
        e = 'last_version_object is None'
        logging.Logger('critical').critical(e, exc_info=True)
        return JsonResponse({
            'status_code': 500,
            'error': 'note updation error',
        })

    try:
        note.title = title
        note.save()
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        version.delete(note_id)
        return JsonResponse({
            'status_code': 500,
            'error': 'note updation error',
        })

    return JsonResponse({
        'status_code': 200,
        'count': 1,  # because item always single
        'items': [{
            'note_id': note_id,
            'title': title,
            'date_creation': date,
            'version_id': version_id,
            'text': text,
            'date_modified': date,
            'checksum': checksum,
        }]
    })


def get_notes(request: HttpRequest) -> JsonResponse:
    if request.method != 'GET':
        return JsonResponse({
            'status_code': 400,
            'error': f'{request.method} is not allowed.',
        })

    try:
        queryset = models.Note.objects.values()

        return JsonResponse({
            'status_code': 200,
            'count': queryset.count(),
            'items': list(queryset),
        })
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        return JsonResponse({
            'status_code': 500,
            'error': 'note selection error',
        })


def get_versions(request: HttpRequest) -> JsonResponse:
    if request.method != 'GET':
        return JsonResponse({
            'status_code': 400,
            'error': f'{request.method} is not allowed.',
        })

    id_ = request.GET.get('id')

    if id_ is None:
        return JsonResponse({
            'status_code': 400,
            'error': '"id" attribute is required',
        })

    note_id: int
    try:
        note_id = int(id_)
    except ValueError:
        return JsonResponse({
            'status_code': 400,
            'error': '"id" attribute must be integer',
        })

    try:
        queryset = models.Version.objects.filter(
            note_id=note_id).order_by('date').reverse().values()

        return JsonResponse({
            'status_code': 200,
            'count': queryset.count(),
            'items': list(queryset),
        })
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        return JsonResponse({
            'status_code': 500,
            'error': 'note selection error',
        })


def __compose_title(text: str) -> tuple[str, str]:
    try:
        if len(text) > 47:
            if '\n' in text[:50]:
                return ('', text[:text.find('\n')])

            return ('', f'{text[:46]}...')
        return ('', text)
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        return ('error of title composing', '')


def __get_date() -> tuple[str, int]:
    try:
        today = datetime.datetime.now()

        return ('', int(time.mktime(today.timetuple())))
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        return ('error of getting date', 0)


def __gen_checksum(text: str) -> tuple[str, str]:
    try:
        hash = hashlib.md5(text.encode('utf-8'))
        return ('', hash.hexdigest())
    except Exception as e:
        logging.Logger('critical').critical(e, exc_info=True)
        return ('error of checksum generation', '')
