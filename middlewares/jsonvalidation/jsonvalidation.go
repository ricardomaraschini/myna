package jsonvalidation

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ricardomaraschini/myna/webserver"
	"github.com/xeipuuv/gojsonschema"
)

// New receives a json schema as parameter and returns a middleware that
// validates the input against the provided schema
func New(schemaContent []byte) webserver.MiddleWare {

	schema := gojsonschema.NewBytesLoader(schemaContent)

	return func(body []byte) ([]byte, error) {

		log.Printf("validating input against jsonschema")
		if len(body) == 0 {
			log.Printf("empty body")
			return nil, errors.New("empty body")
		}

		doc := gojsonschema.NewBytesLoader(body)
		res, err := gojsonschema.Validate(schema, doc)
		if err != nil {
			log.Printf("error validating json: %s", err.Error())
			return nil, err
		}

		if res.Valid() {
			return body, nil
		}

		// concat all errors into a map so we can encode it before
		// return
		errfound := make(map[string]string)
		for _, e := range res.Errors() {
			log.Printf(
				"invalid json field %s: %s\n",
				e.Field(),
				e.Description(),
			)
			errfound[e.Field()] = e.Description()
		}

		errbody, err := json.Marshal(errfound)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(errbody))
	}
}
