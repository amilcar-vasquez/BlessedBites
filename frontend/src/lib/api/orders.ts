import { API_BASE, apiPost } from './client';

// NOTE: the backend handler (backend/internal/handlers/orders.go) expects
// item keys `id` and `qty` — not the openapi example's menu_item_id/quantity.
export type OrderItemInput = {
  id: number;
  qty: number;
};

export type CreateOrderInput = {
  items: OrderItemInput[];
  full_name: string;
  phone_no: string;
};

export type CreateOrderResponse = {
  order_id: number;
  message: string;
};

export async function createOrder(input: CreateOrderInput): Promise<CreateOrderResponse> {
  return apiPost<CreateOrderResponse>('/orders', input);
}

export type OrderCreatedEvent = {
  order_id: number;
  user_id: number;
  total_cost: number;
};

export type OrdersStreamHandlers = {
  onOrderCreated?: (event: OrderCreatedEvent) => void;
  onOpen?: () => void;
  onError?: () => void;
};

/** Subscribe to the live SSE order feed. Returns an unsubscribe function. */
export function subscribeOrdersStream(handlers: OrdersStreamHandlers): () => void {
  const source = new EventSource(`${API_BASE}/orders/stream`, { withCredentials: true });

  source.addEventListener('open', () => handlers.onOpen?.());
  source.addEventListener('error', () => handlers.onError?.());
  source.addEventListener('order.created', (evt) => {
    try {
      const payload = JSON.parse((evt as MessageEvent).data) as OrderCreatedEvent;
      handlers.onOrderCreated?.(payload);
    } catch {
      // ignore malformed event payloads
    }
  });

  return () => source.close();
}
