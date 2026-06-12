import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params }) => {
  const id = Number(params.category);
  redirect(308, Number.isFinite(id) ? `/menu?category=${id}` : '/menu');
};
