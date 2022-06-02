from django.urls import path

from . import views

urlpatterns = [
    path('add', views.add),
    path('update', views.update),
    path('get-notes', views.get_notes),
    path('get-versions', views.get_versions),
]
