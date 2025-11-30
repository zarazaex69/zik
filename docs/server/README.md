# Zik Server API Documentation

Complete API reference for the Zik AI Server - an OpenAI-compatible API gateway to GLM-4.6 via Z.AI.

## Base URLs

- **Local Development**: `http://localhost:8802`
- **Production**: `https://zcapi.zarazaex.xyz`

## Authentication

Currently, the API does not require authentication. All endpoints are publicly accessible.

---

## Endpoints

### 1. Health Check

Check server status and version.

**Endpoint**: `GET /health`

**Response**: `200 OK`

```json
{
  "status": "ok",
  "version": "0.1.0"
}
```

**Examples**:

```bash
# Local
curl http://localhost:8802/health

# Production
curl https://zcapi.zarazaex.xyz/health
```

---

### 2. List Models

Retrieve available AI models.

**Endpoint**: `GET /v1/models`

**Response**: `200 OK`

```json
{
  "object": "list",
  "data": [
    {
      "id": "GLM-4-6-API-V1",
      "object": "model",
      "created": 1732989234,
      "owned_by": "zik"
    }
  ]
}
```

**Examples**:

```bash
# Local
curl http://localhost:8802/v1/models

# Production
curl https://zcapi.zarazaex.xyz/v1/models
```

---

### 3. Chat Completions

Create a chat completion using OpenAI-compatible format.

**Endpoint**: `POST /v1/chat/completions`

**Headers**:
- `Content-Type: application/json`

**Request Body**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `model` | string | No | Model ID (default: `GLM-4-6-API-V1`) |
| `messages` | array | Yes | Array of message objects (min: 1) |
| `stream` | boolean | No | Enable streaming (default: `false`) |
| `temperature` | float | No | Sampling temperature (0-2) |
| `max_tokens` | integer | No | Maximum tokens to generate |
| `top_p` | float | No | Nucleus sampling (0-1) |
| `stream_options` | object | No | Streaming configuration |
| `tools` | array | No | Function calling tools |
| `thinking` | boolean | No | Enable reasoning mode |

**Message Object**:

```json
{
  "role": "system|user|assistant",
  "content": "Message text or structured content"
}
```

**Stream Options**:

```json
{
  "include_usage": true
}
```

---

#### 3.1. Non-Streaming Mode

Returns a complete response in a single JSON object.

**Request Example**:

```json
{
  "model": "GLM-4-6-API-V1",
  "messages": [
    {
      "role": "system",
      "content": "You are a helpful assistant."
    },
    {
      "role": "user",
      "content": "What is the capital of France?"
    }
  ],
  "stream": false,
  "temperature": 0.7,
  "max_tokens": 150
}
```

**Response**: `200 OK`

```json
{
  "id": "chatcmpl-abc123",
  "object": "chat.completion",
  "created": 1732989234,
  "model": "GLM-4-6-API-V1",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "The capital of France is Paris.",
        "reasoning_content": ""
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 25,
    "completion_tokens": 8,
    "total_tokens": 33
  }
}
```

**cURL Examples**:

```bash
# Local - Simple request
curl -X POST http://localhost:8802/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Hello! How are you?"}
    ],
    "stream": false
  }'

# Production - Simple request
curl -X POST https://zcapi.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Hello! How are you?"}
    ],
    "stream": false
  }'

# Local - With system prompt and parameters
curl -X POST http://localhost:8802/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "system", "content": "You are a helpful coding assistant."},
      {"role": "user", "content": "Write a hello world in Python"}
    ],
    "stream": false,
    "temperature": 0.7,
    "max_tokens": 200
  }'

# Production - With system prompt and parameters
curl -X POST https://zcapi.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "system", "content": "You are a helpful coding assistant."},
      {"role": "user", "content": "Write a hello world in Python"}
    ],
    "stream": false,
    "temperature": 0.7,
    "max_tokens": 200
  }'

# Local - Multi-turn conversation
curl -X POST http://localhost:8802/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "What is 2+2?"},
      {"role": "assistant", "content": "2+2 equals 4."},
      {"role": "user", "content": "What about 2+3?"}
    ],
    "stream": false
  }'

# Production - Multi-turn conversation
curl -X POST https://zcapi.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "What is 2+2?"},
      {"role": "assistant", "content": "2+2 equals 4."},
      {"role": "user", "content": "What about 2+3?"}
    ],
    "stream": false
  }'
```

