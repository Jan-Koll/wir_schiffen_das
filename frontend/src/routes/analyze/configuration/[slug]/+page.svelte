<script lang="ts">
	import type {
		configuration,
		engine_control_analyzer_response,
		cooling_exhaust_analyzer_response,
		mounting_analyzer_response,
		propulsion_analyzer_response,
		supply_analyzer_response
	} from '$lib/types';

	import { page } from '$app/stores';
	export let data: {
		configuration: configuration;
		engine_control_analyzer: engine_control_analyzer_response[];
		cooling_exhaust_analyzer: cooling_exhaust_analyzer_response[];
		mounting_analyzer: mounting_analyzer_response[];
		propulsion_analyzer: propulsion_analyzer_response[];
		supply_analyzer: supply_analyzer_response[];
	};
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import LoaderCircle from 'lucide-svelte/icons/loader-circle';
	import Play from 'lucide-svelte/icons/play';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { toast } from 'svelte-sonner';
	import { messages } from '$lib/state/sse';
	import StatusIndicator from '$lib/components/StatusIndicator.svelte';
	import {capitalize} from '$lib/strings';
	import {browser} from "$app/environment"
	let awaitingResponse = false;
	
	let latestOrderCreatedAt : Date | null = null
	let maxAmountOrders = 0

	type analysis = Omit<typeof data, "configuration">
	const latestResult : analysis = {
		"engine_control_analyzer":{},
		"cooling_exhaust_analyzer": {}, 
		"mounting_analyzer":{},
		"propulsion_analyzer":{}, 
		"supply_analyzer":{},
	}
	let totalState : string | undefined = undefined 
	let totalSuccess : boolean = false
	let groupedResults = groupResults()
	if(browser)	console.log("grouped", groupedResults)
	function groupResults(){
		const result :any = {}
		for (const [analyzer, analyses] of Object.entries(data)) {
			if(analyzer==="configuration") continue;
			for(const analysis of Object.values(analyses)){
				const order_created_at = analysis.order_created_at;
				if(!result[order_created_at]){
					result[order_created_at] = {}
				}
				if(!result[order_created_at][analyzer]){
					result[order_created_at][analyzer] = {}
				}
				result[order_created_at][analyzer] = analysis
			}
		}
		const sortedResult = Object.fromEntries(
			Object.entries(result).sort((a, b) => new Date(b[0]) - new Date(a[0]))
		);

		return sortedResult
	}
	function getLatestResult(){
		for(const analyzerType in latestResult){
			for(const analyzer of data[analyzerType]){
				const timestamp = new Date(analyzer.order_created_at).getTime()
				if(latestOrderCreatedAt === null || timestamp > latestOrderCreatedAt.getTime()){
					latestOrderCreatedAt = new Date(timestamp);
				}
			}
			if(data[analyzerType].length > maxAmountOrders){
				maxAmountOrders = data[analyzerType].length
			}
		}
		for(const analyzerType in latestResult){
			latestResult[analyzerType] = data[analyzerType].filter((analyzer:any) => new Date(analyzer.order_created_at).getTime() === latestOrderCreatedAt?.getTime())[0] ?? {}
		}
		if(browser) console.log("latestResult", latestResult)
		let totalResult = getTotalState(latestResult)
		totalState = totalResult[0]
		totalSuccess = totalResult[1]
		
		groupedResults = groupResults()
	}
	getLatestResult()

	if(browser) console.log('loaded data', data);

	function getTotalState(analyses:analysis): [string, boolean] {
		const results = Object.values(analyses)
		let state: string, success:boolean = false
		if(results.every(result => result.status === "ready")){
			state = "ready"
			success = results.every(result => result.success)
		} else if(results.findIndex(result => result.status === "failed") > -1){
			state = "failed"
		} else if(results.findIndex(result => result.status === "canceled") > -1){
			state = "canceled"
		} else if(results.findIndex(result => result.status === "running") > -1){
			state = "running"
		} else if(results.findIndex(result => result.status === "queued") > -1){
			state = "queued"
		} else {
			state = undefined
		}
		return [state, success]
	}

	function updateAnalysisState(parsedMessage: any) {
		if(!Object.keys(latestResult).includes(parsedMessage.sender)) {
			console.error("not in there")
			return
		}

		const idx = data[parsedMessage.sender].findLastIndex(
			(analysis_result) => analysis_result.id === parsedMessage.id
		)
		if(idx === -1){
			//insert new job into data
			data[parsedMessage.sender].push(parsedMessage)
		}else{
			//update existing job in data, but prevent updates from outdated events (which might be caused due to network latency)
			let is_outdated = false
			switch(data[parsedMessage.sender][idx].status){
				case "queued": break;
				case "running": 
					is_outdated = ["queued"].includes(parsedMessage.status)
					break;
				case "failed": 
					is_outdated = ["queued", "running", "canceled", "ready"].includes(parsedMessage.status)
					break;
				case "ready":
					is_outdated = ["queued", "running", "canceled", "failed"].includes(parsedMessage.status)
					break;
				case "canceled":
					is_outdated = ["queued", "running"].includes(parsedMessage.status)
					break;
			}
			if(!is_outdated){
				data[parsedMessage.sender][idx] = parsedMessage
			} else {
				console.log("outdated message")
			}
		}
		data[parsedMessage.sender] = data[parsedMessage.sender] //reassign for reactivity
	
		if(browser) console.log(data)
		getLatestResult()
	}
	$: if ($messages) {
		let parsedMessage;
		console.log('message: ', $messages);
		try {
			parsedMessage = JSON.parse($messages);
			parsedMessage = parsedMessage.message;
			if (!parsedMessage.sender) {
				toast($messages);
			} else if ('configuration_id' in parsedMessage)
				if (parsedMessage.configuration_id === data.configuration.id) {
					//own configuration
					updateAnalysisState(parsedMessage);
				} else {
					// other configurations
					toast($messages);
				}
		} catch (e) {
			console.log('error', e);
			toast($messages);
		}
	}
	async function startAnalyzing() {
		let formData = {
			configuration_id: data.configuration.id,
			order_created_at: new Date().toISOString()
		};
		const body = JSON.stringify(formData);
		awaitingResponse = true;
		fetch($page.url.pathname, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body
		})
			.then((response) => {
				console.log(response);
				if (response.status >= 300) {
					toast.error('Error ' + response.status, {
						description: 'Failed starting the analysis'
					});
					console.log(`Error ${response.status}: Failed starting the analysis.`);
					return false;
				}
				return response.json();
			})
			.then((analysesData) => {
				awaitingResponse = false;
				if (analysesData) {
					toast.info('Analysis added to queue');
					console.log(
						'Analysis successfully added to queue:',
						analysesData,
						data.engine_control_analyzer
					);
				}
			})
			.catch((error) => {
				awaitingResponse = false;
				console.error('Error starting analysis:', error);
				toast.error('Error', { description: 'Could not start analysis: ' + error });
			});
	}
