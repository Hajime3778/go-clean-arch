# 環境構築

## 初期設定
```
. scripts/setup.sh
```
## ローカル環境起動

```
docker-compose up --build
```

## 終了

```
docker-compose down -v
```

# ディレクトリ構成



## domain

ドメイン層、どの層からもアクセスされる共通のモデルを配置

## infrastructure

フレームワーク、ライブラリに依存関係がある場合、こちらにinterfaceを作成

## interface

## usecase







