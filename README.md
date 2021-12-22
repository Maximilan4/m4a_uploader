### m4a uploader
uploads m4a tracks to s3 by key {isrc}/track.m4a

### Алгоритм:
- loads search base of amid и isrc (look at source.example.csv, -s flag)
- scan dirs with file (-d flag)
- read file meta
- check tracks data in apple, and when search in loaded source file
- if match found -> uploads file to s3

### run
```bash
cat aws_config.example.json > aws_config.json
go build -o uploader cmd/m4a/main.go
./uploader -d /path/to/m4a/files # source.csv and aws_config.json in same dir
```