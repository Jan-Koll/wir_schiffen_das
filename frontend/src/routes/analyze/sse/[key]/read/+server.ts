import { produce } from 'sveltekit-sse'
import { delay } from '$lib/delay';
import { cachedEmitters, cachedKeys } from '$lib/state/sse';
//eslint-disable-next-line
export async function GET({ request, params }) {
  const { key } = params;

  return produce({
    async start({ emit, lock }) {
      const cachedKey = cachedKeys.get(key) ?? false;

      // Remove the key as soon as possible.
      cachedKeys.get(key);

      if (!cachedKey) {
        // Stop the stream if the key is not valid.
        lock.set(false);
        return;
      }

      cachedEmitters.set(key, emit);
      emit('ready', 'true');
      console.log("ready")

      // eslint-disable-next-line no-constant-condition
      while (true) {
        emit('message', `some message ${Date.now()}, ${JSON.stringify(Array.from(cachedKeys))}`);
        await delay(2000);
      }
    },
    stop() {
      // When this stream ends, remove the emitter it cached,
      // so that subsequent notifications POST requests will fail.
      cachedEmitters.get(key);
    },
  });
}

//eslint-disable-next-line
export async function POST({ request, params }) {
  const { key } = params;
  return produce(
    async function start({ emit, lock }) {
      const cachedKey = cachedKeys.get(key) ?? false;
      if (!cachedKey) {
        // Stop the stream if the key is not valid.
        lock.set(false);
        return;
      }
      cachedEmitters.set(key, emit);
      emit('ready', 'true');
      //eslint-disable-next-line no-constant-condition
      while (true) {
        //emit('message', `New test message at ${new Date().toLocaleString("uk")}`);
        await delay(10000);
      }
    },
    {
      stop(){
      // When this stream ends, remove the emitter it cached,
      // so that subsequent notifications POST requests will fail.
      cachedEmitters.get(key);
      }
    }
  )
}