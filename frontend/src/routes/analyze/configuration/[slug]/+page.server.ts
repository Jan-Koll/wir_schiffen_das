import type { PageServerLoad } from './$types';
import { configuration_manager_url, cooling_exhaust_analyzer_url, engine_control_analyzer_url, mounting_analyzer_url, propulsion_analyzer_url, supply_analyzer_url } from '$lib/endpoints';
export const load: PageServerLoad = async ({ params }) => {
	const configuration_manager_promise = await fetch(new URL(`/v1/configuration/${params.slug}`, configuration_manager_url), {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})
		.then(async (response) => {
			return response.json()
		})
		.catch((error) => {
			console.error('Error creating configuration:', error);
			//return new Response(JSON.stringify({ "error": error }), { status: 500 });
			return {}
		});

	const cooling_exhaust_analyzer_promise = await fetch(new URL(`/v1/configuration/${params.slug}/jobs`, cooling_exhaust_analyzer_url), {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})
		.then(async (response) => {
			return response.json()
		})
		.catch((error) => {
			console.error('Error creating configuration:', error);
			//return new Response(JSON.stringify({ "error": error }), { status: 500 });
			return []
		});
	const engine_control_analyzer_promise = await fetch(new URL(`/v1/configuration/${params.slug}/jobs`, engine_control_analyzer_url), {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})
		.then(async (response) => {
			return response.json()
		})
		.catch((error) => {
			console.error('Error creating configuration:', error);
			//return new Response(JSON.stringify({ "error": error }), { status: 500 });
			return []
		});
	const mounting_analyzer_promise = await fetch(new URL(`/v1/configuration/${params.slug}/jobs`, mounting_analyzer_url), {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})
		.then(async (response) => {
			return response.json()
		})
		.catch((error) => {
			console.error('Error creating configuration:', error);
			//return new Response(JSON.stringify({ "error": error }), { status: 500 });
			return []
		});

	const propulsion_analyzer_promise = await fetch(new URL(`/v1/configuration/${params.slug}/jobs`, propulsion_analyzer_url), {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})
		.then(async (response) => {
			return response.json()
		})
		.catch((error) => {
			console.error('Error creating configuration:', error);
			//return new Response(JSON.stringify({ "error": error }), { status: 500 });
			return []
		});
	const supply_analyzer_promise = await fetch(new URL(`/v1/configuration/${params.slug}/jobs`, supply_analyzer_url), {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})
		.then(async (response) => {
			return response.json()
		})
		.catch((error) => {
			console.error('Error creating configuration:', error);
			//return new Response(JSON.stringify({ "error": error }), { status: 500 });
			return []
		});
	const result = await Promise.all([configuration_manager_promise, engine_control_analyzer_promise, cooling_exhaust_analyzer_promise, mounting_analyzer_promise, propulsion_analyzer_promise, supply_analyzer_promise])
	return {
		configuration: result[0],
		engine_control_analyzer: result[1],
		cooling_exhaust_analyzer: result[2],
		mounting_analyzer: result[3],
		propulsion_analyzer: result[4],
		supply_analyzer: result[5]
	}
};