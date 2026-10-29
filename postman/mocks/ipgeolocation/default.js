const http = require('node:http');

const endpoints = new Set([
  '/getip', '/v3/ipgeo', '/v3/security', '/v3/abuse', '/v3/asn', '/v3/timezone',
  '/v3/timezone/convert', '/v3/astronomy', '/v3/astronomy/timeSeries',
  '/v3/ipgeo-bulk', '/v3/security-bulk', '/v3/user-agent', '/v3/user-agent-bulk',
]);

http.createServer((req, res) => {
  const url = new URL(req.url, 'http://localhost');
  const reply = (status, body) => { res.writeHead(status, { 'Content-Type': 'application/json' }); res.end(JSON.stringify(body)); };
  if (!endpoints.has(url.pathname)) return reply(501, { message: 'Endpoint not modeled' });
  const isPost = url.pathname.endsWith('-bulk') || url.pathname === '/v3/user-agent';
  if (req.method !== (isPost ? 'POST' : 'GET')) return reply(400, { message: 'Unexpected HTTP method' });
  const key = url.searchParams.get('apiKey');
  const statuses = { BAD_KEY: 401, FORBIDDEN: 403, RATE_LIMIT: 429, SERVER_ERROR: 503 };
  if (statuses[key]) return reply(statuses[key], { message: `Test failure for ${key}` });
  if (key === 'BAD_JSON') { res.writeHead(200); return res.end('{'); }
  if (url.pathname !== '/getip' && !['TEST_KEY', 'TEST_KEY &+?=#%'].includes(key)) return reply(401, { message: 'Invalid test key' });
  let payload = null;
  if (isPost) {
    if (req.headers['content-type'] !== 'application/json') return reply(400, { message: 'Expected JSON' });
    payload = req.body;
    if (!payload || typeof payload !== 'object') return reply(400, { message: 'Invalid JSON payload' });
  }
  const query = Object.fromEntries([...url.searchParams.keys()].map((name) => [name, url.searchParams.getAll(name)]));
  const result = {
    ip: '2001:db8::1', location: { city: 'Test City' },
    _request: { method: req.method, path: url.pathname, query, payload },
  };
  reply(200, url.pathname.endsWith('-bulk') ? [result] : result);
}).listen(process.env.PORT || 4500);
