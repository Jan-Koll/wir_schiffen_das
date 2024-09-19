<script lang="ts">
	import type { configuration } from '$lib/types';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import LoaderCircle from 'lucide-svelte/icons/loader-circle';

	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { goto } from '$app/navigation';
	import { writable, derived } from 'svelte/store';
	import type { PageData } from './$types.js';
	import { toast } from 'svelte-sonner';
	export let data: PageData;
	import { page } from '$app/stores';
	let maxAmountOfBadgesInSelect: number = 2;
	type item = {
		value: string;
		label: string;
	};

	const engines = [
		{
			value: '2000 M96 10V',
			label: '2000 M96 10V'
		},
		{
			value: '2000 M96 12V 1',
			label: '2000 M96 12V 1'
		},
		{
			value: '2000 M96 12V 2',
			label: '2000 M96 12V 2'
		},
		{
			value: '2000 M96 16V',
			label: '2000 M96 16V'
		}
	];
	const gearboxTypes: item[] = [
		{
			value: '10V - ZF 2050',
			label: '10V - ZF 2050'
		},
		{
			value: '12V - ZF 2060 1',
			label: '12V - ZF 2060 1'
		},
		{
			value: '12V - ZF 2060 2',
			label: '12V - ZF 2060 2'
		},
		{
			value: '16V - ZF 3060',
			label: '16V - ZF 3060'
		}
	];
	const monitoringControlSystems: item[] = [
		{ value: '', label: 'None' },
		{ value: 'BlueVision', label: 'BlueVision' },
		{
			value: 'BlueVision New Generation',
			label: 'BlueVision New Generation'
		}
	];
	const auxiliaryPtos: item[] = [
		{ value: 'Alternator 120A, 28V, 2 pole', label: 'Alternator 120A, 28V, 2 pole' },
		{ value: 'Alternator 190A, 28V, 2 pole', label: 'Alternator 190A, 28V, 2 pole' },
		{ value: 'Bilgepump', label: 'Bilgepump' },
		{ value: 'On-engine PTOs', label: 'On-engine PTOs' }
	];
	const oilSystems: item[] = [
		{ value: 'Oil replenishment system', label: 'Oil replenishment system' },
		{ value: 'Diverter valve for duplex filter', label: 'Diverter valve for duplex filter' },
		{ value: 'Monitoring fuel leakage', label: 'Monitoring fuel leakage' }
	];
	const fuelSystems: item[] = [
		{ value: 'Duplex fuel pre-filter', label: 'Duplex fuel pre-filter' },
		{ value: 'Diverter valve for fuel filter', label: 'Diverter valve for fuel filter' }
	];
	const coolingSystems: item[] = [
		{ value: '', label: 'None' },
		{
			value: 'Coolant preheating system freestanding',
			label: 'Coolant preheating system freestanding'
		},
		{
			value: 'Coolant preheating system engine mounted',
			label: 'Coolant preheating system engine mounted'
		},
		{ value: 'Integr. seawater gearbox piping', label: 'Integr. seawater gearbox piping' }
	];
	const gearboxOptions: item[] = [
		{ value: 'Reverse reduction gearbox', label: 'Reverse reduction gearbox' },
		{ value: 'El. actuated', label: 'El. actuated' },
		{ value: 'Gearbox mounts', label: 'Gearbox mounts' },
		{
			value: 'Trolling mode for dead-slow-propulsion',
			label: 'Trolling mode for dead-slow-propulsion'
		},
		{ value: 'Free auxiliary PTO', label: 'Free auxiliary PTO' },
		{ value: 'Hydraulic pump drives', label: 'Hydraulic pump drives' }
	];
	let selectedEngine: item | null = null;
	let selectedGearboxType: item | null = null;
	let selectedMonitoringControlSystem: item | null = null;
	let selectedCoolingSystem: item | null = null;

	let selectedAuxiliaryPtos: item[] = [];
	let selectedOilSystems: item[] = [];
	let selectedFuelSystems: item[] = [];
	let selectedGearboxOptions: item[] = [];

	let isAirStarterCompatible: boolean = false;

	let initialFormData: configuration = {
		description: '',
		engine: '',
		gearbox_type: '',
		monitoring_control_system: '',
		cooling_system: '',
		auxiliary_pto: [],
		oil_system: [],
		fuel_system: [],
		gearbox_option: [],
		engine_management_system: false,
		starting_system: false,
		power_transmission: false,
		exhaust_system: false,
		mounting_system: false
	};
	let formData = writable(structuredClone(initialFormData));
	const formIsValid = derived(formData, ($formData) => {
		let requiredFields = [$formData.description, $formData.engine, $formData.gearbox_type];
		//console.log($formData);
		return requiredFields.includes('');
	});
	let awaitingResponse = false;
	$: isAirStarterCompatible = ['12V', '16V'].some(
		(compatibleEngine) => selectedEngine?.value.includes(compatibleEngine) ?? false
	);
	$: if (!isAirStarterCompatible) {
		$formData.starting_system = false;
	}
	$: if (selectedAuxiliaryPtos) {
		$formData.auxiliary_pto = getValuesFromItemList(selectedAuxiliaryPtos);
	}
	$: if (selectedOilSystems) {
		$formData.oil_system = getValuesFromItemList(selectedOilSystems);
	}
	$: if (selectedGearboxOptions) {
		$formData.gearbox_option = getValuesFromItemList(selectedGearboxOptions);
	}
	$: if (selectedFuelSystems) {
		$formData.fuel_system = getValuesFromItemList(selectedFuelSystems);
	}
	function getValuesFromItemList(itemList: item[]): string[] {
		return itemList.map((item) => item.value);
	}
	function resetForm() {
		console.log(initialFormData);
		$formData = initialFormData;
		selectedEngine = null;
		selectedGearboxType = null;
		selectedMonitoringControlSystem = null;
		selectedCoolingSystem = null;

		selectedAuxiliaryPtos = [];
		selectedOilSystems = [];
		selectedFuelSystems = [];
		selectedGearboxOptions = [];
	}
	async function submit() {
		const body = JSON.stringify($formData);
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
						description: 'Failed creating configuration.'
					});
					console.log(`Error ${response.status}: Failed creating configuration.`);
					return false;
				}
				return response.json();
			})
			.then((data) => {
				console.log(data);
				awaitingResponse = false;
				if (data) {
					toast.success('Configuration saved');
					console.log('Configuration created successfully:', data);
					goto(`/analyze/configuration/${data.id}`);
				}
			})
			.catch((error) => {
				awaitingResponse = false;
				console.error('Error creating configuration:', error);
				toast.error('Error', { description: 'Could not store configuration: ' + error });
			});
	}
