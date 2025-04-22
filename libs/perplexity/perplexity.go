package perplexity

import (
	"fmt"

	"github.com/ezydark/ezHead/libs/perplexity/request"
	"github.com/go-rod/rod"
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
	script, err := new_perplex.ReqScript.Update(new_perplex.ReqHeaders, new_perplex.ReqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to update request script:\n%w", err)
	}
	new_perplex.ReqScript = script

	return new_perplex, nil
}

func (p *PerplexityReq) SetChatSession(ChatUUID string) (*PerplexityReq, error) {
	p.ReqBody.ToFollowup(ChatUUID)
	script, err := p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to update request script:\n%w", err)
	}
	p.ReqScript = script

	return p, nil
}

func (p *PerplexityReq) SendRequest(page *rod.Page, query string) error {
	p.ReqBody.QueryStr = query
	script, err := p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
	if err != nil {
		return fmt.Errorf("failed to update request script:\n%w", err)
	}
	p.ReqScript = script

	// Wait for the page to fully load
	page.MustWaitStable()

	// Run JS script and send the request
	_ = page.MustEval(string(*p.ReqScript))

	// Convert result from the JS script to JSON
	// _, err = result.MarshalJSON()
	// if err != nil {
	// 	log.Fatal().Msgf("Could not convert result to JSON:\n%v", err)
	// }

	// Parse the JSON response to proper Golang struct
	// var perplexityResponse Response
	// if err := json.Unmarshal(resultJSON, &perplexityResponse); err != nil {
	// 	log.Fatal().Msgf("Error parsing result:\n%v", err)
	// }

	// // Set specific chat session identifier to allow followup messages
	// if !p.ReqBody.IsFollowup() {
	// 	p.ReqBody.ToFollowup(perplexityResponse.FinalMessage.BackendUUID)
	// 	p.ReqScript = p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
	// }

	return nil
}

func (p *PerplexityReq) SetHeaders(headers *request.Headers) error {
	p.ReqHeaders = headers

	script, err := p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
	if err != nil {
		return fmt.Errorf("failed to update request script:\n%w", err)
	}
	p.ReqScript = script

	return nil
}

func (p *PerplexityReq) SetBody(body *request.Body) error {
	p.ReqBody = body

	script, err := p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
	if err != nil {
		return fmt.Errorf("failed to update request script:\n%w", err)
	}
	p.ReqScript = script

	return nil
}
