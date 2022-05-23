import datetime
import time
import hashlib

from django.http import HttpRequest, JsonResponse
from . import models


def add(request: HttpRequest) -> JsonResponse:
    if request.method != 'POST':
        return JsonResponse({
            'status_code': 400,
            'error': f'{request.method} is not allowed.',
        })

    if request.POST.get('text') is None:
        return JsonResponse({
            'status_code': 400,
            'error': 'attribute "text" is required',
        })
    text = str(request.POST.get('text'))
    title = __compose_title(text)
    date = __get_date()

    note = models.Note()
    note.create(title, date)

    note_id = models.Note.objects.last().id

    checksum = __gen_checksum(text)

    version = models.Version()
    version.create(text, date, checksum, note_id)

    return JsonResponse({
        'status_code': 200
    })

def __compose_title(text: str) -> str:
    if len(text) > 47:
        if '\n' in text[:50]:
            return text[:text.find('\n')]

        return f'{text[:46]}...'
    return text

def __get_date() -> int:
    today = datetime.datetime.now()
    
    return int(time.mktime(today.timetuple()))

def __gen_checksum(text: str) -> str:
    hash = hashlib.md5(text.encode('utf-8'))
    return hash.hexdigest()
