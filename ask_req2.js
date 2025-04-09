(async () => {
  let query = "Lol";
  let onlyFinal = true;

  const url = "https://www.perplexity.ai/rest/sse/perplexity_ask";
  const headers = {
    accept: "text/event-stream",
    "accept-language": "cs,en-US;q=0.9,en;q=0.8",
    "content-type": "application/json",
    priority: "u=1, i",
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "same-origin",
  };

  const body = {
    params: {
      last_backend_uuid: "9f0a26f5-48c4-4c1d-b97c-49141d4c1946",
      read_write_token: "a25454be-e94d-4a63-b04e-05116abfa1d1",
      attachments: [],
      language: "en-US",
      timezone: "Europe/Prague",
      search_focus: "internet",
      sources: ["web"],
      frontend_uuid: "0447638c-f6fe-4fe8-9cf6-d3b518634b45",
      mode: "copilot",
      model_preference: "claude37sonnetthinking",
      version: "2.18",
    },
    query_str: query,
  };

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(body),
      credentials: "include",
      mode: "cors",
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    const messages = onlyFinal ? null : [];
    let finalMessage = null;
    let buffer = "";

    try {
      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          break;
        }

        buffer += decoder.decode(value, { stream: true });

        const lines = buffer.split("\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const data = line.substring(6);

            if (data && data !== ":heartbeat") {
              try {
                if (data === "[DONE]") {
                  continue;
                }

                const parsedData = JSON.parse(data);

                // Only store all messages if onlyFinal is false
                if (!onlyFinal) {
                  messages.push(parsedData);
                }

                // Always track the final message
                if (parsedData.final_sse_message === true) {
                  finalMessage = parsedData;
                }
              } catch (parseError) {
                console.warn("Failed to parse SSE data:", parseError.message);
              }
            }
          }
        }
      }

      // Return structure depends on onlyFinal parameter
      if (onlyFinal) {
        return {
          success: true,
          finalMessage: finalMessage,
        };
      } else {
        return {
          success: true,
          messages: messages,
          finalMessage: finalMessage,
        };
      }
    } catch (streamError) {
      throw new Error(`Stream processing error: ${streamError.message}`);
    } finally {
      reader.releaseLock();
    }
  } catch (error) {
    console.error("Error fetching Perplexity data:", error);
    return {
      success: false,
      error: error.message,
    };
  }
})();
