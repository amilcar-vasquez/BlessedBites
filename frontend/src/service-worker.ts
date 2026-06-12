/// <reference types="@sveltejs/kit" />
import { build, files, version } from '$service-worker';

declare const self: ServiceWorkerGlobalScope;

const CACHE = `bb-cache-${version}`;
const ASSETS = [...build, ...files, '/'];

self.addEventListener('install', (event) => {
  // Activate the new worker immediately instead of waiting for all
  // old tabs to close — otherwise an outdated worker keeps serving
  // stale cached API data indefinitely.
  self.skipWaiting();
  event.waitUntil(caches.open(CACHE).then((cache) => cache.addAll(ASSETS)));
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    (async () => {
      // Drop caches from older versions.
      const keys = await caches.keys();
      await Promise.all(keys.filter((key) => key !== CACHE).map((key) => caches.delete(key)));
      // Purge any dynamic responses (API/uploads) cached by previous
      // service-worker versions so fresh data is always fetched.
      const cache = await caches.open(CACHE);
      const cachedRequests = await cache.keys();
      await Promise.all(
        cachedRequests
          .filter((req) => {
            const path = new URL(req.url).pathname;
            return path.startsWith('/api/') || path.startsWith('/uploads/');
          })
          .map((req) => cache.delete(req))
      );
      await self.clients.claim();
    })()
  );
});

self.addEventListener('fetch', (event) => {
  const req = event.request;
  if (req.method !== 'GET') return;

  const url = new URL(req.url);
  // Never intercept dynamic data: API responses and uploaded images must
  // always come from the network so menu/category changes show immediately.
  if (url.origin !== self.location.origin) return;
  if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/uploads/')) return;

  // Cache-first for static assets and pages.
  event.respondWith(
    caches.match(req).then((cached) => {
      if (cached) return cached;
      return fetch(req)
        .then((res) => {
          const copy = res.clone();
          caches.open(CACHE).then((cache) => cache.put(req, copy));
          return res;
        })
        .catch(() => caches.match('/'));
    })
  );
});
