package Resources

type GameHallResource struct {
	ID           uint64 `json:"id"`
	FormattedID  string `json:"formatted_id"` 
	Status       string `json:"status"`
	SortOrder    int    `json:"sort_order"`
	CreatedBy    string `json:"created_by"`
	UpdatedBy    string `json:"updated_by"`
	EnName       string `json:"en_name"`
	ZhName       string `json:"zh_name"`
	KmName       string `json:"km_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
 type GameHallListResource struct {
	ID           uint64 `json:"id"`
	//Name         string `json:"string"`
	EnName       string `json:"en_name"`
	ZhName       string `json:"zh_name"`
	KmName       string `json:"km_name"`
 }
 type GameHallListFEResource struct {
	ID           uint64 `json:"id"`
	Name         string `json:"string"`
	SortOrder    uint `json:"sort_order"`
	Status       string `json:"status"`
	//EnName       string `json:"en_name"`
	//ZhName       string `json:"zh_name"`
	//KmName       string `json:"km_name"`
 }