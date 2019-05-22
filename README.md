# cron-file-reader-via-http

```github.com/pkg/errors```, ```context.WithCancel``` 使いたかっただけ

- 指定したファイルを http で返す
- 定期的に指定したファイルの中身を取得

## 使い方

```
$ go run main.go -v sample.txt
```

```
$ curl localhost:8080/
helloworld
```
