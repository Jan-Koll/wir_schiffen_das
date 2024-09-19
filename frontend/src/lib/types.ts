export type configuration = {
    id?: number;
    description: string;
    engine: string;
    gearbox_type: string;
    monitoring_control_system: string;
    cooling_system: string;
    auxiliary_pto: string[];
    oil_system: string[];
    fuel_system: string[];
    gearbox_option: string[];
    engine_management_system: boolean;
    starting_system: boolean;
    power_transmission: boolean;
    exhaust_system: boolean;
    mounting_system: boolean;
    created_at?: string;
    last_modified?:string;
};

export type cooling_exhaust_analyzer_response = {
    id: number,
    configuration_id: number,
    success: boolean, 
    status: string,
    cooling_system: string,
    exhaust_system: boolean,
    created_at: string|Date,
    last_modified: string|Date,
    order_created_at: string|Date,
}
export type engine_control_analyzer_response = {
    id: number,
    configuration_id: number,
    success: boolean, 
    status: string,
    engine_management: boolean,
    monitoring_control_system: string,
    created_at: string|Date,
    last_modified: string|Date,
    order_created_at: string|Date,
}
export type mounting_analyzer_response = {
    id: number,
    configuration_id: number,
    success: boolean, 
    status: string,
    mounting_system: boolean,
    gearbox_options: string[],
    created_at: string|Date,
    last_modified: string|Date,
    order_created_at: string|Date,
}
export type propulsion_analyzer_response = {
    id: number,
    configuration_id: number,
    success: boolean, 
    status: string,
    starting_system: boolean,
    power_transmission: boolean,
    created_at: string|Date,
    last_modified: string|Date,
    order_created_at: string|Date,
}
export type supply_analyzer_response = {
    id: number,
    configuration_id: number,
    success: boolean, 
    status: string,
    oil_system: string[],
    fuel_system: string[],
    created_at: string|Date,
    last_modified: string|Date,
    order_created_at: string|Date,
}