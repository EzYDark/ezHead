package perplexity

import (
	"fmt"

	"github.com/ezydark/ezHead/libs"
	"github.com/ezydark/ezHead/libs/perplexity/request"
	"github.com/ezydark/ezHead/libs/perplexity/response"
	"github.com/go-rod/rod"
	"github.com/rs/zerolog/log"
	"github.com/ysmood/gson"
)

type PerplexityReq struct {
	ReqHeaders *request.Headers
	ReqBody    *request.Body
	ReqScript  *request.Script
	RodPage    *rod.Page
}

// Initialize new chat session on Perplexity
func Init() (*PerplexityReq, error) {
	new_perplex := &PerplexityReq{
		ReqHeaders: new(request.Headers),
		ReqBody:    new(request.Body),
		ReqScript:  new(request.Script),
		RodPage:    nil,
	}

	new_perplex.ReqHeaders = new_perplex.ReqHeaders.Default()

	new_perplex.ReqBody = new_perplex.ReqBody.Default()
	script, err := new_perplex.ReqScript.Update(new_perplex.ReqHeaders, new_perplex.ReqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to update request script:\n%w", err)
	}

	new_perplex.ReqScript = script

	rod_browser := libs.GetRodBrowser()
	new_perplex.RodPage = rod_browser.MustPage("https://www.perplexity.ai/")

	libs.SetPageSettings(new_perplex.RodPage)
	libs.ExposeGoLogger(new_perplex.RodPage)
	ExposeProcStreamChunk(new_perplex.RodPage)

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
	result := page.MustEval(string(*p.ReqScript))

	// Convert result from the JS script to JSON
	_, err = result.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to convert result to JSON:\n%w", err)
	}

	// // Parse the JSON response to proper Golang struct
	// var perplexityResponse types.Response
	// if err := json.Unmarshal(resultJSON, &perplexityResponse); err != nil {
	// 	return fmt.Errorf("failed to parse result:\n%w", err)
	// }

	// // Set specific chat session identifier to allow followup messages
	// if !p.ReqBody.IsFollowup() {
	// 	p.ReqBody.ToFollowup(perplexityResponse.BackendUUID)
	// 	script, err := p.ReqScript.Update(p.ReqHeaders, p.ReqBody)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to update request script:\n%w", err)
	// 	}
	// 	p.ReqScript = script
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

func ExposeProcStreamChunk(page *rod.Page) *rod.Page {
	_ = page.MustExpose("goProcessStreamChunk", func(chunk gson.JSON) (any, error) {
		_, err := response.ProcessStreamChunk(chunk)
		if err != nil {
			log.Fatal().Msgf("Error processing stream chunk:\n%v", err)
		}

		return nil, nil
	})

	return page
}

// GetResponseChunks returns the current chunks from the response package
func GetResponseChunks() []string {
	return response.ChunkStorage
}

// ClearResponseChunks clears the chunks in the response package
func ClearResponseChunks() {
	response.ChunkStorage = nil
}
