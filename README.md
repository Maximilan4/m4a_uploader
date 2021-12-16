### m4a uploader
загружает треки m4a в s3 хранилище по ключу {isrc}/track.m4a

### Алгоритм:
- загружаем базу искомых пар amid и isrc (см source.example.csv, можно указать флагом -s)
- сканируем переданную директорию (флаг -d)
- вычленяем из тега трека имя автора и название
- на каждый трек с тегом сделается запрос в apple, чтобы соотнести результаты поиска и пару в бд
- если соответствие найдено - загружаем файл с диска на s3 (см aws_config.example.json и флаг -a)

### run
```bash
cat aws_config.example.json > aws_config.json # создать конфиг и заменить все значения на свои
go build -o uploader cmd/m4a/main.go
./uploader -d /path/to/m4a/files # при условии что source.csv и aws_config.json лежат в директории исполенения)
```