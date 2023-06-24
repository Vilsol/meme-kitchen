import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { browser } from "$app/environment";
import { PUBLIC_BACKEND_API_BASE } from "$env/static/public";

export const prerender = false;

export const load = (async ({ params, fetch }) => {
	if (!browser) {
		return {};
	}

  return fetch(`${PUBLIC_BACKEND_API_BASE}/api/templates/${params.id}`)
		.then((response) => response.json())
		.catch((e) => {
      console.error(e);
			throw error(404);
		});
}) satisfies PageLoad;