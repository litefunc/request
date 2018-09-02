package request

import (
	"bytes"
	"cloud/lib/logger"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func req(method, url string, header map[string]string, body map[string]interface{}, obj interface{}) error {
	logger.Debug(url)
	//body
	jData, err := json.Marshal(body)
	if err != nil {
		logger.Error(err)
		return err
	}
	jBytes := bytes.NewReader(jData)
	//req
	req, err := http.NewRequest(method, url, jBytes)
	if err != nil {
		logger.Error(err)
		return err
	}
	//header
	for key, val := range header {
		req.Header.Set(key, val)
	}

	client := &http.Client{}
	//resp-status
	resp, err := client.Do(req)
	if err != nil {
		if resp != nil {
			logger.Warn(resp)
		}
		logger.Error(err)
		return err
	}
	if resp.StatusCode == 404 {
		logger.Error("404 page not found")
		return errors.New("404 page not found")
	}

	//resp-body
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Debug(string(msg))
	//resp-body-Unmarshal
	if obj != nil {
		if err := json.Unmarshal(msg, obj); err != nil {
			logger.Error(err)
			return err
		}
	} else {
		var i interface{}
		if err := json.Unmarshal(msg, &i); err != nil {
			logger.Error(err)
			return err
		}
		if i != nil {
			logger.Error("msg != nil", i)
			return errors.New("msg != nil")
		}
	}

	return nil
}

func Get(url string, header map[string]string, obj interface{}) error {
	return req("GET", url, header, nil, obj)
}

func Post(url string, header map[string]string, body map[string]interface{}, obj interface{}) error {
	return req("POST", url, header, body, obj)
}

func Put(url string, header map[string]string, body map[string]interface{}, obj interface{}) error {
	return req("PUT", url, header, body, obj)
}

func Delete(url string, header map[string]string, body map[string]interface{}, obj interface{}) error {
	return req("DELETE", url, header, body, obj)
}

type FormFile struct {
	Field string
	Path  string
}

func reqForm(method, url string, header map[string]string, params map[string]string, formfile FormFile, obj interface{}) error {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.Println(url)

	body := bytes.Buffer{}

	//body - create writer
	writer := multipart.NewWriter(&body)

	//body-params
	for key, val := range params {
		if err := writer.WriteField(key, val); err != nil {
			log.Println(err)
			return err
		}
	}
	//body-formfile
	if formfile != (FormFile{}) {
		file, err := os.Open(formfile.Path)
		if err != nil {
			log.Println(err)
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(formfile.Field, filepath.Base(formfile.Path))
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	//body - close writer
	if err := writer.Close(); err != nil {
		log.Println(err)
		return err
	}
	//req
	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		log.Println(err)
		return err
	}
	//header
	for key, val := range header {
		req.Header.Set(key, val)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	//resp-status
	resp, err := client.Do(req)
	if err != nil {
		log.Println(resp.StatusCode, err)
		return err
	}
	if resp.StatusCode == 404 {
		return errors.New("404 page not found")
	}

	//resp-body
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	logger.Debug(string(msg))
	//resp-body-Unmarshal
	if obj != nil {
		if err := json.Unmarshal(msg, obj); err != nil {
			log.Println(err)
			return err
		}
	} else {
		var i interface{}
		if err := json.Unmarshal(msg, &i); err != nil {
			log.Println(err)
			return err
		}
		if i != nil {
			log.Println("msg != nil", i)
			return errors.New("msg != nil")
		}
	}

	return nil
}

func PostForm(url string, header map[string]string, params map[string]string, obj interface{}) error {
	return reqForm("POST", url, header, params, FormFile{}, obj)
}

func PutForm(url string, header map[string]string, params map[string]string, obj interface{}) error {
	return reqForm("PUT", url, header, params, FormFile{}, obj)
}

func DeleteForm(url string, header map[string]string, params map[string]string, obj interface{}) error {
	return reqForm("DELETE", url, header, params, FormFile{}, obj)
}

func PostFormFile(url string, header map[string]string, params map[string]string, formfile FormFile, obj interface{}) error {
	return reqForm("POST", url, header, params, formfile, obj)
}

func PutFormFile(url string, header map[string]string, params map[string]string, formfile FormFile, obj interface{}) error {
	return reqForm("PUT", url, header, params, formfile, obj)
}
