package streams

type Change struct {
	Project         string       `json:"project,omitempty"`
	Branch          string       `json:"branch,omitempty"`
	Topic           string       `json:"topic,omitempty"`
	ID              string       `json:"id,omitempty"`
	Number          int          `json:"number,omitempty"`
	Subject         string       `json:"subject,omitempty"`
	Owner           Account      `json:"owner"`
	URL             string       `json:"url,omitempty"`
	CommitMessage   string       `json:"commitMessage,omitempty"`
	Hashtags        []string     `json:"hashtags,omitempty"`
	CreatedOn       int          `json:"createdOn,omitempty"`
	LastUpdated     int          `json:"lastUpdated,omitempty"`
	Open            bool         `json:"open,omitempty"`
	Status          string       `json:"status,omitempty"`
	Private         bool         `json:"private,omitempty"`
	Wip             bool         `json:"wip,omitempty"`
	Comments        []Message    `json:"comments,omitempty"`
	TrackingIDs     []TrackingID `json:"trackingIds,omitempty"`
	CurrentPatchSet PatchSet     `json:"currentPatchSet"`
	PatchSets       []PatchSet   `json:"patchSets,omitempty"`
	DependsOn       []Dependency `json:"dependsOn,omitempty"`
	NeededBy        []Dependency `json:"neededBy,omitempty"`
	SubmitRecords   SubmitRecord `json:"submitRecords"`
	AllReviewers    []Account    `json:"allReviewers,omitempty"`
}

type TrackingID struct {
	System string `json:"system,omitempty"`
	ID     string `json:"id,omitempty"`
}

type Account struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

type PatchSet struct {
	Number         int               `json:"number,omitempty"`
	Revision       string            `json:"revision,omitempty"`
	Parents        []string          `json:"parents,omitempty"`
	Ref            string            `json:"ref,omitempty"`
	Uploader       Account           `json:"uploader"`
	Author         Account           `json:"author"`
	CreatedOn      int               `json:"createdOn,omitempty"`
	Kind           string            `json:"kind,omitempty"`
	Approvals      Approval          `json:"approvals"`
	Comments       []PatchsetComment `json:"comments,omitempty"`
	Files          []File            `json:"files,omitempty"`
	SizeInsertions int               `json:"sizeInsertions,omitempty"`
	SizeDeletions  int               `json:"sizeDeletions,omitempty"`
}

type Approval struct {
	Type        string  `json:"type,omitempty"`
	Description string  `json:"description,omitempty"`
	Value       string  `json:"value,omitempty"`
	OldValue    string  `json:"oldValue,omitempty"`
	GrantedOn   int     `json:"grantedOn,omitempty"`
	By          Account `json:"by"`
}

type RefUpdate struct {
	OldRev  string `json:"oldRev,omitempty"`
	NewRev  string `json:"newRev,omitempty"`
	RefName string `json:"refName,omitempty"`
	Project string `json:"project,omitempty"`
}

type SubmitRecord struct {
	Status       string        `json:"status,omitempty"`
	Labels       []Label       `json:"labels,omitempty"`
	Requirements []Requirement `json:"requirements,omitempty"`
}

type Requirement struct {
	FallbackText string      `json:"fallbackText,omitempty"`
	Type         string      `json:"type,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type Label struct {
	Label  string  `json:"label,omitempty"`
	Status string  `json:"status,omitempty"`
	By     Account `json:"by"`
}

type Dependency struct {
	ID                string `json:"id,omitempty"`
	Number            int    `json:"number,omitempty"`
	Revision          string `json:"revision,omitempty"`
	Ref               string `json:"ref,omitempty"`
	IsCurrentPatchSet bool   `json:"isCurrentPatchSet,omitempty"`
}

type Message struct {
	Timestamp int     `json:"timestamp,omitempty"`
	Reviewer  Account `json:"reviewer"`
	Message   string  `json:"message,omitempty"`
}

type PatchsetComment struct {
	File     string  `json:"file,omitempty"`
	Line     string  `json:"line,omitempty"`
	Reviewer Account `json:"reviewer"`
	Message  string  `json:"message,omitempty"`
}

type File struct {
	File       string `json:"file,omitempty"`
	FileOld    string `json:"fileOld,omitempty"`
	Type       string `json:"type,omitempty"`
	Insertions int    `json:"insertions,omitempty"`
	Deletions  int    `json:"deletions,omitempty"`
}
