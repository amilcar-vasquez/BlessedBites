const CACHE_NAME = 'blessed-bites-v1';
const PRECACHE_URLS = [
  '/',
  '/static/css/styles.css',
  '/static/css/responsiveStyles.css',
  '/static/js/app.js',
  '/static/img/BlessedBitesIcon.png',
  '/offline.html'
];

self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => cache.addAll(PRECACHE_URLS))
  );
  self.skipWaiting();
});

self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(keys => Promise.all(
      keys.map(k => {
        if (k !== CACHE_NAME) return caches.delete(k);
      })
    ))
  );
  self.clients.claim();
});

self.addEventListener('fetch', event => {
  // HTML navigation requests: network-first, fallback to cache/offline
  if (event.request.mode === 'navigate') {
    event.respondWith(
      fetch(event.request)
        .then(resp => {
          return resp;
        })
        .catch(() => caches.match('/offline.html'))
    );
    return;
  }

  // For other requests, try cache-first
  event.respondWith(
    caches.match(event.request).then(cached => {
      if (cached) return cached;
      return fetch(event.request).then(resp => {
        // put a copy in cache
        return caches.open(CACHE_NAME).then(cache => {
          cache.put(event.request, resp.clone());
          return resp;
        });
      }).catch(() => {
        // fallback for images
        if (event.request.destination === 'image') {
          return caches.match('/static/img/BlessedBitesIcon.png');
        }
      });
    })
  );
});