</script>

<svelte:head>
	<title>Configuration {$page.params.slug} | WirSchiffenDas</title>
</svelte:head>
<div class="mx-auto grid max-w-[59rem] flex-1 auto-rows-max gap-4">
	<div class="flex items-center gap-4">
		<Button variant="outline" size="icon" class="h-7 w-7" href="/analyze">
			<ChevronLeft class="h-4 w-4" />
			<span class="sr-only">Back</span>
		</Button>
		<h1 class="flex-1 shrink-0 whitespace-nowrap text-xl font-semibold tracking-tight sm:grow-0">
			View Configuration #{data.configuration.id}
		</h1>
		<div class="hidden items-center gap-2 md:ml-auto md:flex">
			<Button
				type="reset"
				variant="outline"
				size="sm"
				href="/analyze/edit/{data.configuration.id}">Edit</Button
			>
			<Button size="sm" disabled={awaitingResponse} on:click={startAnalyzing}>
				{#if !awaitingResponse}
					<Play class="mr-2 h-4 w-4"></Play>
					Start Analysis
				{:else}
					<LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
					Please wait
				{/if}
			</Button>
		</div>
		<!-- todo: when viewing, show badge with id: -->
	</div>
	<div class="grid gap-4 md:grid-cols-[1fr_250px] lg:grid-cols-3 lg:gap-8">
		<div class="grid-auto-rows-max grid items-start gap-4 lg:col-span-2 lg:gap-8">
			<Card.Root>
				<Card.Header>
					<Card.Title>General</Card.Title>
					<Card.Description
						>General information about the configuration and the used engine</Card.Description
					>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={data.configuration.description === ''}>
							<Label>Description <span class="text-xs">(required)</span></Label>
							{#if data.configuration.description !== ''}
								{data.configuration.description}
							{:else}
								<span class="text-red-500">No description</span>
							{/if}
						</div>
						<!--<div class="flex flex-col gap-6 justify-normal space-y-6 sm:flex-row sm:space-y-0">-->
						<div class="grid gap-8 sm:grid-cols-2">
							<div class="w grid gap-3" class:opacity-60={data.configuration.engine === ''}>
								<Label>Engine <span class="text-xs">(required)</span></Label>
								{#if data.configuration.engine !== ''}
									{data.configuration.engine}
								{:else}
									<span class="text-red-500">No engine</span>
								{/if}
							</div>
							<div class="grid gap-3" class:opacity-60={data.configuration.gearbox_type === ''}>
								<Label for="gearboxType">Gearbox type <span class="text-xs">(required)</span></Label
								>
								{#if data.configuration.gearbox_type !== ''}
									{data.configuration.gearbox_type}
								{:else}
									<span class="text-red-500">No gearbox type</span>
								{/if}
							</div>
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Engine Control</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={!data.configuration.engine_management_system}>
							<Label>Engine Management System</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									id="engineManagementSystem"
									bind:checked={data.configuration.engine_management_system}
									aria-labelledby="engineManagementSystemLabel"
									disabled
								/>
								<Label
									id="engineManagementSystemLabel"
									for="engineManagementSystem"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
								>
									Compliance with Classification Society Regulations
								</Label>
							</div>
						</div>
						<div
							class="grid gap-3"
							class:opacity-60={data.configuration.monitoring_control_system === ''}
						>
							<Label for="monitoringControlSystem">Monitoring/Control System</Label>
							{#if data.configuration.monitoring_control_system !== ''}
								{data.configuration.monitoring_control_system}
							{:else}
								None
							{/if}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Propulsion</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={!data.configuration.starting_system}>
							<Label>Starting System</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									disabled
									id="startingSystem"
									bind:checked={data.configuration.starting_system}
									aria-labelledby="startingSystemLabel"
								/>
								<Label
									id="startingSystemLabel"
									for="startingSystem"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
									>Air starter
								</Label>
							</div>
						</div>
						<div class="grid gap-3" class:opacity-60={!data.configuration.power_transmission}>
							<Label>Power Transmission</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									id="powerTransmission"
									disabled
									bind:checked={data.configuration.power_transmission}
									aria-labelledby="powerTransmissionLabel"
								/>
								<Label
									id="powerTransmissionLabel"
									for="powerTransmission"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
									>Torsionally resilient coupling</Label
								>
							</div>
						</div>
						<div
							class="grid gap-3"
							class:opacity-60={!data.configuration.auxiliary_pto ||
								(data.configuration.auxiliary_pto && data.configuration.auxiliary_pto.length === 0)}
						>
							<Label for="auxiliaryPto">Auxiliary PTO</Label>
							{#if data.configuration.auxiliary_pto && data.configuration.auxiliary_pto.length > 0}
								<ul class="ml-5 list-disc">
									{#each data.configuration.oil_system as item}
										<li>{item}</li>
									{/each}
								</ul>
							{:else}
								None
							{/if}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Supply</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div
							class="grid gap-3"
							class:opacity-60={!data.configuration.oil_system ||
								(data.configuration.oil_system && data.configuration.oil_system.length === 0)}
						>
							<Label for="oilSystems">Oil system</Label>
							{#if data.configuration.oil_system && data.configuration.oil_system.length > 0}
								<ul class="ml-5 list-disc">
									{#each data.configuration.oil_system as item}
										<li>{item}</li>
									{/each}
								</ul>
							{:else}
								None
							{/if}
						</div>
						<div
							class="grid gap-3"
							class:opacity-60={!data.configuration.fuel_system ||
								(data.configuration.fuel_system && data.configuration.fuel_system.length === 0)}
						>
							<Label for="fuelSystems">Fuel system</Label>
							{#if data.configuration.fuel_system && data.configuration.fuel_system.length > 0}
								<ul class="ml-5 list-disc">
									{#each data.configuration.fuel_system as item}
										<li>{item}</li>
									{/each}
								</ul>
							{:else}
								None
							{/if}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Cooling & Exhaust</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={data.configuration.cooling_system === ''}>
							<Label for="coolingSystem">Cooling system</Label>
							{#if data.configuration.cooling_system !== ''}
								{data.configuration.cooling_system}
							{:else}
								None
							{/if}
						</div>
						<div class="grid gap-3" class:opacity-60={!data.configuration.exhaust_system}>
							<Label>Exhaust system</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									id="exhaustSystem"
									disabled
									bind:checked={data.configuration.exhaust_system}
									aria-labelledby="exhaustSystemLabel"
								/>
								<Label
									id="exhaustSystemLabel"
									for="exhaustSystem"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
									>90° Exhaust bellows discharge rotatable</Label
								>
							</div>
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Mounting & Gearbox</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={!data.configuration.mounting_system}>
							<Label>Mounting system</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									disabled
									id="mountingSystem"
									bind:checked={data.configuration.mounting_system}
									aria-labelledby="mountingSystemLabel"
								/>
								<Label
									id="mountingSystemLabel"
									for="mountingSystem"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
									>Resilient mounts at driving end</Label
								>
							</div>
						</div>
						<div
							class="grid gap-3"
							class:opacity-60={!data.configuration.gearbox_option ||
								(data.configuration.gearbox_option &&
									data.configuration.gearbox_option.length === 0)}
						>
							<Label for="gearboxOptions">Gearbox options</Label>
							{#if data.configuration.gearbox_option && data.configuration.gearbox_option.length > 0}
								<ul class="ml-5 list-disc">
									{#each data.configuration.gearbox_option as item}
										<li>{item}</li>
									{/each}
								</ul>
							{:else}
								None
							{/if}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
		</div>

		<div class="grid auto-rows-max items-start gap-4 lg:gap-8">
			<Card.Root>
				<Card.Header>
					<Card.Title class="flex content-center items-center justify-between"
						>Status #{maxAmountOrders}
						<StatusIndicator
							status={totalState}
							success={totalSuccess}
							name="Total result"
							side="right"
						/></Card.Title
					>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3">
							{#each Object.keys(latestResult) as analyzer}
								{#if data[analyzer].length > 0 }
									<div class="flex gap-3">
										<StatusIndicator
											status={latestResult[analyzer].status}
											success={latestResult[analyzer].success}
											name="{capitalize(analyzer, ["_"])} #{maxAmountOrders}"
											side="left"
										></StatusIndicator>
										{capitalize(analyzer, ["_"])}
									</div>
								{/if}
							{/each}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root class="overflow-hidden">
				<Card.Header>
					<Card.Title>Result history</Card.Title>
					<Card.Description>Results of past analyses of this configuration</Card.Description>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-2">
						{#each Object.values(groupedResults) as analysis, index}
							<div
								class:dark:border-emerald-800="{getTotalState(analysis)[1]}" 
								class:dark:bg-emerald-950="{getTotalState(analysis)[1]}" 
								class:border-emerald-400="{getTotalState(analysis)[1]}"
								class="flex items-center justify-start gap-x-2 rounded border p-1 dark:border-emerald-800"
							>
								<span class="mx-2 text-muted-foreground min-w-8 text-right">#{Object.values(groupedResults).length - index}</span>
								<!-- {JSON.stringify(analysis)} -->
								 <div class="mt-1 -mb-1">
								{#each Object.entries(analysis) as [analyzer, analyzer_result]}
									<StatusIndicator 
										status={analyzer_result.status} 
										success={analyzer_result.success}
										side="top"
										name="{capitalize(analyzer, "_")} #{Object.values(groupedResults).length - index}"
									></StatusIndicator>
								{/each}
								</div>
							</div>
						{/each}
						
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	</div>
</div>
