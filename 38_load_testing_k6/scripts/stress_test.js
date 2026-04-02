// Stress Test — Mencari batas maksimum sistem sebelum degradasi/crash
// Menggunakan Ramping VUs yang agresif hingga 100 VU
// Jalankan: k6 run scripts/stress_test.js

import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  // Ramping VUs — naik agresif untuk menemukan breaking point
  stages: [
    { duration: '10s', target: 10 },   // Warm up
    { duration: '20s', target: 50 },   // Naikkan ke 50 VU
    { duration: '30s', target: 50 },   // Sustain 50 VU
    { duration: '20s', target: 100 },  // Naikkan ke 100 VU (stress zone)
    { duration: '30s', target: 100 },  // Sustain 100 VU
    { duration: '10s', target: 0 },    // Recovery — turun ke 0
  ],

  thresholds: {
    http_req_failed: ['rate<0.05'],        // toleransi error lebih tinggi (5%)
    http_req_duration: ['p(95)<1000'],      // 95% request < 1 detik
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Request campuran untuk mengsimulasikan traffic realistis
  const responses = http.batch([
    ['GET', `${BASE_URL}/accounts`, null, { tags: { name: 'GetAccounts' } }],
    ['GET', `${BASE_URL}/health`, null, { tags: { name: 'HealthCheck' } }],
  ]);

  // Validasi response pertama
  check(responses[0], {
    'GET /accounts → not 5xx': (r) => r.status < 500,
  });

  check(responses[1], {
    'GET /health → 200': (r) => r.status === 200,
  });

  sleep(0.3);
}