</script>

<svelte:head>
	<title>Create Configuration | WirSchiffenDas</title>
</svelte:head>
<form method="POST">
	<div class="mx-auto grid max-w-[59rem] flex-1 auto-rows-max gap-4">
		<div class="flex items-center gap-4">
			<Button variant="outline" size="icon" class="h-7 w-7" href="/analyze">
				<ChevronLeft class="h-4 w-4" />
				<span class="sr-only">Back</span>
			</Button>
			<h1 class="flex-1 shrink-0 whitespace-nowrap text-xl font-semibold tracking-tight sm:grow-0">
				Create Configuration
			</h1>
			<!-- todo: when viewing, show badge with id: <Badge variant="outline" class="ml-auto sm:ml-0">#id</Badge>-->
		</div>
		<div class="grid gap-4 md:grid-cols-[1fr_250px] lg:grid-cols-3 lg:gap-8">
			<div class="col-span-full grid auto-rows-max items-start gap-4 lg:gap-8">
				<Card.Root>
					<Card.Header>
						<Card.Title>General</Card.Title>
						<Card.Description
							>General information about the configuration and the used engine</Card.Description
						>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6">
							<div class="grid gap-3">
								<Label for="description">Description <span class="text-xs">(required)</span></Label>
								<Input
									id="description"
									type="text"
									class="w-full"
									placeholder="A brief description what you analyze"
									bind:value={$formData.description}
								/>
							</div>
							<div class="flex flex-col justify-between gap-4 space-y-6 sm:flex-row sm:space-y-0">
								<div class="grid gap-3">
									<Label for="engine">Engine <span class="text-xs">(required)</span></Label>

									<Select.Root
										portal={null}
										multiple={false}
										bind:selected={selectedEngine}
										onSelectedChange={(v) => {
											v && ($formData.engine = v.value);
										}}
									>
										<Select.Trigger class="w-[250px] md:w-[300px] lg:w-[400px]">
											<div class="flex items-center justify-start">
												{#if !selectedEngine}
													Select an engine
												{:else}
													{selectedEngine.label}
												{/if}
											</div>
										</Select.Trigger>
										<Select.Content>
											<Select.Group>
												{#each engines as engine}
													<Select.Item value={engine.value} label={engine.label}
														>{engine.label}</Select.Item
													>
												{/each}
											</Select.Group>
										</Select.Content>
										<Select.Input id="engine" name="engine" />
									</Select.Root>
								</div>
								<div class="grid gap-3">
									<Label for="gearboxType"
										>Gearbox type <span class="text-xs">(required)</span></Label
									>
									<Select.Root
										portal={null}
										multiple={false}
										bind:selected={selectedGearboxType}
										onSelectedChange={(v) => {
											v && ($formData.gearbox_type = v.value);
										}}
									>
										<Select.Trigger class="w-[250px] md:w-[300px] lg:w-[400px]">
											<div class="flex items-center justify-start">
												{#if !selectedGearboxType}
													Select a gearbox type
												{:else}
													{selectedGearboxType.label}
												{/if}
											</div>
										</Select.Trigger>
										<Select.Content>
											<Select.Group>
												{#each gearboxTypes as gearboxType}
													<Select.Item value={gearboxType.value} label={gearboxType.label}
														>{gearboxType.label}</Select.Item
													>
												{/each}
											</Select.Group>
										</Select.Content>
										<Select.Input id="gearboxType" name="gearboxType" />
									</Select.Root>
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
							<div class="grid gap-3">
								<Label>Engine Management System</Label>
								<div class="flex items-center space-x-2">
									<Checkbox
										id="engineManagementSystem"
										bind:checked={$formData.engine_management_system}
										aria-labelledby="engineManagementSystemLabel"
									/>
									<Label
										id="engineManagementSystemLabel"
										for="engineManagementSystem"
										class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
									>
										Compliance with Classification Society Regulations
									</Label>
								</div>
							</div>
							<div class="grid gap-3">
								<Label for="monitoringControlSystem">Monitoring/Control System</Label>
								<Select.Root
									portal={null}
									multiple={false}
									bind:selected={selectedMonitoringControlSystem}
									onSelectedChange={(v) => {
										v && ($formData.monitoring_control_system = v.value);
									}}
								>
									<Select.Trigger class="w-[500px]">
										<div class="flex items-center justify-start">
											{#if !selectedMonitoringControlSystem}
												Select a monitoring/control system
											{:else}
												{selectedMonitoringControlSystem.label}
											{/if}
										</div>
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											{#each monitoringControlSystems as monitoringControlSystem}
												<Select.Item
													value={monitoringControlSystem.value}
													label={monitoringControlSystem.label}
													>{monitoringControlSystem.label}</Select.Item
												>
											{/each}
										</Select.Group>
									</Select.Content>
									<Select.Input id="monitoringControlSystem" name="monitoringControlSystem" />
								</Select.Root>
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
							<div class="grid gap-3">
								<Label>Starting System</Label>
								<div class="flex items-center space-x-2">
									<Checkbox
										disabled={!isAirStarterCompatible}
										id="startingSystem"
										bind:checked={$formData.starting_system}
										aria-labelledby="startingSystemLabel"
									/>
									<Label
										id="startingSystemLabel"
										for="startingSystem"
										class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
										>Air starter
										{#if !isAirStarterCompatible}
											(only available for 12V and 16V cylinder configurations)
										{/if}
									</Label>
								</div>
								<div class="grid gap-3">
									<Label>Power Transmission</Label>
									<div class="flex items-center space-x-2">
										<Checkbox
											id="powerTransmission"
											bind:checked={$formData.power_transmission}
											aria-labelledby="powerTransmissionLabel"
										/>
										<Label
											id="powerTransmissionLabel"
											for="powerTransmission"
											class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
											>Torsionally resilient coupling</Label
										>
									</div>
								</div>
								<div class="grid gap-3">
									<Label for="auxiliaryPto">Auxiliary PTO</Label>
									<Select.Root portal={null} multiple={true} bind:selected={selectedAuxiliaryPtos}>
										<Select.Trigger class="w-[500px]">
											<div class="flex items-center justify-start">
												{#if selectedAuxiliaryPtos.length === 0}
													Select an auxiliary PTO
												{:else if selectedAuxiliaryPtos.length <= maxAmountOfBadgesInSelect}
													{#each selectedAuxiliaryPtos as item}
														<Badge variant="secondary">{item.label}</Badge>
													{/each}
												{:else}
													{#each { length: maxAmountOfBadgesInSelect } as _, i}
														<Badge variant="secondary">{selectedAuxiliaryPtos[i].label}</Badge>
													{/each}
													<span class="text-xs text-muted-foreground"
														>(+ {selectedAuxiliaryPtos.length - maxAmountOfBadgesInSelect} others)</span
													>
												{/if}
											</div>
										</Select.Trigger>
										<Select.Content>
											<Select.Group>
												{#each auxiliaryPtos as auxiliaryPto}
													<Select.Item value={auxiliaryPto.value} label={auxiliaryPto.label}
														>{auxiliaryPto.label}</Select.Item
													>
												{/each}
											</Select.Group>
										</Select.Content>
										<Select.Input id="auxiliaryPtos" name="auxiliaryPtos" />
									</Select.Root>
								</div>
							</div>
						</div></Card.Content
					>
				</Card.Root>
				<Card.Root>
					<Card.Header>
						<Card.Title>Supply</Card.Title>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6">
							<div class="grid gap-3">
								<Label for="oilSystems">Oil system</Label>
								<Select.Root portal={null} multiple={true} bind:selected={selectedOilSystems}>
									<Select.Trigger class="w-[500px]">
										<div class="flex items-center justify-start">
											{#if selectedOilSystems.length === 0}
												Select an oil system
											{:else if selectedOilSystems.length <= maxAmountOfBadgesInSelect}
												{#each selectedOilSystems as item}
													<Badge variant="secondary">{item.label}</Badge>
												{/each}
											{:else}
												{#each { length: maxAmountOfBadgesInSelect } as _, i}
													<Badge variant="secondary">{selectedOilSystems[i].label}</Badge>
												{/each}
												<span class="text-xs text-muted-foreground"
													>(+ {selectedOilSystems.length - maxAmountOfBadgesInSelect} others)</span
												>
											{/if}
										</div>
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											{#each oilSystems as oilSystem}
												<Select.Item value={oilSystem.value} label={oilSystem.label}
													>{oilSystem.label}</Select.Item
												>
											{/each}
										</Select.Group>
									</Select.Content>
									<Select.Input id="oilSystems" name="oilSystems" />
								</Select.Root>
							</div>
							<div class="grid gap-3">
								<Label for="fuelSystems">Fuel system</Label>
								<Select.Root portal={null} multiple={true} bind:selected={selectedFuelSystems}>
									<Select.Trigger class="w-[500px]">
										<div class="flex items-center justify-start">
											{#if selectedFuelSystems.length === 0}
												Select a fuel system
											{:else if selectedFuelSystems.length <= maxAmountOfBadgesInSelect}
												{#each selectedFuelSystems as item}
													<Badge variant="secondary">{item.label}</Badge>
												{/each}
											{:else}
												{#each { length: maxAmountOfBadgesInSelect } as _, i}
													<Badge variant="secondary">{selectedFuelSystems[i].label}</Badge>
												{/each}
												<span class="text-xs text-muted-foreground"
													>(+ {selectedFuelSystems.length - maxAmountOfBadgesInSelect} others)</span
												>
											{/if}
										</div>
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											{#each fuelSystems as fuelSystem}
												<Select.Item value={fuelSystem.value} label={fuelSystem.label}
													>{fuelSystem.label}</Select.Item
												>
											{/each}
										</Select.Group>
									</Select.Content>
									<Select.Input id="fuelSystems" name="fuelSystems" />
								</Select.Root>
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
							<div class="grid gap-3">
								<Label for="coolingSystem">Cooling system</Label>
								<Select.Root
									portal={null}
									multiple={false}
									bind:selected={selectedCoolingSystem}
									onSelectedChange={(v) => {
										v && ($formData.cooling_system = v.value);
									}}
								>
									<Select.Trigger class="w-[500px]">
										<div class="flex items-center justify-start">
											{#if !selectedCoolingSystem}
												Select a cooling system
											{:else}
												{selectedCoolingSystem.label}
											{/if}
										</div>
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											{#each coolingSystems as coolingSystem}
												<Select.Item value={coolingSystem.value} label={coolingSystem.label}
													>{coolingSystem.label}</Select.Item
												>
											{/each}
										</Select.Group>
									</Select.Content>
									<Select.Input id="coolingSystem" name="coolingSystem" />
								</Select.Root>
							</div>
							<div class="grid gap-3">
								<Label>Exhaust system</Label>
								<div class="flex items-center space-x-2">
									<Checkbox
										id="exhaustSystem"
										bind:checked={$formData.exhaust_system}
										aria-labelledby="exhaustSystemLabel"
									/>
									<Label
										id="exhaustSystemLabel"
										for="exhaustSystem"
										class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
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
							<div class="grid gap-3">
								<Label>Mounting system</Label>
								<div class="flex items-center space-x-2">
									<Checkbox
										id="moungingSystem"
										bind:checked={$formData.mounting_system}
										aria-labelledby="moungingSystemLabel"
									/>
									<Label
										id="mountingSystemLabel"
										for="mountingSystem"
										class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
										>Resilient mounts at driving end</Label
									>
								</div>
							</div>
							<div class="grid gap-3">
								<Label for="gearboxOptions">Gearbox options</Label>
								<Select.Root portal={null} multiple={true} bind:selected={selectedGearboxOptions}>
									<Select.Trigger class="w-[500px]">
										<div class="flex items-center justify-start">
											{#if selectedGearboxOptions.length === 0}
												Select a gearbox option
											{:else if selectedGearboxOptions.length <= maxAmountOfBadgesInSelect}
												{#each selectedGearboxOptions as item}
													<Badge variant="secondary">{item.label}</Badge>
												{/each}
											{:else}
												{#each { length: maxAmountOfBadgesInSelect } as _, i}
													<Badge variant="secondary">{selectedGearboxOptions[i].label}</Badge>
												{/each}
												<span class="text-xs text-muted-foreground"
													>(+ {selectedGearboxOptions.length - maxAmountOfBadgesInSelect} others)</span
												>
											{/if}
										</div>
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											{#each gearboxOptions as gearboxOption}
												<Select.Item value={gearboxOption.value} label={gearboxOption.label}
													>{gearboxOption.label}</Select.Item
												>
											{/each}
										</Select.Group>
									</Select.Content>
									<Select.Input id="gearboxOptions" name="gearboxOptions" />
								</Select.Root>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</div>
		<div class="flex items-center justify-center gap-2 md:justify-end">
			<Button type="reset" variant="outline" on:click={resetForm}>Discard</Button>
			<Button disabled={$formIsValid || awaitingResponse} on:click={submit}>
				{#if !awaitingResponse}
					Submit
				{:else}
					<LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
					Please wait
				{/if}
			</Button>
		</div>
	</div>
</form>
