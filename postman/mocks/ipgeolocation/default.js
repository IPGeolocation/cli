const http = require('node:http');
const { Readable } = require('node:stream');

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
  if (url.pathname !== '/getip' && !['TEST_KEY', 'TEST_KEY &+?=#%', 'SLOW_HEADERS', 'SLOW_BODY', 'LIMIT_RESPONSE', 'BULK_50K'].includes(key)) return reply(401, { message: 'Invalid test key' });
  let payload = null;
  if (isPost) {
    if (req.headers['content-type'] !== 'application/json') return reply(400, { message: 'Expected JSON' });
    payload = req.body;
    if (!payload || typeof payload !== 'object') return reply(400, { message: 'Invalid JSON payload' });
  }
  const query = Object.fromEntries([...url.searchParams.keys()].map((name) => [name, url.searchParams.getAll(name)]));
  const result = {
    ip: '2001:db8::1', location: { city: url.searchParams.get('lang') === 'ar' ? 'مدينة' : 'Test City' },
    _request: { method: req.method, path: url.pathname, query, payload },
  };
  if (key === 'BULK_50K') {
    if (!Array.isArray(payload?.ips) || payload.ips.length !== 50000) return reply(400, { message: 'Expected 50000 fixture addresses' });
    const companyName = 'x'.repeat(4000);
    const records = Readable.from((function* () {
      yield '[';
      for (let i = 0; i < payload.ips.length; i++) {
        if (i) yield ',';
        yield JSON.stringify({ ip: payload.ips[i], location: result.location, company: { name: companyName } });
      }
      yield ']';
    })());
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.once('close', () => records.destroy());
    records.pipe(res);
    return;
  }
  if (key === 'LIMIT_RESPONSE') result.padding = 'x'.repeat(1100000);
  if (key === 'SLOW_HEADERS' || key === 'SLOW_BODY') {
    const encoded = JSON.stringify(url.pathname.endsWith('-bulk') ? [result] : result);
    if (key === 'SLOW_BODY') {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.write(encoded.slice(0, 1));
    }
    const timer = setTimeout(() => {
      if (key === 'SLOW_HEADERS') res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(key === 'SLOW_BODY' ? encoded.slice(1) : encoded);
    }, 750);
    res.once('close', () => clearTimeout(timer));
    return;
  }
  reply(200, url.pathname.endsWith('-bulk') ? [result] : result);
}).listen(process.env.PORT || 4500);
