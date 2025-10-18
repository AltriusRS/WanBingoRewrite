# WAN Bingo LLM Moderation Worker

This document provides comprehensive guidance for implementing and deploying a Cloudflare Worker for LLM-based content moderation in the WAN Bingo chat system.

## Overview

The WAN Bingo chat system includes comprehensive content moderation that works out-of-the-box:

- **Keyword Filtering**: Blocks common slurs, hate speech, and inappropriate language
- **Markdown Filtering**: Prevents abuse of formatting (tables, headers, images, etc.)
- **Pattern Detection**: Catches excessive caps, repeated characters, and bypass attempts

For enhanced AI-powered moderation, you can optionally deploy the LLM Moderation Worker in one of two modes:

1. **Cloud API Mode**: Uses OpenAI's GPT models via API calls
2. **Local Mode**: Uses LM Studio with local models on your hardware

The base system provides reliable moderation without any external dependencies. The LLM worker adds sophisticated AI analysis for detecting subtle forms of toxicity that keyword filters might miss.

## Moderation Levels

### Base Moderation (Always Active)
The chat system includes robust moderation that works without any external setup:

- **Keyword filtering** for slurs and hate speech
- **Markdown filtering** to prevent formatting abuse
- **Pattern detection** for spam and bypass attempts
- **Zero configuration** required

### Enhanced Moderation (Optional)
Add AI-powered analysis for complex cases:

#### Cloud API Mode (OpenAI)
- Pay per request (~$0.15/1K messages)
- 99.9% uptime, zero maintenance
- Data sent to OpenAI (privacy consideration)

#### Local Mode (LM Studio)
- Free after setup (electricity only)
- Data stays local, maximum privacy
- Requires capable hardware (GPU recommended)

### When to Add Enhanced Moderation

**Consider enhanced moderation when:**
- You want to catch subtle forms of toxicity
- Community standards require high accuracy
- You have budget for AI analysis
- Privacy allows external API usage

**The base system is sufficient for:**
- Most community chat platforms
- Cost-conscious deployments
- Environments where privacy is critical
- Initial testing and development

## Architecture

```
User Message → Go Server → Cloudflare Worker → OpenAI API → Moderation Result
                      ↓
                Database + SSE Broadcast
```

## Prerequisites

- Cloudflare account with Workers enabled
- OpenAI API key
- Node.js/npm for development (optional)

## Model Recommendations

### Primary Models

#### GPT-4o Mini (Recommended)
- **Model ID**: `gpt-4o-mini`
- **Cost**: ~$0.15 per 1M input tokens, ~$0.60 per 1M output tokens
- **Performance**: Excellent balance of accuracy and cost
- **Use Case**: Production deployment for cost-effective moderation

#### GPT-4o
- **Model ID**: `gpt-4o`
- **Cost**: ~$2.50 per 1M input tokens, ~$10 per 1M output tokens
- **Performance**: Highest accuracy for complex content analysis
- **Use Case**: High-stakes moderation or when maximum accuracy is required

#### GPT-3.5 Turbo
- **Model ID**: `gpt-3.5-turbo`
- **Cost**: ~$0.50 per 1M input tokens, ~$1.50 per 1M output tokens
- **Performance**: Good accuracy with lower cost than GPT-4
- **Use Case**: Development/testing or budget-conscious production

### Model Selection Guidelines

| Scenario | Recommended Model | Reasoning |
|----------|------------------|-----------|
| Production (cost-optimized) | `gpt-4o-mini` | Best balance of accuracy and cost |
| Production (max accuracy) | `gpt-4o` | Highest precision for complex content |
| Development/Testing | `gpt-3.5-turbo` | Lower cost for experimentation |
| High-volume chat | `gpt-4o-mini` | Cost-effective for frequent requests |

## Prompt Engineering

### Core Moderation Prompt

The worker uses a carefully crafted prompt designed to:

1. **Conservative Approach**: Only flag clearly toxic content
2. **Context Awareness**: Consider jokes, sarcasm, and cultural references
3. **Allow Common Content**: Mild profanity, product criticism, technical discussions
4. **Flag Harmful Content**: Hate speech, threats, harassment, explicit content

