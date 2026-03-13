import { apiGet } from './client';

export type Category = {
  id: number;
  name: string;
};

type CategoriesResponse = {
  items: Category[];
};

export async function fetchCategories(): Promise<Category[]> {
  const payload = await apiGet<CategoriesResponse>('/categories');
  return payload.items;
}
