async () => {
  window.goLogInfo("Sending request to Perplexity API");

  const perplexityApiConfig = {
    url: "https://www.perplexity.ai/rest/sse/perplexity_ask",
    headers: %s,
    body: %s,
    method: "POST",
    mode: "cors",
    credentials: "include",
  };

  function handleError(error, context) {
    window.goLogFatal(`Error encountered during ${context}:\n${error.message}`);
    // TODO: Implement more robust error handling (e.g., UI feedback, logging service)
  }

  async function makePerplexityRequest(config) {
    try {
      const response = await fetch(config.url, {
        method: config.method,
        headers: config.headers,
        body: JSON.stringify(config.body),
        mode: config.mode,
        credentials: config.credentials,
        referrer: config.referrer,
        referrerPolicy: config.referrerPolicy,
      });

      // Check if the HTTP status code indicates success (e.g., 200-299)[3]
      if (!response.ok) {
        // Throw an error for bad responses (like 4xx or 5xx)[3]
        throw new Error(
          `HTTP error! Status: ${response.status} ${response.statusText}`,
        );
      }

      window.goLogInfo("Request successful, received response headers.");
      // Return the Response object which contains the streamable body[3][6]
      return response;
    } catch (error) {
      // Catch network errors or errors thrown from the response check[3]
      handleError(error, "API request");
      return null; // Indicate failure
    }
  }

  async function* readStreamChunks(stream) {
    const reader = stream.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    try {
      while (true) {
        const { done, value } = await reader.read();
        if (done) {
          window.goLogInfo("Stream finished.");
          break;
        }

        // Add new text to buffer
        buffer += decoder.decode(value, { stream: true });

        // Process complete SSE messages
        let eventEndIndex;
        while ((eventEndIndex = buffer.indexOf('\r\n\r\n')) !== -1) {
          const event = buffer.substring(0, eventEndIndex);
          buffer = buffer.substring(eventEndIndex + 4); // Skip '\r\n\r\n'

          // Look for the data part
          const dataMatch = event.match(/data: ({.*})/);
          if (dataMatch && dataMatch[1]) {
            try {
              // Parse JSON and yield it
              const jsonData = JSON.parse(dataMatch[1]);
              yield jsonData;
            } catch (e) {
              window.goLogError(`Error parsing JSON: ${e.message}`);
            }
          }
        }
      }
    } catch (error) {
      handleError(error, "reading stream chunk");
    } finally {
      reader.releaseLock();
      window.goLogInfo("Stream reader released.");
    }
  }

  async function executeRequestAndStreamResponse() {
    window.goLogInfo("Starting Perplexity API request...");
    const response = await makePerplexityRequest(perplexityApiConfig);

    // Proceed only if the initial request was successful and we have a response with a body
    if (response && response.body) {
      window.goLogInfo("Request successful, starting to process stream...");
      try {
        // Iterate over the async generator to get chunks as they arrive[4][5][6]
        for await (const chunk of readStreamChunks(response.body)) {
          // window.goLogInfo(`Received chunk:\n${chunk}`);
          window.goProcessStreamChunk(chunk);
        }
        window.goLogInfo("Finished processing all stream chunks.");
      } catch (error) {
        // Catch errors specifically related to iterating or processing the stream chunks
        handleError(error, "processing stream");
      }
    } else {
      window.goLogInfo(
        "Failed to get a valid response or response body to stream.",
      );
    }
  }

  await executeRequestAndStreamResponse();
};
