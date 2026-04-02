import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// Configuration: Ramping VUs
export const options = {
  stages: [
    { duration: '10s', target: 20 },  // Ramp-up: 0 to 20 users
    { duration: '30s', target: 20 },  // Steady state: 20 users
    { duration: '10s', target: 0 },   // Ramp-down
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],   // Error rate should be < 1%
    http_req_duration: ['p(95)<500'], // 95% of requests should be < 500ms
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Scenario 1: Fetch all accounts (GET) - Most frequent
  const getRes = http.get(`${BASE_URL}/accounts`, {
    tags: { name: 'GetAccounts' },
  });
  check(getRes, {
    'GET /accounts status is 200': (r) => r.status === 200,
  });

  sleep(1);

  // Scenario 2: Create a Transaction / Transfer (POST) - Less frequent
  // We simulate idempotency here
  if (Math.random() < 0.2) { // 20% chance to perform a transfer
    const payload = JSON.stringify({
      from_account_id: "00000000-0000-0000-0000-000000000001", // Assuming seeds or previous runs
      to_account_id: "00000000-0000-0000-0000-000000000002",
      amount: 1000,
    });

    const params = {
      headers: {
        'Content-Type': 'application/json',
        'Idempotency-Key': `k6-test-${randomString(10)}`, // Unique key per test iteration
      },
      tags: { name: 'Transfer' },
    };

    const postRes = http.post(`${BASE_URL}/transfer`, payload, params);
    
    check(postRes, {
      'POST /transfer status is 200 or 400': (r) => r.status === 200 || r.status === 400,
    });
  }

  sleep(1);
}
