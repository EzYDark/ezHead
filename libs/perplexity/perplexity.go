package perplexity

import (
	"encoding/json"
	"errors"

	"github.com/go-rod/rod"
	"github.com/rs/zerolog/log"
)

type Perplexity struct {
	RequestHeaders *RequestHeaders
	RequestBody    *RequestBody
	RequestScript  *RequestScript
}

var perplex *Perplexity

func Init() (*Perplexity, error) {
	if perplex != nil {
		return nil, errors.New("perplexity struct already initialized")
	}

	perplex = &Perplexity{
		RequestHeaders: new(RequestHeaders),
		RequestBody:    new(RequestBody),
		RequestScript:  new(RequestScript),
	}

	perplex.RequestHeaders = perplex.RequestHeaders.Default()
	perplex.RequestBody = perplex.RequestBody.Default()
	perplex.RequestScript = perplex.RequestScript.Update(perplex.RequestHeaders, perplex.RequestBody)

	return perplex, nil
}

func (p *Perplexity) SendRequest(page *rod.Page, query string) (Response, error) {
	p.RequestBody.QueryStr = query
	p.RequestScript = p.RequestScript.Update(p.RequestHeaders, p.RequestBody)

	// Wait for the page to fully load
	page.MustWaitStable()

	// Run JS script and send the request
	result := page.MustEval(string(*p.RequestScript))

	// Convert result from the JS script to JSON
	resultJSON, err := result.MarshalJSON()
	if err != nil {
		log.Fatal().Msgf("Could not convert result to JSON:\n%v", err)
	}

	// Parse the JSON response to proper Golang struct
	var perplexityResponse Response
	if err := json.Unmarshal(resultJSON, &perplexityResponse); err != nil {
		log.Fatal().Msgf("Error parsing result:\n%v", err)
	}

	return perplexityResponse, nil
}

func (p *Perplexity) SetHeaders(headers *RequestHeaders) {
	p.RequestHeaders = headers
	p.RequestScript = p.RequestScript.Update(p.RequestHeaders, p.RequestBody)
}

func (p *Perplexity) SetBody(body *RequestBody) {
	p.RequestBody = body
	p.RequestScript = p.RequestScript.Update(p.RequestHeaders, p.RequestBody)
}