### Prompt Structure

```
You are a content moderation AI for a live streaming chat platform. Your task is to analyze the following message and determine if it contains toxic, harmful, or inappropriate content.

Guidelines:
- Be conservative: only flag content that is clearly toxic, hateful, or harmful
- Consider context: jokes, sarcasm, or cultural references may not be toxic
- Allow: mild profanity in casual conversation, criticism of products/companies, technical discussions
- Flag: hate speech, threats, harassment, explicit content, spam, doxxing, self-harm promotion
- Confidence: rate 0.0-1.0 how certain you are this content violates community guidelines

Message to analyze: "[CONTENT]"

Respond with ONLY a JSON object in this exact format:
{"toxic": boolean, "confidence": number, "reason": "brief explanation if toxic"}
```

### Prompt Optimization

#### Key Principles
- **Clear Instructions**: Explicit guidelines reduce false positives
- **Context Examples**: Specific examples of allowed/forbidden content
- **Confidence Scoring**: Numerical confidence helps with threshold tuning
- **Structured Output**: JSON format ensures consistent parsing

#### Customization Options

For different communities, you can adjust the prompt:

```javascript
// Stricter moderation
const strictPrompt = `...Flag: hate speech, threats, harassment, explicit content, spam, doxxing, self-harm promotion, political extremism...`;

// More permissive
const permissivePrompt = `...Allow: strong language, controversial opinions, edgy humor...`;
```

## Deployment

### 1. Create Cloudflare Worker

#### Using Wrangler CLI

```bash
# Install Wrangler
npm install -g wrangler

# Login to Cloudflare
wrangler auth login

# Create new worker
wrangler init moderation-worker
cd moderation-worker
```

#### Manual Setup

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. Navigate to Workers & Pages
3. Click "Create Worker"
4. Name it (e.g., `wan-bingo-moderation`)

### 2. Deploy Worker Code

Replace the default worker code with the content from `cloudflare-worker-moderation.js`.

### 3. Configure Environment Variables

#### Using Wrangler

```bash
# Set OpenAI API key
wrangler secret put OPENAI_API_KEY

# Set model (optional, defaults to gpt-4o-mini)
wrangler secret put MODEL
```

#### Using Dashboard

1. Go to Worker Settings → Variables
2. Add the following secrets:
   - `OPENAI_API_KEY`: Your OpenAI API key
   - `MODEL`: Model ID (optional, defaults to `gpt-4o-mini`)

### 4. Deploy

```bash
wrangler deploy
```

### 5. Get Worker URL

After deployment, note the worker URL (e.g., `https://wan-bingo-moderation.your-subdomain.workers.dev`)

## Local LM Studio Setup

For local deployment using LM Studio, follow these steps:

### 1. Install LM Studio

