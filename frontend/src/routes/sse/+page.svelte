<script>
	let name = 'world';
	import { writable } from 'svelte/store';
	import { onMount } from 'svelte';
	const messages = writable([]);
    let events = ["asd"]
	onMount(() => {
		const source = new EventSource(`/test`, {
			withCredentials: false
		});
		source.onmessage = (event) => {
            console.log(event)
            events = [...events, event.data]
		};
        source.addEventListener("userconnect", (event)=>{
            events = [...events, event.data]
            console.log("USERCONNECT", event.data)
        });
        source.addEventListener("usermessage", (event)=>{
            events = [...events, event.data]
            console.log("USERMESSAGE", event.data)
        });
		return () => {
			source.close();
		};
	});
</script>

<h1>Hello {name}!</h1>


{#each events as e}
	<p>
		{e}
	</p>
{/each}
