# とりあえずgolangのバージョンは最新を使用
FROM golang:latest
# アップデートとgitのインストール
RUN apt-get update && apt-get install -y git
WORKDIR /go/src