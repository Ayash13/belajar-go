// Load Test — Mengukur performa API dibawah beban normal
// Menggunakan Ramping VUs: naik bertahap → sustain → turun bertahap
// Jalankan: k6 run scripts/load_test.js

import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  // Scenario: Ramping VUs
  stages: [
    { duration: '10s', target: 10 },  // Ramp up: 0 → 10 VU dalam 10 detik
    { duration: '30s', target: 10 },  // Sustain: 10 VU selama 30 detik
    { duration: '10s', target: 0 },   // Ramp down: 10 → 0 VU dalam 10 detik
  ],

  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    'http_req_duration{name:GetAccounts}': ['p(95)<300'],
    'http_req_duration{name:GetAccountByID}': ['p(95)<300'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Skenario 1: GET semua akun
  const allAccountsRes = http.get(`${BASE_URL}/accounts`, {
    tags: { name: 'GetAccounts' },
  });
  check(allAccountsRes, {
    'GET /accounts → 200': (r) => r.status === 200,
    'GET /accounts → has success status': (r) => JSON.parse(r.body).status === 'success',
  });

  sleep(0.5);

  // Skenario 2: GET akun by ID
  const accountID = 'acc-001';
  const singleAccountRes = http.get(`${BASE_URL}/accounts/${accountID}`, {
    tags: { name: 'GetAccountByID' },
  });
  check(singleAccountRes, {
    'GET /accounts/{id} → 200 or 404': (r) => r.status === 200 || r.status === 404,
  });

  sleep(0.5);

  // Skenario 3: POST transfer (10% dari iterasi saja)
  if (Math.random() < 0.1) {
    const transferPayload = JSON.stringify({
      from_account_id: 'acc-001',
      to_account_id: 'acc-002',
      amount: Math.floor(Math.random() * 50000) + 10000,
    });

    const transferRes = http.post(`${BASE_URL}/transfer`, transferPayload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: 'Transfer' },
    });
    check(transferRes, {
      'POST /transfer → 200': (r) => r.status === 200,
    });
  }

  sleep(1);
}
