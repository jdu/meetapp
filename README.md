MeetApp
==============

[![Circle CI](https://circleci.com/gh/shumipro/meetapp.svg?style=svg)](https://circleci.com/gh/shumipro/meetapp)

## ローカル環境

### Go環境ない人向けのGo環境構築（Mac & brew版）

```
$ brew update
$ brew install go
$ export GOPATH=$HOME
$ export PATH=$GOPATH/bin:$PATH
```

exportはbashなりzshに追記してください。

### install & build

```
$ go get -u github.com/shumipro/meetapp
$ cd $GOPATH/src/github.com/shumipro/meetapp
$ go build
$ ./meetapp
```

godepがinstallされていない場合 

```
$ go get github.com/kr/godep
$ godep get
```

### 環境変数

ローカルで動かす場合に、環境変数に以下の指定が必要です。

※FakeLogin的なのを実装して、いずれ指定無しで動くようにします...

- `FACEBOOK_APPID`: 開発用に作成したFB AppのAppID
- `FACEBOOK_SECRET`: 開発用に作成したFB AppのSECRET
- `BASE_URL`: FB AuthのcallbackでアクセスするURL (localhostはダメだったと思うのでhostsで適当にドメインを指定)
- `CLOUDINARY_URL`: 開発用に取得したcloudinaryのURL (アカウントのページに表示されているcloudinary://から始まるURL)
- (optional) `AIRBRAKE_PROJECT_ID`: 開発用に取得したAirbrakeのID
- (optional) `AIRBRAKE_API_KEY`: 開発用に取得したAirbrakeのkey
- (optional) `STATIC_URL`: staticリソース(imgなど)のベースURL (開発時はブランクでok)

### try

[http://localhost:8000](http://localhost:8000)

## heroku環境

自分用のheroku作りたい人向け

```
$ heroku create -b https://github.com/kr/heroku-buildpack-go.git meetapp-xxx
$ git push heroku master
```

- `xxx`: は適当にかぶらない文字列（いっそランダムでもいいけど）

## Front-end Dev Setup

JS/CSS(stylus)を変更する場合は以下の手順（nodeがinstallされている前提）

### install gulp

```
npm install -g gulp
```

### run npm install
meetappディレクトリで一度以下を実行
```
npm install
```

### run gulp
JS/CSSを変更した場合watchしてbuildするタスク
```
gulp
```

development用にJSをuncompressedする場合
```
gulp --dev
```
