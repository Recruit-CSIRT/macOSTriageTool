# macOS Triage Tool

macOS向けのアーティファクト収集ツール

<img src="imgs/mtt_cli.gif" width="700">
<img src="imgs/mtt_gui.gif" width="700">

## 概要

`macOS Triage Tool`(以下、`MTT`)は、macOSのアーティファクトを収集するフォレンジックツールです。
このツールは、ディレクトリ構造を保持したまま、調査に必要なファイルだけを収集することを目的としています。

`MTT`はCLI版だけでなくGUI版もあり、初めての人でも簡単に利用できるようになっています。
インストール作業等は不要で、1つのバイナリもしくはAppを対象において実行するだけでアーティファクトの収集を開始することができます。


## 特徴
- GUI版とCLI版
    - リモートで使用したい場合はCLI版を選択するなど、ユースケースに応じて利用できます。

- ディレクトリ階層を保持
    - 既存のディレクトリ構造を保持したままファイルを取得するため、様々なパースツールで取得結果を処理することが可能です。
    
- Extended attributeを保持
    - 書き出し先のデバイスがExtended attributeに対応していなくてもDisk Image(dmg)形式で保存できるため、出力先のファイルシステムに依存しません。

- 4つのプリセット
    - 以下の4つの状況に応じて、取得するアーティファクトを選択できます。
        - `AllList`: ツールで定義しているすべてのアーティファクトを収集
        - `Malware`: マルウェアの調査で取得した方が良いアーティファクトを収集
        - `Fraud`: 内部不斉の調査で取得した方が良いアーティファクトを収集
        - `macripper`: 解析ツール(macripper)で対応しているアーティファクトを収集
    - また、自身で定義したファイルも読み込んでアーティファクトを収集することも可能です。
    
## クイックスタート
デフォルト値が設定されているため、何も指定せずに実行することも可能です。
ツールは、Mojave、Catalinaで動作します。


### CLI
CLIの場合は、sudo権限で実行した方が取得できるアーティファクトが多くなります。  
デフォルトの場合、ルートディレクトリを基点として、ファイルを取得しはじめ、現在のディレクトリ配下に結果を出力します。  
アーティファクトは設定しているものすべてを取得します。
```
$ sudo ./mtt
```

### GUI
MTT.appをクリック、設定項目を入力し、Runボタンを押すと実行が開始されます。


## 使い方

### CLI
```
Usage of mtt:
  -c    Calc the file hash. (default: false)
  -d    Save files into a dmg. (default: false)
  -file string
        Set user custom file list path.
  -i    Get the system information. (default: false)
  -output string
        Set the output path. (default ".")
  -preset int
        Choose forensics type: 0(AllList), 1(Malware), 2(Fraud), 3(macripper), 4(only custom). (default: 0)
  -root string
        Set the evidence root path. (default "/")
  -s    Get the stat info. (default: false)
```

例: 
```
$ sudo ./mtt -root / -output /Volumes/USB -preset 2 -file /Volumes/USB/udf.txt -s -c -i -d
```

### GUI

![gui](imgs/gui.png)

- Root path setting: 取得対象のルートディレクトリを指定
- Output path setting: アーティファクト保存先を指定
- Preset setting:
    - Preset: AllList, Malware, Fraud, macripper, only custom から選択
    - Select your file list: 自分で定義したカスタムファイルリストを選択
- Options:
    - Get system information: macOSのシステム情報を取得
    - Get file hash: アーティファクトのhash(md5)を取得
    - Get stat info: アーティファクトのstat情報(最終アクセス日時、作成日時など)を取得
    - Save file into dmg: Disk Imageにアーティファクトを保存
    
### カスタムファイルリストの作成方法
テキストファイルに取得したいファイルを1行ずつ記入してください。  
ユーザ名などの固定できないパスには、`*`を使用することができます。

例: udf.txt
```
/Users/*/Downloads/snork*
/Users/*/Documents/
/Library/LaunchAgents/*.plist
```

### Tips
取得対象ファイルを Disk Image(dmg 形式)に保存する場合は、 「Options」の
「Save file into dmg」にチェックを入れます。 手元の検証では、このオプションを使用すると、より多くのディスク容量を必要とする代わりに、
より早くコピーを終えることができます。 コピー先のファイルシステムがmacOS の extended attributes に対応して いない場合、diito によるコピーに失敗する場合かがありますが、
dmg オプションを使用すると取得できることがあります。

## コンパイル

### qtのインストール
```
$ brew install qt
$ export QT_HOMEBREW=true
$ export GO111MODULE=off; xcode-select --install; go get -v github.com/therecipe/qt/cmd/... && $(go env GOPATH)/bin/qtsetup -test=false
```
### CLI
```
# コンパイル
$ make cli
```

### GUI
```
$ make gui
```

## License
This repository is available under the GNU General Public License v3.0

## Author
moniik