Download and install LM Studio from [https://lmstudio.ai/](https://lmstudio.ai/)

### 2. Download Models

#### Recommended Models for Moderation

**Primary Recommendation:**
- **Qwen2.5-7B-Instruct** (7B parameters)
  - Excellent balance of speed and accuracy
  - Good at following instructions
  - ~4GB VRAM required

**High Accuracy (slower):**
- **Qwen2.5-14B-Instruct** (14B parameters)
  - Best accuracy for complex content
  - ~8GB VRAM required

**Fast & Lightweight:**
- **Phi-3-mini-4k-instruct** (3.8B parameters)
  - Very fast inference
  - Good for high-volume chat
  - ~2GB VRAM required

#### Download Instructions

1. Open LM Studio
2. Go to "My Models" tab
3. Search for the model name
4. Click download
5. Wait for download to complete

### 3. Configure Local Server

1. Go to "Local Server" tab in LM Studio
2. Load your downloaded model
3. Set server configuration:
   - **Port**: 1234 (default)
   - **Context Length**: 2048-4096 (sufficient for moderation)
   - **GPU Layers**: Maximum available (for best performance)
   - **Threads**: Match your CPU cores

4. Start the local server

### 4. Test Local API

```bash
# Test the local API
curl -X POST http://localhost:1234/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "local-model",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

### 5. Deploy Worker for Local Mode

Use the `cloudflare-worker-moderation-local.js` file instead of the OpenAI version:

```bash
# Copy the local version
cp cloudflare-worker-moderation-local.js dist/worker.js

# Deploy
wrangler deploy
```

### 6. Configure Environment Variables

```bash
# Set local LM Studio URL
wrangler secret put LM_STUDIO_URL "http://your-local-ip:1234"

# Optional: Set model name (as configured in LM Studio)
wrangler secret put MODEL "qwen2.5-7b-instruct"
```

### 7. Network Access

Since LM Studio runs locally, you need to make it accessible to the Cloudflare Worker:

#### Option A: Expose Local Server (Development)
```bash
# Use ngrok or similar to expose local port
ngrok http 1234
# Use the ngrok URL as LM_STUDIO_URL
```

#### Option B: Run on VPS/Cloud Server
- Deploy LM Studio on a cloud VPS with GPU
- Configure security (firewall, authentication)
- Use the cloud server IP as LM_STUDIO_URL

#### Option C: Hybrid Approach
- Run LM Studio locally for development/testing
- Use OpenAI for production
- Switch environment variables as needed

### Hardware Requirements

#### Minimum (for basic functionality)
- **CPU**: 4+ cores
- **RAM**: 8GB
- **GPU**: Optional (CPU-only inference possible)

#### Recommended (for good performance)
- **CPU**: 8+ cores
- **RAM**: 16GB+
- **GPU**: NVIDIA GPU with 4GB+ VRAM (Tesla P4 is excellent)

#### Performance Expectations

| Hardware | Model Size | Inference Speed | Quality |
|----------|------------|----------------|---------|
| CPU-only | 7B | 5-10 seconds | Good |
| NVIDIA P4 | 7B | 1-2 seconds | Excellent |
| RTX 3080 | 14B | 0.5-1 second | Best |

## Configuration

### Environment Variables

#### Cloud API Mode (OpenAI)

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `OPENAI_API_KEY` | Yes | - | OpenAI API key for authentication |
| `MODEL` | No | `gpt-4o-mini` | OpenAI model to use for moderation |

#### Local Mode (LM Studio)

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LM_STUDIO_URL` | Yes | `http://localhost:1234` | URL of your LM Studio local server |
| `MODEL` | No | `local-model` | Model name as configured in LM Studio |

**Note**: Use either Cloud API variables OR Local mode variables, not both. The worker code automatically detects which mode to use based on available environment variables.

### Cost Optimization

#### Token Usage Estimation

- **Average message length**: ~50-100 characters
- **Tokens per request**: ~100-200 (including prompt)
- **Monthly cost estimate**: $5-20 for moderate chat activity

#### Rate Limiting

Consider implementing rate limiting to control costs:

```javascript
// Add to worker for rate limiting
const rateLimit = new Map();

export default {
  async fetch(request, env, ctx) {
    // Basic rate limiting by IP
    const clientIP = request.headers.get('CF-Connecting-IP');
    const now = Date.now();
    const windowMs = 60000; // 1 minute
    const maxRequests = 10; // 10 requests per minute per IP

    // Implementation details...
  }
}
```

## API Usage

### Request Format

```http
POST / HTTP/1.1
Content-Type: application/json

{
  "content": "This is a test message to moderate"
}
```

### Response Format

#### Successful Moderation

```json
{
  "toxic": false,
  "confidence": 0.05,
  "reason": ""
}
```

#### Toxic Content Detected

```json
{
  "toxic": true,
  "confidence": 0.95,
  "reason": "Contains hate speech and threats"
}
```

#### Error Response

```json
{
  "toxic": false,
  "confidence": 0,
  "reason": "Moderation service temporarily unavailable"
}
```

### Integration with Go Server

Set the environment variable in your Go server:

```bash
export LLM_MODERATION_ENDPOINT="https://your-worker.your-subdomain.workers.dev"
```

The Go code will automatically use this endpoint when configured.

## Testing

### Manual Testing

```bash
# Test with curl
curl -X POST https://your-worker.workers.dev \
  -H "Content-Type: application/json" \
  -d '{"content": "This is a test message"}'
```

### Test Cases

#### Should Pass (non-toxic)
- "Hello everyone!"
- "This product is terrible, I hate it"
- "LOL that's hilarious"
- "Check out this link: https://example.com"

#### Should Fail (toxic)
- "You are all worthless pieces of shit"
- "Kill yourself, nobody cares"
- "I hate all [protected group]"
- Threats or harassment

### Performance Testing

```bash
# Load testing with hey
hey -n 100 -c 10 -m POST \
  -H "Content-Type: application/json" \
  -d '{"content": "Test message"}' \
  https://your-worker.workers.dev
```

## Monitoring & Analytics

### Cloudflare Analytics

Monitor your worker through the Cloudflare dashboard:

1. **Requests/Second**: Track usage patterns
2. **Error Rates**: Monitor API failures
3. **Latency**: Ensure response times stay under 2-3 seconds
4. **Costs**: Track OpenAI API usage

### Custom Logging

Add logging to track moderation decisions:

```javascript
// Add to moderateContent function
console.log(`Moderation: ${content.substring(0, 50)}... -> toxic: ${result.toxic}, confidence: ${result.confidence}`);
```

### Alerting

Set up alerts for:
- High error rates (>5%)
- Increased latency (>3 seconds)
- Unusual request patterns

## Troubleshooting

### Common Issues

#### 401 Unauthorized from OpenAI
- Check `OPENAI_API_KEY` is correctly set
- Verify API key has sufficient credits
- Ensure key has appropriate permissions

#### 429 Rate Limited
- OpenAI API rate limits apply
- Implement client-side rate limiting
- Consider upgrading OpenAI plan for higher limits

#### Worker Timeouts
- OpenAI requests can take 2-5 seconds
- Increase worker timeout in dashboard (max 30 seconds)
- Consider caching for repeated content

#### CORS Issues
- Worker includes CORS headers for all origins
- For production, restrict to your domain:
```javascript
'Access-Control-Allow-Origin': 'https://yourdomain.com'
```

### Debug Mode

Enable debug logging:

```javascript
// Add to worker
console.log('Request received:', await request.text());
console.log('OpenAI response:', JSON.stringify(data, null, 2));
```

## Security Considerations

### API Key Protection
- Store OpenAI key as Cloudflare secret (never in code)
- Rotate keys regularly
- Monitor for unauthorized usage

### Input Validation
- Worker validates JSON input
- Content is truncated to prevent token abuse
- Rate limiting prevents DoS attacks

### Privacy
- Messages are processed in memory only
- No persistent storage of user content
- OpenAI may retain data per their privacy policy

## Cost Management

### Budget Planning

#### Base Moderation Costs
| Component | Monthly Cost | Notes |
|-----------|-------------|-------|
| Server Processing | $0 | Included in existing infrastructure |
| Maintenance | Minimal | No external dependencies |

#### Enhanced Moderation Costs

##### Cloud API Costs
| Model | Cost per 1K requests | Monthly (10K messages) |
|-------|---------------------|----------------------|
| GPT-4o Mini | ~$0.02 | ~$0.20 |
| GPT-3.5 Turbo | ~$0.06 | ~$0.60 |
| GPT-4o | ~$0.35 | ~$3.50 |

##### Local Mode Costs
| Component | Monthly Cost | Notes |
|-----------|-------------|-------|
| Electricity | $1-5 | Depends on usage and local rates |
| Hardware | $0 (existing) | Or VPS costs if using cloud GPU |
| Maintenance | Minimal | Software updates, occasional hardware |

**Recommendation**: Start with base moderation, then add enhanced moderation when needed for higher accuracy.

### Cost Optimization

1. **Model Selection**: Use GPT-4o Mini for production
2. **Caching**: Cache results for identical messages
3. **Batch Processing**: Group multiple messages (if applicable)
4. **Rate Limiting**: Prevent abuse that increases costs

### Monitoring Costs

```javascript
// Track costs (approximate)
const tokensUsed = data.usage?.total_tokens || 0;
const estimatedCost = (tokensUsed / 1000) * 0.002; // Adjust based on model
console.log(`Estimated cost: $${estimatedCost.toFixed(4)}`);
```

## Advanced Configuration

### Custom Models

For using fine-tuned models:

```javascript
const MODEL = env.CUSTOM_MODEL || 'ft:gpt-4o-mini:your-org:moderation-model';
```

### Multi-Model Fallback

```javascript
async function moderateContent(content, env) {
  // Try primary model
  try {
    return await moderateWithModel(content, env.MODEL || 'gpt-4o-mini', env);
  } catch (error) {
    // Fallback to cheaper model
    console.warn('Primary model failed, using fallback');
    return await moderateWithModel(content, 'gpt-3.5-turbo', env);
  }
}
```

### Regional Deployment

Deploy workers in multiple regions for lower latency:

```bash
# Deploy to different regions
wrangler deploy --regions ewr,lax
```

## Maintenance

### Regular Updates

1. **Monitor Performance**: Check latency and error rates weekly
2. **Update Models**: Test new OpenAI models for improved accuracy
3. **Review Logs**: Analyze false positives/negatives monthly
4. **Cost Analysis**: Review spending patterns quarterly

### Backup Plan

Always have a fallback when LLM service is unavailable:

```javascript
// In Go server, implement fallback logic
if (llmError) {
  // Use only keyword filtering
  return keywordModerationResult;
}
```

## Support

### Resources

- [Cloudflare Workers Documentation](https://developers.cloudflare.com/workers/)
- [OpenAI API Reference](https://platform.openai.com/docs/api-reference)
- [Wrangler CLI](https://developers.cloudflare.com/workers/wrangler/)

### Getting Help

1. Check Cloudflare worker logs in dashboard
2. Monitor OpenAI API usage and errors
3. Test with various message types
4. Review prompt effectiveness regularly

---

## Quick Start Checklist

### Base Moderation (No Setup Required)
The chat system works out-of-the-box with keyword and markdown filtering:
- [x] Keyword filtering active
- [x] Markdown filtering active
- [x] Pattern detection active
- [ ] Test with sample messages (optional)

### Enhanced Moderation (Optional)

#### Cloud API Mode
- [ ] Create Cloudflare account
- [ ] Set up OpenAI API key
- [ ] Deploy `cloudflare-worker-moderation.js`
- [ ] Configure `OPENAI_API_KEY` secret
- [ ] Set `LLM_MODERATION_ENDPOINT` in Go server
- [ ] Test enhanced moderation
- [ ] Monitor performance and costs

#### Local Mode
- [ ] Install LM Studio
- [ ] Download a moderation model (Qwen2.5-7B recommended)
- [ ] Configure and start local server
- [ ] Deploy `cloudflare-worker-moderation-local.js`
- [ ] Configure `LM_STUDIO_URL` secret
- [ ] Set `LLM_MODERATION_ENDPOINT` in Go server
- [ ] Test enhanced moderation
- [ ] Monitor local hardware performance

## Example Deployment

### Cloud API Mode

```bash
# 1. Initialize worker
wrangler init moderation-worker
cd moderation-worker

# 2. Copy worker code
cp /path/to/cloudflare-worker-moderation.js .

# 3. Set secrets
wrangler secret put OPENAI_API_KEY
wrangler secret put MODEL  # Optional

# 4. Deploy
wrangler deploy

# 5. Get URL from output
# Example: https://moderation-worker.your-subdomain.workers.dev

# 6. Configure Go server
export LLM_MODERATION_ENDPOINT="https://moderation-worker.your-subdomain.workers.dev"
```

### Local Mode

```bash
# 1. Install and setup LM Studio
# - Download from https://lmstudio.ai/
# - Download Qwen2.5-7B-Instruct model
# - Start local server on port 1234

# 2. Initialize worker
wrangler init moderation-worker-local
cd moderation-worker-local

# 3. Copy local worker code
cp /path/to/cloudflare-worker-moderation-local.js .

# 4. Set secrets
wrangler secret put LM_STUDIO_URL "http://your-local-ip:1234"

# 5. Deploy
wrangler deploy

# 6. Configure Go server
export LLM_MODERATION_ENDPOINT="https://moderation-worker-local.your-subdomain.workers.dev"
```</content>
</xai:function_call">LLM_MODERATION_WORKER.md