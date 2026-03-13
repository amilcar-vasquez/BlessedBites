import { apiGet } from './client';

export type MenuItem = {
  id: number;
  name: string;
  description: string;
  price: number;
  category_id: number;
  image_url?: string;
  is_active?: boolean;
  popular?: boolean;
};

type MenuResponse = {
  items: MenuItem[];
};

export async function fetchMenu(activeOnly = true, categoryId?: number): Promise<MenuItem[]> {
  const params = new URLSearchParams();
  if (activeOnly) {
    params.set('active', 'true');
  }
  if (typeof categoryId === 'number' && Number.isFinite(categoryId)) {
    params.set('category', String(categoryId));
  }

  const query = params.size > 0 ? `?${params.toString()}` : '';
  const payload = await apiGet<MenuResponse>(`/menu${query}`);
  return payload.items.map((item) => ({
    ...item,
    image_url: normalizeImageUrl(item.image_url)
  }));
}

function normalizeImageUrl(raw?: string): string | undefined {
  if (!raw) return raw;

  // DB records may contain Windows-style separators and legacy static prefixes.
  const normalized = raw.replaceAll('\\', '/').trim();

  if (normalized.startsWith('/ui/static/img/')) {
    return normalized.replace('/ui/static/img/', '/');
  }

  if (normalized.startsWith('ui/static/img/')) {
    return `/${normalized.replace('ui/static/img/', '')}`;
  }

  const uploadsIndex = normalized.indexOf('/uploads/');
  if (uploadsIndex >= 0) {
    return normalized.slice(uploadsIndex);
  }

  const relativeUploadsIndex = normalized.indexOf('uploads/');
  if (relativeUploadsIndex >= 0) {
    return `/${normalized.slice(relativeUploadsIndex)}`;
  }

  return normalized;
}
