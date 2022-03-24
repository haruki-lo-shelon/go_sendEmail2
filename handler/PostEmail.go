package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/joho/godotenv"

	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

//メールの構造体
type Email struct {
	// Name		string
	// Subject	string
	Tokumemo_API_KEY	string	//追加
	Text    string
	Email   string
}

//送信する関数
func PostEmail(w rest.ResponseWriter, r *rest.Request) {
	SendToEmail := Email{}
	err := r.DecodeJsonPayload(&SendToEmail)//sendにpost値を入れる
	if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	//送信元のメールアドレスが空ならエラー
	if SendToEmail.Email == "" {
			rest.Error(w, "mail required", 400)
			return
	}

	err_read := godotenv.Load()
	if err_read != nil {
		log.Fatalf("error: %v", err_read)
	}

	// .envから環境変数読み込み
	SendGrid_API_KEY := os.Getenv("SendGrid_API_KEY")
	TOS := strings.Split(os.Getenv("TOS"), ",")
	fr := SendToEmail.Email//postした値に含めたメールアドレス（送信者のメアド）
	Tokumemo_API_KEY := os.Getenv("Tokumemo_API_KEY")

	//TokumemoのAPIキーを検証
	if SendToEmail.Tokumemo_API_KEY != Tokumemo_API_KEY {
		rest.Error(w, "correct apiKey required", 400)
		return
	} else {
		//検証が終わったら隠す
		SendToEmail.Tokumemo_API_KEY = "xxxxxxxxx"
	}

	// メッセージの構築
	message := mail.NewV3Mail()
	// 送信元を設定
	from := mail.NewEmail("", fr)
	message.SetFrom(from)

	// 宛先と対応するSubstitutionタグを指定(宛先は複数指定可能)
	EmailDestination := mail.NewPersonalization()
	to := mail.NewEmail("", TOS[0])
	EmailDestination.AddTos(to)

	//残りのpostされた値を変数に格納
	// name := SendToEmail.Name
	// subject := SendToEmail.Subject
	text := SendToEmail.Text

	// EmailDestination.SetSubstitution("%name%", name)
	// EmailDestination.SetSubstitution("%m_subject%", subject)
	EmailDestination.SetSubstitution("%m_text%", text)
	message.AddPersonalizations(EmailDestination)

	// 2つ目の宛先と、対応するSubstitutionタグを指定
	EmailDestination2 := mail.NewPersonalization()
	to2 := mail.NewEmail("", TOS[1])
	EmailDestination2.AddTos(to2)
	// EmailDestination2.SetSubstitution("%name%", name)
	// EmailDestination2.SetSubstitution("%m_subject%", subject)
	EmailDestination2.SetSubstitution("%m_text%", text)
	message.AddPersonalizations(EmailDestination2)

	// 3つ目の宛先と、対応するSubstitutionタグを指定
	// EmailDestination3 := mail.NewPersonalization()
	// to3 := mail.NewEmail("", TOS[2])
	// EmailDestination3.AddTos(to3)
	// EmailDestination3.SetSubstitution("%name%", name)
	// EmailDestination3.SetSubstitution("%m_subject%", subject)
	// EmailDestination3.SetSubstitution("%m_text%", text)
	// message.AddPersonalizations(EmailDestination3)

	// 件名を設定
	message.Subject = "ユーザーからのお問い合わせ"
	// テキストパートを設定
	c := mail.NewContent("text/plain", "%m_text%\r\n")
	message.AddContent(c)
	// HTMLパートを設定
	//c = mail.NewContent("text/html", "<strong> %name% さんは何をしていますか？</strong><br>　文章ー－－。")
	// message.AddContent(c)

	// カテゴリ情報を付加
	// message.AddCategories("category1")
	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")
	// 画像ファイルを添付
	// a := mail.NewAttachment()
	// file, _ := os.OpenFile("./gif.gif", os.O_RDONLY, 0600)
	// defer file.Close()
	// data, _ := ioutil.ReadAll(file)
	// data_enc := base64.StdEncoding.EncodeToString(data)
	// a.SetContent(data_enc)
	// a.SetType("image/gif")
	// a.SetFilename("owl.gif")
	// a.SetDisposition("attachment")
	// message.AddAttachment(a)

	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(SendGrid_API_KEY)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	w.WriteJson(&SendToEmail)
}
