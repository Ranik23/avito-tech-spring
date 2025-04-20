import http from 'k6/http';
import { check } from 'k6';

export const options = {
    scenarios: {
        moderate_load: {
            executor: 'constant-arrival-rate',
            rate: 50,
            timeUnit: '1s',
            duration: '30s',
            preAllocatedVUs: 10,
            maxVUs: 50,
        },
    },
};

const BASE_URL = 'http://localhost:8080';
const HEADERS = { 'Content-Type': 'application/json' };

export default function () {
    const loginRes = http.post(`${BASE_URL}/dummyLogin`, JSON.stringify({ 
        role: 'moderator' 
    }), {
        headers: HEADERS,
    });

    check(loginRes, {
        'dummyLogin 200': (r) => r.status === 200,
    });

    const token = loginRes.json('token');
    const authHeaders = {
        ...HEADERS,
        Authorization: `Bearer ${token}`,
    };

    const pvzRes = http.post(`${BASE_URL}/pvz`, JSON.stringify({ city: 'Moscow' }), {
        headers: authHeaders,
    });

    check(pvzRes, {
        'CreatePVZ 201': (r) => r.status === 201 && r.json('city') === 'Moscow',
    });

    const pvzId = pvzRes.json('id');

    const productRes = http.post(`${BASE_URL}/products`, JSON.stringify({
        pvzId: pvzId,
        type: "box",
    }), {
        headers: authHeaders,
    });

    check(productRes, {
        'AddProduct 201': (r) => r.status === 201,
    });
}
