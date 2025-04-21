import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
      { duration: '30s', target: 20 },
      { duration: '1m30s', target: 10 },
      { duration: '20s', target: 0 },
    ],
  };
  
const BASE_URL = 'http://localhost:8080';
const HEADERS = { 'Content-Type': 'application/json' };

export default function () {
    const loginResModerator = http.post(`${BASE_URL}/dummyLogin`, JSON.stringify({ 
        role: 'moderator' 
    }), {
        headers: HEADERS,
    });

    check(loginResModerator, {
        'dummyLoginModerator 200': (r) => r.status === 200,
    });

    const moderator_token = loginResModerator.json('token');
    const authHeadersModerator = {
        ...HEADERS,
        Authorization: `Bearer ${moderator_token}`,
    };

    const pvzRes = http.post(`${BASE_URL}/pvz`, JSON.stringify({ city: 'Moscow' }), {
        headers: authHeadersModerator,
    });

    check(pvzRes, {
        'CreatePVZ 201': (r) => r.status === 201 && r.json('city') === 'Moscow',
    });

    const loginResEmployee = http.post(`${BASE_URL}/dummyLogin`, JSON.stringify({
        role: "employee"
    }), {
        headers: HEADERS,
    });

    check(loginResEmployee, {
        'dummyLoginEmployee 200': (r) => r.status === 200,
    });

    const employee_token = loginResEmployee.json('token');
    const authHeadersEmployee = {
        ...HEADERS,
        Authorization: `Bearer ${employee_token}`,
    };
}
