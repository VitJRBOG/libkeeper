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
