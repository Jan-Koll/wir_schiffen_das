import type { PageServerLoad } from './$types';
import { configuration_manager_url } from '$lib/endpoints';
export const load: PageServerLoad = async ({ params }) => {
	return {configurations: await fetch(new URL("/v1/configuration", configuration_manager_url), {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})
		.then(async (response) => {
			return response.json()
		}).then(data => {
			if(data===null){
				return []
			}
			return data
		})
		.catch((error) => {
			console.error('Error creating configuration:', error);
			return []
		})
	}
};