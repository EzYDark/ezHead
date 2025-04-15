package perplexity

import (
	"encoding/json"

	"github.com/ezydark/ezHead/libs/perplexity/request"
	"github.com/go-rod/rod"
	"github.com/rs/zerolog/log"
)

type PerplexityReq struct {
	ReqHeaders *request.Headers
	ReqBody    *request.Body
	ReqScript  *request.Script
}

// Initialize new chat session on Perplexity
func Init() (*PerplexityReq, error) {
	new_perplex := &PerplexityReq{
		ReqHeaders: new(request.Headers),
		ReqBody:    new(request.Body),
		ReqScript:  new(request.Script),
	}

	new_perplex.ReqHeaders = new_perplex.ReqHeaders.Default()
	new_perplex.ReqBody = new_perplex.ReqBody.Default()
	new_perplex.ReqScript = new_perplex.ReqScript.Update(new_perplex.ReqHeaders, new_perplex.ReqBody)

	return new_perplex, nil
}

func (p *PerplexityReq) SetChatSession(ChatUUID string) *PerplexityReq {
	p.ReqBody.ToFollowup(ChatUUID)
	p.ReqScript = p.ReqScript.Update(p.ReqHeaders, p.ReqBody)

	return p
}

func (p *PerplexityReq) SendRequest(page *rod.Page, query string) (Response, error) {
	p.ReqBody.QueryStr = query
	p.ReqScript = p.ReqScript.Update(p.ReqHeaders, p.ReqBody)

	// Wait for the page to fully load
	page.MustWaitStable()

	// Run JS script and send the request
	result := page.MustEval(string(*p.ReqScript))

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

	// Set specific chat session identifier to allow followup messages
	if !p.ReqBody.IsFollowup() {
		p.ReqBody.ToFollowup(perplexityResponse.FinalMessage.BackendUUID)
		p.ReqScript = p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
	}

	return perplexityResponse, nil
}

func (p *PerplexityReq) SetHeaders(headers *request.Headers) {
	p.ReqHeaders = headers
	p.ReqScript = p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
}

func (p *PerplexityReq) SetBody(body *request.Body) {
	p.ReqBody = body
	p.ReqScript = p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
}
