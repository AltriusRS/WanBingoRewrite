/**
 * Test script for the LLM Moderation Worker
 * Run with: node test-moderation.js
 */

const testMessages = [
  // Should pass (non-toxic)
  { content: "Hello everyone!", expected: false },
  { content: "This product is terrible, I hate it", expected: false },
  { content: "LOL that's hilarious", expected: false },
  { content: "Check out this link: https://example.com", expected: false },
  { content: "The new GPU is amazing for gaming", expected: false },

  // Should fail (toxic)
  { content: "You are all worthless pieces of shit", expected: true },
  { content: "Kill yourself, nobody cares", expected: true },
  { content: "I hate all niggers and jews", expected: true },
  { content: "Go die in a fire, terrorist scum", expected: true },
  { content: "This is fucking bullshit and I want to murder someone", expected: true },
];

async function testModeration(endpoint) {
  console.log(`Testing moderation endpoint: ${endpoint}\n`);

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

// Usage
const endpoint = process.argv[2] || 'http://localhost:8787'; // Default to wrangler dev
testModeration(endpoint).catch(console.error);