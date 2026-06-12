import { apiGet, apiPost } from './client';

export type AverageRating = {
  menu_item_id: number;
  average: number;
};

export async function getAverageRating(menuItemId: number): Promise<AverageRating> {
  return apiGet<AverageRating>(`/ratings/${menuItemId}`);
}

export async function submitRating(
  menuItemId: number,
  rating: number,
  userId?: number
): Promise<{ message: string }> {
  return apiPost<{ message: string }>('/ratings', {
    menu_item_id: menuItemId,
    rating,
    ...(userId ? { user_id: userId } : {})
  });
}
