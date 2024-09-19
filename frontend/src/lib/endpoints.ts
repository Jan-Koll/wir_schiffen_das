import { env } from "$env/dynamic/private"
const configuration_manager_host = env.CONFIGURATION_MANAGER_PORT ? "wir_schiffen_das-configuration_manager-1" : "127.0.0.1" // containers hosts are accessible by their name
const configuration_manager_port = env.CONFIGURATION_MANAGER_PORT ?? "8081"
export const configuration_manager_url = `http://${configuration_manager_host}:${configuration_manager_port}`

const cooling_exhaust_analyzer_host = env.COOLING_EXHAUST_ANALYZER_PORT ? "wir_schiffen_das-cooling_exhaust_analyzer-1" : "127.0.0.1" // containers hosts are accessible by their name
const cooling_exhaust_analyzer_port = env.COOLING_EXHAUST_ANALYZER_PORT ?? "8082"
export const cooling_exhaust_analyzer_url = `http://${cooling_exhaust_analyzer_host}:${cooling_exhaust_analyzer_port}`

const engine_control_analyzer_host = env.ENGINE_CONTROL_ANALYZER_PORT ? "wir_schiffen_das-engine_control_analyzer-1" : "127.0.0.1" // containers hosts are accessible by their name
const engine_control_analyzer_port = env.ENGINE_CONTROL_ANALYZER_PORT ?? "8083"
export const engine_control_analyzer_url = `http://${engine_control_analyzer_host}:${engine_control_analyzer_port}`

const mounting_analyzer_host = env.MOUNTING_ANALYZER_PORT ? "wir_schiffen_das-mounting_analyzer-1" : "127.0.0.1" // containers hosts are accessible by their name
const mounting_analyzer_port = env.MOUNTING_ANALYZER_PORT ?? "8084"
export const mounting_analyzer_url = `http://${mounting_analyzer_host}:${mounting_analyzer_port}`

const propulsion_analyzer_host = env.PROPULSION_ANALYZER_PORT ? "wir_schiffen_das-propulsion_analyzer-1" : "127.0.0.1" // containers hosts are accessible by their name
const propulsion_analyzer_port = env.PROPULSION_ANALYZER_PORT ?? "8085"
export const propulsion_analyzer_url = `http://${propulsion_analyzer_host}:${propulsion_analyzer_port}`

const supply_analyzer_host = env.SUPPLY_ANALYZER_PORT ? "wir_schiffen_das-supply_analyzer-1" : "127.0.0.1" // containers hosts are accessible by their name
const supply_analyzer_port = env.SUPPLY_ANALYZER_PORT ?? "8086"
export const supply_analyzer_url = `http://${supply_analyzer_host}:${supply_analyzer_port}`