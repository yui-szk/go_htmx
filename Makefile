# Makefile for GTH Stack Demo

# 生成するバイナリ名
BINARY_NAME=app

# デフォルトの動作（make と打つだけで実行される）
.PHONY: all
all: run

# 1. 依存ツールのインストール（初回のみ使います）
.PHONY: setup
setup:
	go install github.com/a-h/templ/cmd/templ@latest
	go mod tidy

# 2. Templファイルの生成（.templ -> .go）
.PHONY: generate
generate:
	templ generate

# 3. アプリの実行（生成してから実行する完全版）
.PHONY: run
run: generate
	@echo "========================================"
	@echo "  GTH App is running at port 8080"
	@echo "========================================"
	go run .

# 4. ngrokの起動（別ターミナルで実行用）
.PHONY: ngrok
ngrok:
	ngrok http 8080

# 5. お掃除（生成されたファイルやバイナリを消す）
.PHONY: clean
clean:
	go clean
	rm -f *_templ.go
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe