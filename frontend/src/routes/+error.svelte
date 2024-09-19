<script>
	import { page } from '$app/stores';
	import Disclaimer from '../lib/components/Disclaimer.svelte';
</script>

{#if $page.status === 404}
	<Disclaimer />
{/if}
<h1>{$page.status}: {$page.error.message}</h1>
