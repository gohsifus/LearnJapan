package yandexTranslateApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	token = "Bearer  t1.9euelZrNj5ieyZfPzJiNk4-Sj5WQnO3rnpWaxovGjJ3GzJTKyovGnpKKkszl9Pc6R31q-e9GagfY3fT3enV6avnvRmoH2A.xbMm64cJkpML_OzOGa6P7MKmQo3FAGqGYvMKG3owP5T20nGWQ-l-l7rB7-TW6lZKOd5EISO-lwF51J_MV7eHAw"
	folderId = "b1gfhitl9di88phga4fk"
	globalUrl = "https://translate.api.cloud.yandex.net/"
)

//callApiMethod Отправляет запрос с параметрами на конкретный endpoint
func callApiMethod(params map[string]interface{}, method string, url string) (*http.Response, error){
	url = globalUrl + url
	client := http.Client{}
	if method == "POST"{
		jsonParams, err := json.Marshal(params)
		if err != nil{
			return &http.Response{}, errors.New("Ошибка кодирования параметров в json")
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
		if err != nil{
			return &http.Response{}, errors.New("Ошибка создания запроса")
		}
		req.Header.Set("Authorization", token)

		response, err := client.Do(req)
		if err != nil{
			return &http.Response{}, err
		}

		return response, nil
	} else if method == "GET"{

	}
	return &http.Response{}, errors.New("Передан некорректный аргумент method, ожидалось: POST или GET")
}

type Language struct{
	Name string `json:"name"`
	Code string `json:"code"`
}

func (l Language) String() string{
	return fmt.Sprintf("code: %s, name: %s", l.Code, l.Name)
}

//GetLanguageList Вернет список поддерживаемых языков
func GetLanguageList() ([]Language, error){
	url := "translate/v2/languages"

	params := map[string]interface{}{
		"folderId": folderId,
	}
	response, err := callApiMethod(params, "POST", url)
	if err != nil{
		return []Language{}, err
	}

	languages := make(map[string][]Language)
	bytes, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(bytes, &languages)

	return languages["languages"], nil
}

//Translate Переведет строку
func Translate(str, srcLang, dstLang string) (string, error){
	url := "translate/v2/translate"
	params := map[string]interface{}{
		"folderId": folderId,
		"sourceLanguageCode": srcLang,
		"targetLanguageCode": dstLang,
		"texts": []string{str},
	}

	response, err := callApiMethod(params, "POST", url)
	if err != nil{
		return "", err
	}

	result := make(map[string][]map[string]string)
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	json.Unmarshal(body, &result)

	return result["translations"][0]["text"], nil
}