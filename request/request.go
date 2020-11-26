package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	base    string
	headers Headers
}
type request struct {
	url     string
	method  string
	body    interface{}
	headers Headers
}
type Headers map[string]string

func New(url string, headers ...Headers) *Request {
	h := Headers{}
	if len(headers) > 0 {
		h = headers[0]
	}
	return &Request{base: url, headers: h}
}

// GET hacer petición
func (r *Request) GET(url string, headers Headers, output interface{}) error {
	if headers != nil {
		for k, v := range headers {
			r.headers[k] = v
		}
	}
	req := &request{
		url:     r.base + url,
		method:  "GET",
		headers: r.headers,
	}
	return req.execute(output)
}

// POST enviar datos
func (r *Request) POST(url string, headers Headers, body, output interface{}) error {
	if headers != nil {
		for k, v := range headers {
			r.headers[k] = v
		}
	}
	req := &request{
		url:     r.base + url,
		method:  "POST",
		headers: r.headers,
		body:    body,
	}
	return req.execute(output)
}

// PUT .
func (r *Request) PUT(url string, headers Headers, body, output interface{}) error {
	if headers != nil {
		for k, v := range headers {
			r.headers[k] = v
		}
	}
	req := &request{
		url:     r.base + url,
		method:  "PUT",
		headers: r.headers,
		body:    body,
	}
	return req.execute(output)
}

func (r *request) execute(output interface{}) error {
	// Crear Cliente HTTP
	httpClient := &http.Client{}

	// Codificar el cuerpo
	var body = &bytes.Buffer{}
	if r.body != nil {
		jsonCode, err := json.Marshal(r.body)
		if err != nil {
			log.Println("Error en el marshal json", err)
			return err
		}
		body = bytes.NewBuffer(jsonCode)
	}

	// Crear la petición
	req, err := http.NewRequest(r.method, r.url, body)
	if err != nil {
		log.Println("Error creando petición", err)
		return err
	}

	for key, value := range r.headers {
		req.Header.Add(key, value)
	}
	if body != nil {
		req.Header.Add("Content-type", "Application/Json")
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error ejecutando petición", res)
		return err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error obteniendo respuesta", err)
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New(string(resBody))
	}
	err = json.Unmarshal(resBody, output)
	if err != nil {
		log.Println("Error decodificando respuesta", err)
		return err
	}
	return nil
}
