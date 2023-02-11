package traefik_jwt_internal

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"
	"strings"
	"text/template"
	"time"
)

type JWT struct {
	next           http.Handler
	name           string
	config         *Config
	claimsTemplate *template.Template
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	j := &JWT{
		config: config,
		next:   next,
		name:   name,
	}

	t, err := template.New("").Funcs(template.FuncMap{
		"unmarshalJson": func(s string) interface{} {
			var data interface{}

			json.Unmarshal([]byte(s), &data)

			return data
		},
	}).Parse(j.config.Claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse claims template")
	}

	j.claimsTemplate = t

	return j, nil
}

func (j *JWT) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	header := map[string]interface{}{
		"alg": j.config.Alg,
		"typ": "JWT",
	}

	headerJson, err := json.Marshal(header)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload := make(map[string]interface{})
	payload["iat"] = time.Now().Unix()
	payload["exp"] = time.Now().Unix() + j.config.TTL

	data := make(map[string]interface{})
	data["Headers"] = req.Header

	writer := &bytes.Buffer{}

	err = j.claimsTemplate.Execute(writer, req)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var claims map[string]interface{}
	err = json.Unmarshal(writer.Bytes(), &claims)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload = MergeJSONMaps(claims, payload)

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := j.MakeJWT(headerJson, payloadJson)

	headerValue := strings.Join([]string{j.config.Header.Prefix, token}, " ")
	req.Header.Add(j.config.Header.Name, headerValue)

	j.next.ServeHTTP(rw, req)
}

func (j *JWT) MakeJWT(header []byte, payload []byte) string {
	message := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)

	var h hash.Hash
	switch j.config.Alg {
	case "HS512":
		h = hmac.New(sha512.New, []byte(j.config.Secret))
	case "HS256":
		fallthrough
	default:
		h = hmac.New(sha256.New, []byte(j.config.Secret))
	}

	h.Write([]byte(message))

	return message + "." + base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
