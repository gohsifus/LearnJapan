package yandexTranslateApi

import (
	_ "LearnJapan.com/configs"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type yApiConf struct {
	folderId   string `json:"folder_id"`
	globalUrl  string `json:"global_url"`
	oAuthToken string `json:"o_auth_token"`
	iamToken   string `json:"iam_token"`
}

func (y *yApiConf) convertFromMap(m map[string]string) {
	y.folderId = m["folder_id"]
	y.globalUrl = m["global_url"]
	y.oAuthToken = m["o_auth_token"]
	y.iamToken = m["iam_token"]
}

var apiConf yApiConf

func init() {
	/*data := map[string]string{
		"folder_id":    configs.Cfg.YandexApiFolderId,
		"global_url":   configs.Cfg.YandexApiUrl,
		"o_auth_token": configs.Cfg.YandexApiToken,
	}

	apiConf.convertFromMap(data)
	updateIamToken()*/
	//fmt.Println("apiYandex")
}

func updateIamToken() {
	client := http.Client{}

	data := map[string]string{
		"yandexPassportOauthToken": apiConf.oAuthToken,
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Println("Ошибка обновления iam токена")
	}

	request, err := http.NewRequest("POST", "https://iam.api.cloud.yandex.net/iam/v1/tokens", bytes.NewBuffer(body))
	response, _ := client.Do(request)

	reqBytes, _ := ioutil.ReadAll(response.Body)
	iamToken := make(map[string]string)
	json.Unmarshal(reqBytes, &iamToken)
	apiConf.iamToken = iamToken["iamToken"]
	fmt.Println(iamToken["iamToken"])
}

// callApiMethod Отправляет запрос с параметрами на конкретный endpoint
func callApiMethod(params map[string]interface{}, method string, url string) (*http.Response, error) {
	url = apiConf.globalUrl + url
	client := http.Client{}
	if method == "POST" {
		jsonParams, err := json.Marshal(params)
		if err != nil {
			return &http.Response{}, errors.New("Ошибка кодирования параметров в json")
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
		if err != nil {
			return &http.Response{}, errors.New("Ошибка создания запроса")
		}
		req.Header.Set("Authorization", "Bearer "+apiConf.iamToken)

		response, err := client.Do(req)
		if err != nil {
			return &http.Response{}, err
		}

		return response, nil
	} else if method == "GET" {

	}
	return &http.Response{}, errors.New("Передан некорректный аргумент method, ожидалось: POST или GET")
}

type Language struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (l Language) String() string {
	return fmt.Sprintf("code: %s, name: %s", l.Code, l.Name)
}

// GetLanguageList Вернет список поддерживаемых языков
func GetLanguageList() ([]Language, error) {
	url := "translate/v2/languages"

	params := map[string]interface{}{
		"folderId": apiConf.folderId,
	}
	response, err := callApiMethod(params, "POST", url)
	if err != nil {
		return []Language{}, err
	}

	languages := make(map[string][]Language)
	bytes, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(bytes, &languages)

	return languages["languages"], nil
}

// Translate Переведет строку
func Translate(str, srcLang, dstLang string) (string, error) {
	url := "translate/v2/translate"
	params := map[string]interface{}{
		"folderId":           apiConf.folderId,
		"sourceLanguageCode": srcLang,
		"targetLanguageCode": dstLang,
		"texts":              []string{str},
	}

	response, err := callApiMethod(params, "POST", url)
	if err != nil {
		return "", err
	}

	result := make(map[string][]map[string]string)
	body, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	json.Unmarshal(body, &result)

	fmt.Println(result)

	return result["translations"][0]["text"], nil
}
