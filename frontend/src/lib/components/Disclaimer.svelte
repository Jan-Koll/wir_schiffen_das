<script lang="ts">
	import CircleAlert from 'lucide-svelte/icons/circle-alert';
	import * as Alert from '$lib/components/ui/alert';
	import { Button } from '$lib/components/ui/button/index.js';
</script>

<Alert.Root class="bg-red-200 dark:bg-red-800">
	<CircleAlert class="h-4 w-4 stroke-red-600 dark:stroke-white" />
	<Alert.Title class="text-red-600 dark:text-white">Achtung!</Alert.Title>
	<Alert.Description class="text-black-700 dark:text-red-200">
		Die Startseite sowie die Seiten <code>Orders</code>, <code>Products</code> und
		<code>Customers</code>
		sind nur nur Beispiele bzw. Platzhalter, um zusätzlichen Kontext zum Unternehmen WirSchiffenDas hinzuzufügen.<br
		/>
		<b class="text-bold"
			>Die eigentliche Projektarbeit befindet sich unter
			<a class="text-black underline dark:text-white" href="/analyze">Analyze</a>.</b
		>
		<br />
		<Button class="mt-4" size="sm" href="/analyze">Go to Analyze</Button>
	</Alert.Description>
</Alert.Root>
