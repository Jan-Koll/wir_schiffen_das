package main

import (
	"log"

	"github.com/lib/pq"
)

func getConfigurationById(id int) (result configuration, err error) {
	row := db.QueryRow("SELECT * FROM configurations WHERE configuration_id = $1", id)
	err = row.Scan(&result.ConfigurationId, &result.Description, &result.Engine, &result.GearboxType, &result.EngineManagementSystem, &result.MonitoringControlSystem, &result.StartingSystem, &result.PowerTransmission, pq.Array(&result.AuxiliaryPto), pq.Array(&result.OilSystem), pq.Array(&result.FuelSystem), &result.CoolingSystem, &result.ExhaustSystem, &result.MountingSystem, pq.Array(&result.GearboxOption), &result.CreatedAt, &result.LastModified)
	return
}

func createConfiguration(input configuration) (result configuration, err error) {
	result = input
	row := db.QueryRow(
		`INSERT INTO configurations (description, engine, gearbox_type, engine_management_system, monitoring_control_system, starting_system, power_transmission, auxiliary_pto, oil_system, fuel_system, cooling_system, exhaust_system, mounting_system, gearbox_option )
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		 RETURNING configuration_id, created_at, last_modified`,
		input.Description,
		input.Engine,
		input.GearboxType,
		input.EngineManagementSystem,
		input.MonitoringControlSystem,
		input.StartingSystem,
		input.PowerTransmission,
		pq.Array(input.AuxiliaryPto),
		pq.Array(input.OilSystem),
		pq.Array(input.FuelSystem),
		input.CoolingSystem,
		input.ExhaustSystem,
		input.MountingSystem,
		pq.Array(input.GearboxOption),
	)
	err = row.Scan(&result.ConfigurationId, &result.CreatedAt, &result.LastModified)
	if err != nil {
		result = configuration{}
	}
	return
}
func updateConfiguration(input configuration) (result configuration, err error) {
	result = input
	row := db.QueryRow(
		`UPDATE configurations SET 
				description = $2,
				engine = $3,
				gearbox_type = $4,
				engine_management_system = $5,
				monitoring_control_system = $6,
				air_starter = $7,
				power_transmission = $8,
				auxiliary_pto = $9, 
				oil_system = $10, 
				fuel_system = $11, 
				cooling_system = $12,
				exhaust_system = $13, 
				mounting_system = $14, 
				gearbox_option = $15,
				last_modified = now()
			WHERE configuration_id = $1
			RETURNING created_at, last_modified;`,
		input.ConfigurationId,
		input.Description,
		input.Engine,
		input.GearboxType,
		input.EngineManagementSystem,
		input.MonitoringControlSystem,
		input.StartingSystem,
		input.PowerTransmission,
		pq.Array(input.AuxiliaryPto),
		pq.Array(input.OilSystem),
		pq.Array(input.FuelSystem),
		input.CoolingSystem,
		input.ExhaustSystem,
		input.MountingSystem,
		pq.Array(input.GearboxOption),
	)
	err = row.Scan(&result.CreatedAt, &result.LastModified)
	if err != nil {
		result = configuration{}
	}
	return

}
func createOrUpdate(input configuration) (result configuration, err error) {
	create_new_configuration := input.ConfigurationId == 0 // 0 means that no id was set, so a new configuration will be added

	// In case the user specifies an id (indicating to update an resource), we check if the resource exists in the database.
	// In that way, we can still create the new resource if it doesn't exist yet.
	if !create_new_configuration {
		_, err := getConfigurationById(input.ConfigurationId)
		if err != nil {
			create_new_configuration = false
		}
	}

	if create_new_configuration {
		result, err = createConfiguration(input)
	} else {
		result, err = updateConfiguration(input)
	}
	return
}
func deleteById(id int) (err error) {
	_, err = db.Exec("DELETE FROM configurations WHERE configuration_id = $1", id)
	return
}

func getAllConfigurations() (result []configuration, err error) {
	rows, err := db.Query("SELECT * FROM configurations")
	if err != nil {
		log.Println("Failed to query database\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var configuration configuration
		err = rows.Scan(&configuration.ConfigurationId, &configuration.Description, &configuration.Engine, &configuration.GearboxType, &configuration.EngineManagementSystem, &configuration.MonitoringControlSystem, &configuration.StartingSystem, &configuration.PowerTransmission, pq.Array(&configuration.AuxiliaryPto), pq.Array(&configuration.OilSystem), pq.Array(&configuration.FuelSystem), &configuration.CoolingSystem, &configuration.ExhaustSystem, &configuration.MountingSystem, pq.Array(&configuration.GearboxOption), &configuration.CreatedAt, &configuration.LastModified)
		if err != nil {
			log.Println("Failed to read row\n", err)
			return
		}
		result = append(result, configuration)
	}
	return
}
