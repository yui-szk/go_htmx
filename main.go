package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/a-h/templ"
)

// 簡易的なデータストア（メモリ保存）
var (
	posts  []Post
	nextID = 1
	mu     sync.Mutex // 並列アクセス対策
)

func main() {
	// ルーティング
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/posts", handlePostList)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// ホーム画面の表示
func handleHome(w http.ResponseWriter, r *http.Request) {
	// データ取得（新しい順にするため逆順で渡すなどの工夫も可）
	mu.Lock()
	currentPosts := make([]Post, len(posts))
	copy(currentPosts, posts)
	// 表示用に逆順にする（新しいものが上）
	for i, j := 0, len(currentPosts)-1; i < j; i, j = i+1, j-1 {
		currentPosts[i], currentPosts[j] = currentPosts[j], currentPosts[i]
	}
	mu.Unlock()

	// ページ全体を描画
	templ.Handler(Page(currentPosts)).ServeHTTP(w, r)
}

// 投稿を受け取るハンドラ
func handlePost(w http.ResponseWriter, r *http.Request) {
	// フォーム解析
	r.ParseForm()
	name := r.FormValue("name")
	content := r.FormValue("content")

	// データ保存
	mu.Lock()
	newPost := Post{
		ID:        nextID,
		Name:      name,
		Content:   content,
		CreatedAt: time.Now(),
	}
	posts = append(posts, newPost)
	nextID++
	mu.Unlock()
	
	// ページ全体ではなく、新しくできた「PostItem」ひとつだけをHTMLとして返す
	component := PostItem(newPost)
	component.Render(r.Context(), w)
}

// 最新のリスト(HTML)だけを返す関数
func handlePostList(w http.ResponseWriter, r *http.Request) {
    // データ取得（並列アクセス対策のロックを忘れずに）
    mu.Lock()
    currentPosts := make([]Post, len(posts))
    copy(currentPosts, posts)
    // 逆順にする処理
    for i, j := 0, len(currentPosts)-1; i < j; i, j = i+1, j-1 {
        currentPosts[i], currentPosts[j] = currentPosts[j], currentPosts[i]
    }
    mu.Unlock()

    // リストの中身だけをレンダリングして返す
    for _, post := range currentPosts {
        // PostItemコンポーネントを直接書き込む
        PostItem(post).Render(r.Context(), w)
    }
}