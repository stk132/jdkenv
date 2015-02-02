# jdkenv
jdkenv is JAVA_HOME change command written in golang.

## description

JAVA_HOMEを簡単に切り替えるためのコマンドです。
JAVA_HOMEにシンボリックリンクを指定して、シンボリックリンクの向き先を変更します。

## useage

### 初期化

```console
$ jdkenv init
```

実行後、環境変数JAVA_HOMEに~/.jdkenv/java/currentを指定し、JAVA_HOME/binにパスを通す

### jdkのリスト表示

windowsの場合は~/.jdkenv/java/配下にある、プリフィクスがjdkのディレクトリを表示する
macの場合は以下から探す

- "/System/Library/Java/JavaVirtualMachines/"
- "/Library/Java/JavaVirtualMachines/"

```console
$ jdkenv list
```

### JAVA_HOMEの変更

シンボリックリンク~/.jdkenv/java/currentの向き先を変更する

```console
$ jdkenv use "listコマンドで表示されたjdk番号"
```

### 現在のjdkの確認

```console
$ jdkenv current
```

### pecoと組み合わせる

```console
$ jdkenv use `jdkenv list | peco`
```

## TODO

- init時のメッセージが適当なのを修正する　
- helpメッセージの実装
- このREADMEの英語化
