from django.contrib import admin

from .models import Note, Version


admin.site.register(Note)
admin.site.register(Version)
