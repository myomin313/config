package requests

type MaintenanceUpdateRequest struct {
    ApplicationMaintenance *bool `json:"application_maintenance"`
    GameMaintenance        *bool `json:"game_maintenance"`
}
