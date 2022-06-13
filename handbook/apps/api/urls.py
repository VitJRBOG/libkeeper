from django.urls import path

from . import views

urlpatterns = [
    path('add-note', views.add_note),
    path('add-version', views.add_version),
    path('get-notes', views.get_notes),
    path('get-versions', views.get_versions),
]
