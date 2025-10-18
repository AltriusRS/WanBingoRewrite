/**
 * WAN Bingo Chat Moderation Worker (Local LM Studio)
 * Cloudflare Worker for local LLM-based content moderation using LM Studio
 */

export default {
  async fetch(request, env, ctx) {
    // Handle CORS preflight requests
    if (request.method === 'OPTIONS') {
      return new Response(null, {
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'POST, OPTIONS',
          'Access-Control-Allow-Headers': 'Content-Type',
        },
      });
    }

    // Only allow POST requests
    if (request.method !== 'POST') {
      return new Response(JSON.stringify({ error: 'Method not allowed' }), {
        status: 405,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
        },
      });
    }

    try {
      const { content } = await request.json();

      if (!content || typeof content !== 'string') {
        return new Response(JSON.stringify({ error: 'Invalid content' }), {
          status: 400,
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
          },
        });
      }

      // Moderate the content using local LM Studio
      const result = await moderateContentLocal(content, env);

      return new Response(JSON.stringify(result), {
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
        },
      });

    } catch (error) {
      console.error('Moderation error:', error);
      return new Response(JSON.stringify({
        toxic: false,
        confidence: 0,
        reason: 'Error processing request'
      }), {
        status: 500,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
        },
      });
    }
  },
};

/**
 * Moderate content using local LM Studio API
 */
async function moderateContentLocal(content, env) {
  const LM_STUDIO_URL = env.LM_STUDIO_URL || 'http://localhost:1234';
  const MODEL = env.MODEL || 'local-model'; // LM Studio model name

  // Truncate content if too long
  const truncatedContent = content.length > 2000 ? content.substring(0, 2000) + '...' : content;

  const prompt = `You are a content moderation AI for a live streaming chat platform. Your task is to analyze the following message and determine if it contains toxic, harmful, or inappropriate content.

Guidelines:
- Be conservative: only flag content that is clearly toxic, hateful, or harmful
- Consider context: jokes, sarcasm, or cultural references may not be toxic
- Allow: mild profanity in casual conversation, criticism of products/companies, technical discussions
- Flag: hate speech, threats, harassment, explicit content, spam, doxxing, self-harm promotion
- Confidence: rate 0.0-1.0 how certain you are this content violates community guidelines

Message to analyze: "${truncatedContent}"

Respond with ONLY a JSON object in this exact format:
{"toxic": boolean, "confidence": number, "reason": "brief explanation if toxic"}`;

  try {
    const response = await fetch(`${LM_STUDIO_URL}/v1/chat/completions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        model: MODEL,
        messages: [
          {
            role: 'system',
            content: 'You are a content moderation expert. Always respond with valid JSON only.'
          },
          {
            role: 'user',
            content: prompt
          }
        ],
        max_tokens: 150,
        temperature: 0.1, // Low temperature for consistent results
        stream: false,
      }),
    });

    if (!response.ok) {
      throw new Error(`LM Studio API error: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();
    const aiResponse = data.choices[0]?.message?.content;

    if (!aiResponse) {
      throw new Error('No response from LM Studio');
    }

    // Parse the JSON response
    const result = JSON.parse(aiResponse.trim());

    // Validate response format
    if (typeof result.toxic !== 'boolean' || typeof result.confidence !== 'number') {
      throw new Error('Invalid response format from AI');
    }

    return {
      toxic: result.toxic,
      confidence: Math.max(0, Math.min(1, result.confidence)), // Clamp to 0-1
      reason: result.reason || (result.toxic ? 'Detected by AI moderation' : ''),
    };

  } catch (error) {
    console.error('Local AI moderation failed:', error);
    // Return safe defaults on error
    return {
      toxic: false,
      confidence: 0,
      reason: 'Local moderation service temporarily unavailable',
    };
  }
}