import { engine_control_analyzer_url } from '$lib/endpoints';
export async function POST({ request, params }) {
    const body = await request.text();
    const result = await fetch(new URL("/v1/job", engine_control_analyzer_url), {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body
    })
        .catch((error) => {
            console.error('Could not send request', error);
            return new Response(JSON.stringify({ "error": error }), { status: 500 });
        });
    return result
}

