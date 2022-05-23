from django.db import models


class Note(models.Model):
    id = models.AutoField('ID', primary_key=True)
    title = models.CharField('Title', max_length=50)
    date = models.IntegerField('Creation date')

    def create(self, title: str, date: int):
        self.title = title
        self.date = date

        self.save()

    def delete(self, id_: int):
        self.objects.filter(id=id_).delete()

    def __str__(self) -> str:
        return self.title

    class Meta:
        verbose_name = 'Note'
        verbose_name_plural = 'Notes'


class Version(models.Model):
    id = models.AutoField('ID', primary_key=True)
    text = models.TextField('Text', max_length=1000)
    date = models.IntegerField('Creation date')
    checksum = models.CharField('Checksum', max_length=64)
    note_id = models.IntegerField('Note ID')

    def create(self, text: str, date: int, checksum: str, note_id: int):
        self.text = text
        self.date = date
        self.checksum = checksum
        self.note_id = note_id

        self.save()

    def __str__(self) -> str:
        return str(self.checksum)

    class Meta:
        verbose_name = 'Version'
        verbose_name_plural = 'Versions'
