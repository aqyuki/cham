# README

Discord 上で送信されたメッセージリンクを展開します．

## How to use

> [!NOTE]
> この Bot は Docker 上で動作させることを想定しています．

事前に，Discord Token を発行し，Bot をサーバーに招待してください．

1. このリポジトリをローカルにクローンする

2. `.env.sample`を`.env`にファイル名を変更する.

3. `.env`の**DISCORD_TOKEN**をあなたの Discord Token に変更する．

4. `docker compose up -d`で起動する．

## Configuration

| key             | 概要                                              |
| :-------------- | :------------------------------------------------ |
| `DISCORD_TOKEN` | Discord の API を使用するためのトークン           |
| `LOG_MODE`      | `develop`とするとログをテキスト形式で出力します． |
| `LOG_LEVEL`     | 出力するログのレベルを指定します．                |