---

#### 3.2. Streaming Mode

Returns Server-Sent Events (SSE) stream for real-time responses.

**Request Example**:

```json
{
  "model": "GLM-4-6-API-V1",
  "messages": [
    {
      "role": "user",
      "content": "Count from 1 to 5"
    }
  ],
  "stream": true,
  "stream_options": {
    "include_usage": true
  }
}
```

**Response**: `200 OK` with `Content-Type: text/event-stream`

Stream format:

```
data: {"id":"chatcmpl-xyz","object":"chat.completion.chunk","created":1732989234,"model":"GLM-4-6-API-V1","choices":[{"index":0,"delta":{"role":"assistant","content":"1"}}]}

data: {"id":"chatcmpl-xyz","object":"chat.completion.chunk","created":1732989234,"model":"GLM-4-6-API-V1","choices":[{"index":0,"delta":{"content":", 2"}}]}

data: {"id":"chatcmpl-xyz","object":"chat.completion.chunk","created":1732989234,"model":"GLM-4-6-API-V1","choices":[{"index":0,"delta":{"content":", 3, 4, 5"}}]}

data: {"id":"chatcmpl-xyz","object":"chat.completion.chunk","created":1732989234,"model":"GLM-4-6-API-V1","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":"stop"}]}

data: {"id":"chatcmpl-xyz","object":"chat.completion.chunk","created":1732989234,"model":"GLM-4-6-API-V1","choices":[],"usage":{"prompt_tokens":12,"completion_tokens":15,"total_tokens":27}}

data: [DONE]
```

**cURL Examples**:

```bash
# Local - Basic streaming
curl -X POST http://localhost:8802/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Tell me a short story"}
    ],
    "stream": true
  }'

# Production - Basic streaming
curl -X POST https://zcapi.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Tell me a short story"}
    ],
    "stream": true
  }'

# Local - Streaming with usage tracking
curl -X POST http://localhost:8802/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Explain quantum computing"}
    ],
    "stream": true,
    "stream_options": {
      "include_usage": true
    }
  }'

# Production - Streaming with usage tracking
curl -X POST https://zcapi.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Explain quantum computing"}
    ],
    "stream": true,
    "stream_options": {
      "include_usage": true
    }
  }'

# Local - Streaming with temperature control
curl -X POST http://localhost:8802/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Write a creative poem"}
    ],
    "stream": true,
    "temperature": 1.2
  }'

# Production - Streaming with temperature control
curl -X POST https://zcapi.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Write a creative poem"}
    ],
    "stream": true,
    "temperature": 1.2
  }'
```

---

#### 3.3. Reasoning Mode

Enable extended reasoning with the `thinking` parameter.

**Request Example**:

```json
{
  "model": "GLM-4-6-API-V1",
  "messages": [
    {
      "role": "user",
      "content": "Solve this logic puzzle: If all roses are flowers and some flowers fade quickly, can we conclude that some roses fade quickly?"
    }
  ],
  "thinking": true,
  "stream": false
}
```

**Response**: `200 OK`

```json
{
  "id": "chatcmpl-reasoning",
  "object": "chat.completion",
  "created": 1732989234,
  "model": "GLM-4-6-API-V1",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "No, we cannot conclude that some roses fade quickly...",
        "reasoning_content": "Let me analyze this step by step: 1) All roses are flowers (universal statement)..."
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 45,
    "completion_tokens": 120,
    "total_tokens": 165
  }
}
```

**cURL Examples**:

```bash
# Local - Reasoning mode
curl -X POST http://localhost:8802/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "What is 15% of 240?"}
    ],
    "thinking": true,
    "stream": false
  }'

# Production - Reasoning mode
curl -X POST https://zcapi.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "What is 15% of 240?"}
    ],
    "thinking": true,
    "stream": false
  }'

# Local - Reasoning with streaming
curl -X POST http://localhost:8802/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Prove that the square root of 2 is irrational"}
    ],
    "thinking": true,
    "stream": true
  }'

# Production - Reasoning with streaming
curl -X POST https://zcapi.zarazaex.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "GLM-4-6-API-V1",
    "messages": [
      {"role": "user", "content": "Prove that the square root of 2 is irrational"}
    ],
    "thinking": true,
    "stream": true
  }'
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "message": "Error description",
    "type": "invalid_request_error",
    "code": 400
  }
}
```

