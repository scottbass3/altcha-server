import { check, group } from 'k6';
import http from 'k6/http'

export const options = {
  scenarios: {
    load: {
      vus: 1,
      iterations: 1,
      executor: 'per-vu-iterations',
      options: {
        browser: {
          type: 'chromium',
        },
      },
    },
  }
};

http.setResponseCallback(
    http.expectedStatuses(200)
);

export default function () {

    let response;

    group('Request challenge', function () {

        response = http.get(`${__ENV.LOAD_TEST_URL}request`);
        check(response, {
            'is status 200': (r) => r.status === 200, 
            'Contenu souhaité': r => r.body.includes('algorithm') && r.body.includes('challenge') && r.body.includes('maxNumber') && r.body.includes('salt') && r.body.includes('signature'),
        });

    })

    group('Verify challenge', function () {

        const data = {"challenge":"d5656d52c5eadce5117024fbcafc706aad397c7befa17804d73c992d966012a8","salt":"8ec1b7ed694331baeb7416d9?expires=1727963398","signature":"781014d0a7ace7e7ae9d12e2d5c0204b60a8dbf42daa352ab40ab582b03a9dc6","number":219718,};
        response = http.post(`${__ENV.LOAD_TEST_URL}verify`, JSON.stringify(data),
        {
          headers: { 'Content-Type': 'application/json' },
        });
        check(response, {
            'is status 200': (r) => r.status === 200,
            'Challenge verified': (r) => r.body.includes('true'),
        });

    })
}