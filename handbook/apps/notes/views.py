import datetime
import time
import hashlib

from tools import logging
from django.http import HttpRequest, JsonResponse
from . import models


def add(request: HttpRequest) -> JsonResponse:
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
        logging.Logger('critical').critical(e)
        return JsonResponse({
            'status_code': 500,
            'error': 'note creation error',
        })

    note_id = models.Note.objects.last().id  # type: ignore

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
        logging.Logger('critical').critical(e)
        note.delete(note_id)
        return JsonResponse({
            'status_code': 500,
            'error': 'note creation error',
        })

    return JsonResponse({
        'status_code': 200
    })


def update(request: HttpRequest) -> JsonResponse:
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

    try:
        note = models.Note.objects.get(id=id_)
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
        version.create(text, date, checksum, id_)  # type: ignore
    except Exception as e:
        logging.Logger('critical').critical(e)
        return JsonResponse({
            'status_code': 500,
            'error': 'note updation error',
        })

    try:
        note.title = title
        note.save()
    except Exception as e:
        logging.Logger('critical').critical(e)
        version.delete(id_)
        return JsonResponse({
            'status_code': 500,
            'error': 'note updation error',
        })

    return JsonResponse({
        'status_code': 200,
    })


def get_all(request: HttpRequest) -> JsonResponse:
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
        logging.Logger('critical').critical(e)
        return JsonResponse({
            'status_code': 500,
            'error': 'note selection error',
        })


def get_by_id(request: HttpRequest) -> JsonResponse:
    if request.method != 'GET':
        return JsonResponse({
            'status_code': 400,
            'error': f'{request.method} is not allowed.',
        })

    id_ = request.GET.get('id')

    try:
        if id_ is not None:
            queryset = models.Version.objects.filter(note_id=id_).order_by('date').reverse().values()

            return JsonResponse({
                'status_code': 200,
                'items': list(queryset)[0],
            })
        else:
            return JsonResponse({
                'status_code': 400,
                'error': '"id" attribute is required',
            })
    except Exception as e:
        logging.Logger('critical').critical(e)
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
        logging.Logger('critical').critical(e)
        return ('error of title composing', '')

def __get_date() -> tuple[str, int]:
    try:
        today = datetime.datetime.now()
        
        return ('', int(time.mktime(today.timetuple())))
    except Exception as e:
        logging.Logger('critical').critical(e)
        return ('error of getting date', 0)

def __gen_checksum(text: str) -> tuple[str, str]:
    try:
        hash = hashlib.md5(text.encode('utf-8'))
        return ('', hash.hexdigest())
    except Exception as e:
        logging.Logger('critical').critical(e)
        return ('error of checksum generation', '')
