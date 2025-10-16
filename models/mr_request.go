package models

type MRRequest struct {
	AuthToken           string `json:"auth_token" validate:"required"`
	FilePath            string
	ArtifactFormat      string `json:"format" validate:"required"`
	ArtifactFileName    string `json:"artifact_file_name"`
	ProjectId           int    `json:"project_id"`
	JobId               int    `json:"job_id"`
	MergeRequestIid     int    `json:"merge_request_iid"`
	VulnerabilityMgmtId int    `json:"vulnerability_mgmt_id"`
}
