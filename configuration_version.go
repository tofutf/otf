package ots

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	tfe "github.com/leg100/go-tfe"
	"gorm.io/gorm"
)

const (
	DefaultAutoQueueRuns       = true
	DefaultConfigurationSource = "tfe-api"
)

var (
	ErrInvalidConfigurationVersionGetOptions = errors.New("invalid configuration version get options")
)

// ConfigurationVersionList represents a list of configuration versions.
type ConfigurationVersionList struct {
	*tfe.Pagination
	Items []*ConfigurationVersion
}

// ConfigurationVersion is a representation of an uploaded or ingressed
// Terraform configuration in TFE. A workspace must have at least one
// configuration version before any runs may be queued on it.
type ConfigurationVersion struct {
	ExternalID string `gorm:"uniqueIndex"`
	InternalID uint   `gorm:"primaryKey;column:id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	AutoQueueRuns    bool
	Error            string
	ErrorMessage     string
	Source           tfe.ConfigurationSource
	Speculative      bool
	Status           tfe.ConfigurationStatus
	StatusTimestamps *tfe.CVStatusTimestamps `gorm:"embedded;embeddedPrefix:timestamp_"`

	Configuration []byte

	// Configuration Version belongs to a Workspace
	WorkspaceID uint
	Workspace   *Workspace
}

type ConfigurationVersionService interface {
	Create(workspaceID string, opts *tfe.ConfigurationVersionCreateOptions) (*ConfigurationVersion, error)
	Get(id string) (*ConfigurationVersion, error)
	GetLatest(workspaceID string) (*ConfigurationVersion, error)
	List(workspaceID string, opts tfe.ConfigurationVersionListOptions) (*ConfigurationVersionList, error)
	Upload(id string, payload []byte) error
	Download(id string) ([]byte, error)
}

type ConfigurationVersionRepository interface {
	Create(run *ConfigurationVersion) (*ConfigurationVersion, error)
	Get(opts ConfigurationVersionGetOptions) (*ConfigurationVersion, error)
	List(workspaceID string, opts ConfigurationVersionListOptions) (*ConfigurationVersionList, error)
	Update(id string, fn func(*ConfigurationVersion) error) (*ConfigurationVersion, error)
}

// ConfigurationVersionGetOptions are options for retrieving a single config
// version. Either ID *or* WorkspaceID must be specfiied.
type ConfigurationVersionGetOptions struct {
	// ID of config version to retrieve
	ID *string

	// Get latest config version for this workspace ID
	WorkspaceID *string
}

// ConfigurationVersionListOptions are options for paginating and filtering the list of runs to
// retrieve from the ConfigurationVersionRepository ListConfigurationVersions endpoint
type ConfigurationVersionListOptions struct {
	tfe.ListOptions

	// Filter by run statuses (with an implicit OR condition)
	Statuses []tfe.ConfigurationStatus

	// Filter by workspace ID
	WorkspaceID *string
}

type ConfigurationVersionFactory struct {
	WorkspaceService WorkspaceService
}

func (f *ConfigurationVersionFactory) NewConfigurationVersion(workspaceID string, opts *tfe.ConfigurationVersionCreateOptions) (*ConfigurationVersion, error) {
	cv := ConfigurationVersion{
		ExternalID:    NewConfigurationVersionID(),
		AutoQueueRuns: DefaultAutoQueueRuns,
		Status:        tfe.ConfigurationPending,
		Source:        DefaultConfigurationSource,
	}

	if opts.AutoQueueRuns != nil {
		cv.AutoQueueRuns = *opts.AutoQueueRuns
	}

	if opts.Speculative != nil {
		cv.Speculative = *opts.Speculative
	}

	ws, err := f.WorkspaceService.GetByID(workspaceID)
	if err != nil {
		return nil, err
	}
	cv.Workspace = ws
	cv.WorkspaceID = ws.InternalID

	return &cv, nil
}

func NewConfigurationVersionID() string {
	return fmt.Sprintf("cv-%s", GenerateRandomString(16))
}

func (cv *ConfigurationVersion) UploadURL() string {
	return fmt.Sprintf("/configuration-versions/%s/upload", cv.ExternalID)
}

func (cv *ConfigurationVersion) DTO() interface{} {
	dto := &tfe.ConfigurationVersion{
		ID:            cv.ExternalID,
		AutoQueueRuns: cv.AutoQueueRuns,
		Error:         cv.Error,
		ErrorMessage:  cv.ErrorMessage,
		Speculative:   cv.Speculative,
		Source:        cv.Source,
		Status:        cv.Status,
		UploadURL:     cv.UploadURL(),
	}

	if cv.StatusTimestamps != nil && !reflect.ValueOf(cv.StatusTimestamps).Elem().IsZero() {
		dto.StatusTimestamps = cv.StatusTimestamps
	}

	return dto
}

func (cvl *ConfigurationVersionList) DTO() interface{} {
	l := &tfe.ConfigurationVersionList{
		Pagination: cvl.Pagination,
	}
	for _, item := range cvl.Items {
		l.Items = append(l.Items, item.DTO().(*tfe.ConfigurationVersion))
	}

	return l
}