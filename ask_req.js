(async () => {
  try {
    const response = await fetch(
      "https://www.perplexity.ai/rest/sse/perplexity_ask",
      {
        headers: {
          accept: "text/event-stream",
          "accept-language": "cs,en-US;q=0.9,en;q=0.8",
          "content-type": "application/json",
          priority: "u=1, i",
          "sec-ch-ua":
            '"Microsoft Edge";v="135", "Not-A.Brand";v="8", "Chromium";v="135"',
          "sec-ch-ua-arch": '"x86"',
          "sec-ch-ua-bitness": '"64"',
          "sec-ch-ua-full-version": '"135.0.3179.54"',
          "sec-ch-ua-full-version-list":
            '"Microsoft Edge";v="135.0.3179.54", "Not-A.Brand";v="8.0.0.0", "Chromium";v="135.0.7049.42"',
          "sec-ch-ua-mobile": "?0",
          "sec-ch-ua-model": '""',
          "sec-ch-ua-platform": '"Windows"',
          "sec-ch-ua-platform-version": '"19.0.0"',
          "sec-fetch-dest": "empty",
          "sec-fetch-mode": "cors",
          "sec-fetch-site": "same-origin",
        },
        referrer:
          "https://www.perplexity.ai/search/heyyyyy-5hdxDP5STGW52sl_9gJhQQ",
        referrerPolicy: "strict-origin-when-cross-origin",
        body: '{"params":{"last_backend_uuid":"9f0a26f5-48c4-4c1d-b97c-49141d4c1946","read_write_token":"a25454be-e94d-4a63-b04e-05116abfa1d1","attachments":[],"language":"en-US","timezone":"Europe/Prague","search_focus":"internet","sources":["web"],"frontend_uuid":"0447638c-f6fe-4fe8-9cf6-d3b518634b45","mode":"copilot","model_preference":"claude37sonnetthinking","is_related_query":false,"is_sponsored":false,"visitor_id":"7ddb064a-37a6-4058-92f4-f94000fb00cf","user_nextauth_id":"76732ce2-a124-40cb-a0ce-db657d4344b9","prompt_source":"user","query_source":"followup","local_search_enabled":true,"browser_history_summary":[],"is_incognito":false,"use_schematized_api":true,"send_back_text_in_streaming_api":false,"supported_block_use_cases":["answer_modes","media_items","knowledge_cards","inline_entity_cards","place_widgets","finance_widgets","sports_widgets","shopping_widgets","jobs_widgets","search_result_widgets","entity_list_answer","todo_list"],"client_coordinates":null,"is_nav_suggestions_disabled":false,"followup_source":"link","version":"2.18"},"query_str":"Lol"}',
        method: "POST",
        mode: "cors",
        credentials: "include",
      },
    );

    // For SSE (Server-Sent Events), we need to handle the stream
    const reader = response.body.getReader();
    let result = "";

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      result += new TextDecoder().decode(value);
    }

    return result;
  } catch (error) {
    return "Error: " + error.message;
  }
})();
