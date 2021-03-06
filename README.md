## アプリの概要
オンラインでユーザー同士が連想しながら物語を作り上げるアプリケーション

## アプリを作ろうと思った理由
何かアプリを開発したいと考えた時、自分の好きな事や理由を書き出してみた結果「他人の考えている事に興味がある」なんじゃないかという事が分かったので連想しながら物語を作り上げるアプリを開発しました。

## アプリの機能一覧
- ユーザー登録機能
- ユーザー認証機能
- ルーム作成機能
- ランキング機能
- 投票機能

## データベース
- データベースはMySQLを使用しました。
### DataBaseの構成
DataBase名<br >
- STORY
***
Table名<br >
- USER<br >

|Field|Type|Null|Key|Default|Extra|
|:--------|:---------|:--------|:-------|:-------|:-------|
|name|varchar(15)|NO|PRI|NULL|
|password|varchar(32)|YES||NULL|
|mail|varchar(256)|YES||NULL|

- ROOM<br >

|Field|Type|Null|Key|Default|Extra|
|:--------|:---------|:--------|:-------|:-------|:-------|
|room_name|varchar(20)|NO|PRI|NULL|
|title|varchar(20)|NO||NULL|
|create_user|varchar(15)|NO||NULL|
|create_day|datetime|YES||NULL|

- TALKED<br >

|Field|Type|Null|Key|Default|Extra|
|:--------|:---------|:--------|:-------|:-------|:-------|
|room_name|varchar(20)|NO|MUL|NULL|
|user_name|varchar(15)|NO|MUL|NULL|
|talk_word|varchar(30)|NO||NULL|
|talked_time|datetime|NO||NULL|

- VOTING<br >

|Field|Type|Null|Key|Default|Extra|
|:--------|:---------|:--------|:-------|:-------|:-------|
|room_name|varchar(20)|NO|MUL|NULL|
|user_name|varchar(15)|NO|MUL|NULL|

***

## アプリのバックエンド開発で苦労した点
GO言語も初めて触ったので、少し苦労しました。
GO言語公式のチュートリアルで基礎を勉強してから開発に取り組みました。

<br >
