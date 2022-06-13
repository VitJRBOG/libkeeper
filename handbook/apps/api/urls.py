from django.urls import path

from . import views

urlpatterns = [
    path('add-note', views.add_note),
    path('update-note', views.update_note),
    path('get-notes', views.get_notes),
    path('get-versions', views.get_versions),
]
