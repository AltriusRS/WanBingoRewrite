/**
 * Test script for local LM Studio moderation
 * Run with: node test-local-moderation.js [lm-studio-url]
 */

const testMessages = [
  // Should pass (non-toxic)
  { content: "Hello everyone!", expected: false },
  { content: "This product is terrible, I hate it", expected: false },
  { content: "LOL that's hilarious", expected: false },
  { content: "Check out this link: https://example.com", expected: false },

  // Should fail (toxic)
  { content: "You are all worthless pieces of shit", expected: true },
  { content: "Kill yourself, nobody cares", expected: true },
  { content: "I hate all niggers and jews", expected: true },
];

async function testLocalModeration(endpoint) {
  console.log(`Testing local LM Studio endpoint: ${endpoint}\n`);

  for (const test of testMessages) {
    try {
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ content: test.content }),
      });

      if (!response.ok) {
        console.error(`❌ HTTP ${response.status} for: "${test.content}"`);
        continue;
      }

      const result = await response.json();
      const passed = result.toxic === test.expected;

      console.log(`${passed ? '✅' : '❌'} "${test.content}"`);
      console.log(`   Expected: ${test.expected}, Got: ${result.toxic} (${(result.confidence * 100).toFixed(1)}% confidence)`);
      if (result.reason) {
        console.log(`   Reason: ${result.reason}`);
      }
      console.log('');

    } catch (error) {
      console.error(`❌ Error testing "${test.content}": ${error.message}\n`);
    }
  }
}

// Test LM Studio API directly
async function testLMStudioAPI(endpoint) {
  console.log(`Testing LM Studio API directly: ${endpoint}\n`);

  try {
    const response = await fetch(`${endpoint}/v1/models`, {
      method: 'GET',
    });

    if (response.ok) {
      const data = await response.json();
      console.log('✅ LM Studio API is accessible');
      console.log('Available models:', data.data?.map(m => m.id) || 'None listed');
    } else {
      console.error(`❌ LM Studio API returned ${response.status}`);
    }
  } catch (error) {
    console.error(`❌ Cannot connect to LM Studio: ${error.message}`);
    console.log('\nTroubleshooting:');
    console.log('1. Make sure LM Studio is running');
    console.log('2. Check that the Local Server is started');
    console.log('3. Verify the URL is correct (default: http://localhost:1234)');
    console.log('4. Check firewall settings');
  }

  console.log('');
}

// Usage
const lmStudioUrl = process.argv[2] || 'http://localhost:1234';
const workerUrl = process.argv[3]; // Optional worker URL for end-to-end testing

console.log('🚀 WAN Bingo Moderation Testing\n');

if (workerUrl) {
  console.log('Testing Cloudflare Worker...');
  await testLocalModeration(workerUrl);
} else {
  console.log('Testing LM Studio API directly...');
  await testLMStudioAPI(lmStudioUrl);
  console.log('Testing moderation logic...');
  await testLocalModeration(lmStudioUrl.replace('/v1', '') + '/moderation');
}