**Common Error Codes**:

| Code | Description |
|------|-------------|
| `400` | Bad Request - Invalid JSON or validation error |
| `500` | Internal Server Error - Upstream service failure |

**Examples**:

```json
// Missing required field
{
  "error": {
    "message": "Field validation for 'Messages' failed on the 'required' tag",
    "type": "invalid_request_error",
    "code": 400
  }
}

// Invalid temperature
{
  "error": {
    "message": "Field validation for 'Temperature' failed on the 'gte' tag",
    "type": "invalid_request_error",
    "code": 400
  }
}
```

---

## Response Fields Reference

### Chat Completion Object

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique completion identifier |
| `object` | string | Object type (`chat.completion` or `chat.completion.chunk`) |
| `created` | integer | Unix timestamp |
| `model` | string | Model used for completion |
| `choices` | array | Array of completion choices |
| `usage` | object | Token usage statistics (non-streaming or with `include_usage`) |

### Choice Object

| Field | Type | Description |
|-------|------|-------------|
| `index` | integer | Choice index (always 0) |
| `message` | object | Complete message (non-streaming) |
| `delta` | object | Message delta (streaming) |
| `finish_reason` | string | Reason for completion (`stop`, `length`, `null`) |

### Message Object

| Field | Type | Description |
|-------|------|-------------|
| `role` | string | Message role (`assistant`) |
| `content` | string | Main response content |
| `reasoning_content` | string | Extended reasoning (if `thinking: true`) |
| `tool_calls` | array | Function calls (if tools provided) |

### Usage Object

| Field | Type | Description |
|-------|------|-------------|
| `prompt_tokens` | integer | Tokens in the prompt |
| `completion_tokens` | integer | Tokens in the completion |
| `total_tokens` | integer | Total tokens used |

---

## Integration Examples

### Python (OpenAI SDK)

```python
from openai import OpenAI

# Local
client = OpenAI(
    base_url="http://localhost:8802/v1",
    api_key="not-needed"
)

# Production
client = OpenAI(
    base_url="https://zcapi.zarazaex.xyz/v1",
    api_key="not-needed"
)

# Non-streaming
response = client.chat.completions.create(
    model="GLM-4-6-API-V1",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)
print(response.choices[0].message.content)

# Streaming
stream = client.chat.completions.create(
    model="GLM-4-6-API-V1",
    messages=[
        {"role": "user", "content": "Count to 5"}
    ],
    stream=True
)
for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

### JavaScript (Node.js)

```javascript
import OpenAI from 'openai';

// Local
const client = new OpenAI({
  baseURL: 'http://localhost:8802/v1',
  apiKey: 'not-needed'
});

// Production
const client = new OpenAI({
  baseURL: 'https://zcapi.zarazaex.xyz/v1',
  apiKey: 'not-needed'
});

// Non-streaming
const response = await client.chat.completions.create({
  model: 'GLM-4-6-API-V1',
  messages: [{ role: 'user', content: 'Hello!' }]
});
console.log(response.choices[0].message.content);

// Streaming
const stream = await client.chat.completions.create({
  model: 'GLM-4-6-API-V1',
  messages: [{ role: 'user', content: 'Count to 5' }],
  stream: true
});

for await (const chunk of stream) {
  process.stdout.write(chunk.choices[0]?.delta?.content || '');
}
```

### Go

```go
package main

import (
    "context"
    "fmt"
    "github.com/sashabaranov/go-openai"
)

func main() {
    config := openai.DefaultConfig("not-needed")
    
    // Local
    config.BaseURL = "http://localhost:8802/v1"
    
    // Production
    // config.BaseURL = "https://zcapi.zarazaex.xyz/v1"
    
    client := openai.NewClientWithConfig(config)
    
    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: "GLM-4-6-API-V1",
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: "Hello!",
                },
            },
        },
    )
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

---

## Rate Limits

Currently, there are no enforced rate limits. However, please use the API responsibly.

---

## Support

For issues, questions, or contributions:

- **GitHub**: [zarazaex/zik](https://github.com/zarazaex/zik)
- **Issues**: [GitHub Issues](https://github.com/zarazaex/zik/issues)

---

## Changelog

### v0.1.0 (Current)

- OpenAI-compatible `/v1/chat/completions` endpoint
- Streaming and non-streaming modes
- Token usage tracking
- Reasoning mode support
- Health check and models endpoints
