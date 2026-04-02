// Smoke Test — Validasi dasar bahwa API berfungsi tanpa error
// Jalankan: k6 run scripts/smoke_test.js

import http from 'k6/http';
import { check, sleep } from 'k6';

// Konfigurasi: 1 VU, 10 detik — cukup untuk memastikan API tidak crash
export const options = {
  vus: 1,
  duration: '10s',

  // Thresholds: batas minimal yang harus dipenuhi
  thresholds: {
    http_req_failed: ['rate<0.01'],        // error rate < 1%
    http_req_duration: ['p(95)<500'],       // 95% request < 500ms
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Test 1: GET /accounts — ambil semua akun
  const accountsRes = http.get(`${BASE_URL}/accounts`);
  check(accountsRes, {
    'GET /accounts status is 200': (r) => r.status === 200,
    'GET /accounts has data': (r) => {
      const body = JSON.parse(r.body);
      return body.status === 'success' && body.data !== null;
    },
  });

  sleep(1); // jeda 1 detik antar iterasi
}
