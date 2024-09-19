import { cachedEmitters } from '$lib/state/sse';
import { error } from '@sveltejs/kit';

export async function POST({ request, params }) {
    const { key } = params;
    if (!key) {
        error(400, 'Bad request, no key provided.')
    }
    const { message, sender } = await request.json();
    if (!message) {
        error(400, 'Bad request, no message provided')
    }
    if(sender){
        // append sender information to emitted message
        message.sender = sender
    }
    const writeToAll: boolean = key.toLowerCase() === "all"
    if (writeToAll) {
        const emitters = Array.from(cachedEmitters.values());
        if (!emitters) {
            error(404, 'No emitters found.');
        }
        for (const emit of emitters) {
            emit('message', JSON.stringify({message}));
        }
    } else {
        const emit = cachedEmitters.get(key);
        if (!emit) {
            error(401, 'Unauthorized, provided emitter key is not valid.');
        }

        emit('message', JSON.stringify(message));
    }
    const response = new Response("", {status: 200});

    return response;
}

