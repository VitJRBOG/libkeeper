from django.urls import path

from . import views

urlpatterns = [
    path('add', views.add),
    path('update', views.update),
    path('get-all', views.get_all),
    path('get-versions', views.get_versions),
]
