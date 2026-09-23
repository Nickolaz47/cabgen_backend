package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AuditEventLogin                               = "auth.login"
	AuditEventLoginFailed                         = "auth.login_failed"
	AuditEventRegister                            = "auth.register"
	AuditEventRegisterFailed                      = "auth.register_failed"
	AuditEventRefreshToken                        = "auth.refresh"
	AuditEventRefreshTokenFailed                  = "auth.refresh_failed"
	AuditEventForgotPassword                      = "auth.forgotPassword"
	AuditEventForgotPasswordFailed                = "auth.forgotPassword_failed"
	AuditEventResetPassword                       = "auth.resetPassword"
	AuditEventResetPasswordFailed                 = "auth.resetPassword_failed"
	AuditEventContact                             = "auth.contact"
	AuditEventContactFailed                       = "auth.contact_failed"
	AuditEventAnalysesGet                         = "analyses.get"
	AuditEventAnalysesGetFailed                   = "analyses.get_failed"
	AuditEventAnalysesGetByID                     = "analyses.getById"
	AuditEventAnalysesGetByIDFailed               = "analyses.getById_failed"
	AuditEventAnalysesGetFastQCReport             = "analyses.getFastqcReport"
	AuditEventAnalysesGetFastQCReportFailed       = "analyses.getFastqcReport_failed"
	AuditEventAnalysesDownloadZip                 = "analyses.downloadZip"
	AuditEventAnalysesDownloadZipFailed           = "analyses.downloadZip_failed"
	AuditEventAnalysesCreate                      = "analyses.create"
	AuditEventAnalysesCreateFailed                = "analyses.create_failed"
	AuditEventAnalysesDownloadBatchTSV            = "analyses.downloadBatchTsv"
	AuditEventAnalysesDownloadBatchTSVFailed      = "analyses.downloadBatchTsv_failed"
	AuditEventAnalysesDelete                      = "analyses.delete"
	AuditEventAnalysesDeleteFailed                = "analyses.delete_failed"
	AuditEventLogout                              = "auth.logout"
	AuditEventLogoutFailed                        = "auth.logout_failed"
	AuditEventMe                                  = "auth.me"
	AuditEventMeFailed                            = "auth.me_failed"
	AuditEventSamplesGet                          = "samples.get"
	AuditEventSamplesGetFailed                    = "samples.get_failed"
	AuditEventSamplesGetByID                      = "samples.getById"
	AuditEventSamplesGetByIDFailed                = "samples.getById_failed"
	AuditEventSamplesCreate                       = "samples.create"
	AuditEventSamplesCreateFailed                 = "samples.create_failed"
	AuditEventSamplesUpload                       = "samples.upload"
	AuditEventSamplesUploadFailed                 = "samples.upload_failed"
	AuditEventSamplesUpdate                       = "samples.update"
	AuditEventSamplesUpdateFailed                 = "samples.update_failed"
	AuditEventSamplesDelete                       = "samples.delete"
	AuditEventSamplesDeleteFailed                 = "samples.delete_failed"
	AuditEventUsersGet                            = "users.get"
	AuditEventUsersGetFailed                      = "users.get_failed"
	AuditEventUsersUpdate                         = "users.update"
	AuditEventUsersUpdateFailed                   = "users.update_failed"
	AuditEventUsersDelete                         = "users.delete"
	AuditEventUsersDeleteFailed                   = "users.delete_failed"
	AuditEventUsersUpdatePassword                 = "users.updatePassword"
	AuditEventUsersUpdatePasswordFailed           = "users.updatePassword_failed"
	AuditEventUsersRequestEmailUpdate             = "users.requestEmailUpdate"
	AuditEventUsersRequestEmailUpdateFailed       = "users.requestEmailUpdate_failed"
	AuditEventUsersConfirmEmailUpdate             = "users.confirmEmailUpdate"
	AuditEventUsersConfirmEmailUpdateFailed       = "users.confirmEmailUpdate_failed"
	AuditEventAdminAnalysesGet                    = "admin.analyses.get"
	AuditEventAdminAnalysesGetFailed              = "admin.analyses.get_failed"
	AuditEventAdminAnalysesGetByID                = "admin.analyses.getById"
	AuditEventAdminAnalysesGetByIDFailed          = "admin.analyses.getById_failed"
	AuditEventAdminAnalysesDownloadZip            = "admin.analyses.downloadZip"
	AuditEventAdminAnalysesDownloadZipFailed      = "admin.analyses.downloadZip_failed"
	AuditEventAdminAnalysesCreate                 = "admin.analyses.create"
	AuditEventAdminAnalysesCreateFailed           = "admin.analyses.create_failed"
	AuditEventAdminAnalysesDownloadBatchTSV       = "admin.analyses.downloadBatchTsv"
	AuditEventAdminAnalysesDownloadBatchTSVFailed = "admin.analyses.downloadBatchTsv_failed"
	AuditEventAdminAnalysesUpdate                 = "admin.analyses.update"
	AuditEventAdminAnalysesUpdateFailed           = "admin.analyses.update_failed"
	AuditEventAdminAnalysesDelete                 = "admin.analyses.delete"
	AuditEventAdminAnalysesDeleteFailed           = "admin.analyses.delete_failed"
	AuditEventAdminHealthServicesGet              = "admin.healthServices.get"
	AuditEventAdminHealthServicesGetFailed        = "admin.healthServices.get_failed"
	AuditEventAdminHealthServicesGetByID          = "admin.healthServices.getById"
	AuditEventAdminHealthServicesGetByIDFailed    = "admin.healthServices.getById_failed"
	AuditEventAdminHealthServicesSearch           = "admin.healthServices.search"
	AuditEventAdminHealthServicesSearchFailed     = "admin.healthServices.search_failed"
	AuditEventAdminHealthServicesCreate           = "admin.healthServices.create"
	AuditEventAdminHealthServicesCreateFailed     = "admin.healthServices.create_failed"
	AuditEventAdminHealthServicesUpdate           = "admin.healthServices.update"
	AuditEventAdminHealthServicesUpdateFailed     = "admin.healthServices.update_failed"
	AuditEventAdminHealthServicesDelete           = "admin.healthServices.delete"
	AuditEventAdminHealthServicesDeleteFailed     = "admin.healthServices.delete_failed"
	AuditEventAdminCountriesGet                   = "admin.countries.get"
	AuditEventAdminCountriesGetFailed             = "admin.countries.get_failed"
	AuditEventAdminCountriesGetByID               = "admin.countries.getById"
	AuditEventAdminCountriesGetByIDFailed         = "admin.countries.getById_failed"
	AuditEventAdminCountriesSearch                = "admin.countries.search"
	AuditEventAdminCountriesSearchFailed          = "admin.countries.search_failed"
	AuditEventAdminCountriesCreate                = "admin.countries.create"
	AuditEventAdminCountriesCreateFailed          = "admin.countries.create_failed"
	AuditEventAdminCountriesUpdate                = "admin.countries.update"
	AuditEventAdminCountriesUpdateFailed          = "admin.countries.update_failed"
	AuditEventAdminCountriesDelete                = "admin.countries.delete"
	AuditEventAdminCountriesDeleteFailed          = "admin.countries.delete_failed"
	AuditEventAdminLaboratoriesGet                = "admin.laboratories.get"
	AuditEventAdminLaboratoriesGetFailed          = "admin.laboratories.get_failed"
	AuditEventAdminLaboratoriesGetByID            = "admin.laboratories.getById"
	AuditEventAdminLaboratoriesGetByIDFailed      = "admin.laboratories.getById_failed"
	AuditEventAdminLaboratoriesSearch             = "admin.laboratories.search"
	AuditEventAdminLaboratoriesSearchFailed       = "admin.laboratories.search_failed"
	AuditEventAdminLaboratoriesCreate             = "admin.laboratories.create"
	AuditEventAdminLaboratoriesCreateFailed       = "admin.laboratories.create_failed"
	AuditEventAdminLaboratoriesUpdate             = "admin.laboratories.update"
	AuditEventAdminLaboratoriesUpdateFailed       = "admin.laboratories.update_failed"
	AuditEventAdminLaboratoriesDelete             = "admin.laboratories.delete"
	AuditEventAdminLaboratoriesDeleteFailed       = "admin.laboratories.delete_failed"
	AuditEventAdminMetrics                        = "admin.metrics"
	AuditEventAdminMetricsFailed                  = "admin.metrics_failed"
	AuditEventAdminMicroorganismsGet              = "admin.microorganisms.get"
	AuditEventAdminMicroorganismsGetFailed        = "admin.microorganisms.get_failed"
	AuditEventAdminMicroorganismsGetByID          = "admin.microorganisms.getById"
	AuditEventAdminMicroorganismsGetByIDFailed    = "admin.microorganisms.getById_failed"
	AuditEventAdminMicroorganismsSearch           = "admin.microorganisms.search"
	AuditEventAdminMicroorganismsSearchFailed     = "admin.microorganisms.search_failed"
	AuditEventAdminMicroorganismsCreate           = "admin.microorganisms.create"
	AuditEventAdminMicroorganismsCreateFailed     = "admin.microorganisms.create_failed"
	AuditEventAdminMicroorganismsUpdate           = "admin.microorganisms.update"
	AuditEventAdminMicroorganismsUpdateFailed     = "admin.microorganisms.update_failed"
	AuditEventAdminMicroorganismsDelete           = "admin.microorganisms.delete"
	AuditEventAdminMicroorganismsDeleteFailed     = "admin.microorganisms.delete_failed"
	AuditEventAdminOriginsGet                     = "admin.origins.get"
	AuditEventAdminOriginsGetFailed               = "admin.origins.get_failed"
	AuditEventAdminOriginsGetByID                 = "admin.origins.getById"
	AuditEventAdminOriginsGetByIDFailed           = "admin.origins.getById_failed"
	AuditEventAdminOriginsSearch                  = "admin.origins.search"
	AuditEventAdminOriginsSearchFailed            = "admin.origins.search_failed"
	AuditEventAdminOriginsCreate                  = "admin.origins.create"
	AuditEventAdminOriginsCreateFailed            = "admin.origins.create_failed"
	AuditEventAdminOriginsUpdate                  = "admin.origins.update"
	AuditEventAdminOriginsUpdateFailed            = "admin.origins.update_failed"
	AuditEventAdminOriginsDelete                  = "admin.origins.delete"
	AuditEventAdminOriginsDeleteFailed            = "admin.origins.delete_failed"
	AuditEventAdminSamplesGet                     = "admin.samples.get"
	AuditEventAdminSamplesGetFailed               = "admin.samples.get_failed"
	AuditEventAdminSamplesGetByID                 = "admin.samples.getById"
	AuditEventAdminSamplesGetByIDFailed           = "admin.samples.getById_failed"
	AuditEventAdminSamplesCreate                  = "admin.samples.create"
	AuditEventAdminSamplesCreateFailed            = "admin.samples.create_failed"
	AuditEventAdminSamplesUpload                  = "admin.samples.upload"
	AuditEventAdminSamplesUploadFailed            = "admin.samples.upload_failed"
	AuditEventAdminSamplesUpdate                  = "admin.samples.update"
	AuditEventAdminSamplesUpdateFailed            = "admin.samples.update_failed"
	AuditEventAdminSamplesDelete                  = "admin.samples.delete"
	AuditEventAdminSamplesDeleteFailed            = "admin.samples.delete_failed"
	AuditEventAdminSampleSourcesGet               = "admin.sampleSources.get"
	AuditEventAdminSampleSourcesGetFailed         = "admin.sampleSources.get_failed"
	AuditEventAdminSampleSourcesGetByID           = "admin.sampleSources.getById"
	AuditEventAdminSampleSourcesGetByIDFailed     = "admin.sampleSources.getById_failed"
	AuditEventAdminSampleSourcesSearch            = "admin.sampleSources.search"
	AuditEventAdminSampleSourcesSearchFailed      = "admin.sampleSources.search_failed"
	AuditEventAdminSampleSourcesCreate            = "admin.sampleSources.create"
	AuditEventAdminSampleSourcesCreateFailed      = "admin.sampleSources.create_failed"
	AuditEventAdminSampleSourcesUpdate            = "admin.sampleSources.update"
	AuditEventAdminSampleSourcesUpdateFailed      = "admin.sampleSources.update_failed"
	AuditEventAdminSampleSourcesDelete            = "admin.sampleSources.delete"
	AuditEventAdminSampleSourcesDeleteFailed      = "admin.sampleSources.delete_failed"
	AuditEventAdminSequencersGet                  = "admin.sequencers.get"
	AuditEventAdminSequencersGetFailed            = "admin.sequencers.get_failed"
	AuditEventAdminSequencersGetByID              = "admin.sequencers.getById"
	AuditEventAdminSequencersGetByIDFailed        = "admin.sequencers.getById_failed"
	AuditEventAdminSequencersSearch               = "admin.sequencers.search"
	AuditEventAdminSequencersSearchFailed         = "admin.sequencers.search_failed"
	AuditEventAdminSequencersCreate               = "admin.sequencers.create"
	AuditEventAdminSequencersCreateFailed         = "admin.sequencers.create_failed"
	AuditEventAdminSequencersUpdate               = "admin.sequencers.update"
	AuditEventAdminSequencersUpdateFailed         = "admin.sequencers.update_failed"
	AuditEventAdminSequencersDelete               = "admin.sequencers.delete"
	AuditEventAdminSequencersDeleteFailed         = "admin.sequencers.delete_failed"
	AuditEventAdminTicketsGet                     = "admin.tickets.get"
	AuditEventAdminTicketsGetFailed               = "admin.tickets.get_failed"
	AuditEventAdminTicketsGetByID                 = "admin.tickets.getById"
	AuditEventAdminTicketsGetByIDFailed           = "admin.tickets.getById_failed"
	AuditEventAdminTicketsAssign                  = "admin.tickets.assign"
	AuditEventAdminTicketsAssignFailed            = "admin.tickets.assign_failed"
	AuditEventAdminTicketsResolve                 = "admin.tickets.resolve"
	AuditEventAdminTicketsResolveFailed           = "admin.tickets.resolve_failed"
	AuditEventAdminTicketsDelete                  = "admin.tickets.delete"
	AuditEventAdminTicketsDeleteFailed            = "admin.tickets.delete_failed"
	AuditEventAdminUsersGet                       = "admin.users.get"
	AuditEventAdminUsersGetFailed                 = "admin.users.get_failed"
	AuditEventAdminUsersGetByID                   = "admin.users.getById"
	AuditEventAdminUsersGetByIDFailed             = "admin.users.getById_failed"
	AuditEventAdminUsersCreate                    = "admin.users.create"
	AuditEventAdminUsersCreateFailed              = "admin.users.create_failed"
	AuditEventAdminUsersUpdate                    = "admin.users.update"
	AuditEventAdminUsersUpdateFailed              = "admin.users.update_failed"
	AuditEventAdminUsersActivate                  = "admin.users.activate"
	AuditEventAdminUsersActivateFailed            = "admin.users.activate_failed"
	AuditEventAdminUsersDeactivate                = "admin.users.deactivate"
	AuditEventAdminUsersDeactivateFailed          = "admin.users.deactivate_failed"
	AuditEventAdminUsersDelete                    = "admin.users.delete"
	AuditEventAdminUsersDeleteFailed              = "admin.users.delete_failed"
)

var AuditEvents = []string{
	AuditEventLogin,
	AuditEventLoginFailed,
	AuditEventRegister,
	AuditEventRegisterFailed,
	AuditEventRefreshToken,
	AuditEventRefreshTokenFailed,
	AuditEventForgotPassword,
	AuditEventForgotPasswordFailed,
	AuditEventResetPassword,
	AuditEventResetPasswordFailed,
	AuditEventContact,
	AuditEventContactFailed,
	AuditEventAnalysesGet,
	AuditEventAnalysesGetFailed,
	AuditEventAnalysesGetByID,
	AuditEventAnalysesGetByIDFailed,
	AuditEventAnalysesGetFastQCReport,
	AuditEventAnalysesGetFastQCReportFailed,
	AuditEventAnalysesDownloadZip,
	AuditEventAnalysesDownloadZipFailed,
	AuditEventAnalysesCreate,
	AuditEventAnalysesCreateFailed,
	AuditEventAnalysesDownloadBatchTSV,
	AuditEventAnalysesDownloadBatchTSVFailed,
	AuditEventAnalysesDelete,
	AuditEventAnalysesDeleteFailed,
	AuditEventLogout,
	AuditEventLogoutFailed,
	AuditEventMe,
	AuditEventMeFailed,
	AuditEventSamplesGet,
	AuditEventSamplesGetFailed,
	AuditEventSamplesGetByID,
	AuditEventSamplesGetByIDFailed,
	AuditEventSamplesCreate,
	AuditEventSamplesCreateFailed,
	AuditEventSamplesUpload,
	AuditEventSamplesUploadFailed,
	AuditEventSamplesUpdate,
	AuditEventSamplesUpdateFailed,
	AuditEventSamplesDelete,
	AuditEventSamplesDeleteFailed,
	AuditEventUsersGet,
	AuditEventUsersGetFailed,
	AuditEventUsersUpdate,
	AuditEventUsersUpdateFailed,
	AuditEventUsersDelete,
	AuditEventUsersDeleteFailed,
	AuditEventUsersUpdatePassword,
	AuditEventUsersUpdatePasswordFailed,
	AuditEventUsersRequestEmailUpdate,
	AuditEventUsersRequestEmailUpdateFailed,
	AuditEventUsersConfirmEmailUpdate,
	AuditEventUsersConfirmEmailUpdateFailed,
	AuditEventAdminAnalysesGet,
	AuditEventAdminAnalysesGetFailed,
	AuditEventAdminAnalysesGetByID,
	AuditEventAdminAnalysesGetByIDFailed,
	AuditEventAdminAnalysesDownloadZip,
	AuditEventAdminAnalysesDownloadZipFailed,
	AuditEventAdminAnalysesCreate,
	AuditEventAdminAnalysesCreateFailed,
	AuditEventAdminAnalysesDownloadBatchTSV,
	AuditEventAdminAnalysesDownloadBatchTSVFailed,
	AuditEventAdminAnalysesUpdate,
	AuditEventAdminAnalysesUpdateFailed,
	AuditEventAdminAnalysesDelete,
	AuditEventAdminAnalysesDeleteFailed,
	AuditEventAdminHealthServicesGet,
	AuditEventAdminHealthServicesGetFailed,
	AuditEventAdminHealthServicesGetByID,
	AuditEventAdminHealthServicesGetByIDFailed,
	AuditEventAdminHealthServicesSearch,
	AuditEventAdminHealthServicesSearchFailed,
	AuditEventAdminHealthServicesCreate,
	AuditEventAdminHealthServicesCreateFailed,
	AuditEventAdminHealthServicesUpdate,
	AuditEventAdminHealthServicesUpdateFailed,
	AuditEventAdminHealthServicesDelete,
	AuditEventAdminHealthServicesDeleteFailed,
	AuditEventAdminCountriesGet,
	AuditEventAdminCountriesGetFailed,
	AuditEventAdminCountriesGetByID,
	AuditEventAdminCountriesGetByIDFailed,
	AuditEventAdminCountriesSearch,
	AuditEventAdminCountriesSearchFailed,
	AuditEventAdminCountriesCreate,
	AuditEventAdminCountriesCreateFailed,
	AuditEventAdminCountriesUpdate,
	AuditEventAdminCountriesUpdateFailed,
	AuditEventAdminCountriesDelete,
	AuditEventAdminCountriesDeleteFailed,
	AuditEventAdminLaboratoriesGet,
	AuditEventAdminLaboratoriesGetFailed,
	AuditEventAdminLaboratoriesGetByID,
	AuditEventAdminLaboratoriesGetByIDFailed,
	AuditEventAdminLaboratoriesSearch,
	AuditEventAdminLaboratoriesSearchFailed,
	AuditEventAdminLaboratoriesCreate,
	AuditEventAdminLaboratoriesCreateFailed,
	AuditEventAdminLaboratoriesUpdate,
	AuditEventAdminLaboratoriesUpdateFailed,
	AuditEventAdminLaboratoriesDelete,
	AuditEventAdminLaboratoriesDeleteFailed,
	AuditEventAdminMetrics,
	AuditEventAdminMetricsFailed,
	AuditEventAdminMicroorganismsGet,
	AuditEventAdminMicroorganismsGetFailed,
	AuditEventAdminMicroorganismsGetByID,
	AuditEventAdminMicroorganismsGetByIDFailed,
	AuditEventAdminMicroorganismsSearch,
	AuditEventAdminMicroorganismsSearchFailed,
	AuditEventAdminMicroorganismsCreate,
	AuditEventAdminMicroorganismsCreateFailed,
	AuditEventAdminMicroorganismsUpdate,
	AuditEventAdminMicroorganismsUpdateFailed,
	AuditEventAdminMicroorganismsDelete,
	AuditEventAdminMicroorganismsDeleteFailed,
	AuditEventAdminOriginsGet,
	AuditEventAdminOriginsGetFailed,
	AuditEventAdminOriginsGetByID,
	AuditEventAdminOriginsGetByIDFailed,
	AuditEventAdminOriginsSearch,
	AuditEventAdminOriginsSearchFailed,
	AuditEventAdminOriginsCreate,
	AuditEventAdminOriginsCreateFailed,
	AuditEventAdminOriginsUpdate,
	AuditEventAdminOriginsUpdateFailed,
	AuditEventAdminOriginsDelete,
	AuditEventAdminOriginsDeleteFailed,
	AuditEventAdminSamplesGet,
	AuditEventAdminSamplesGetFailed,
	AuditEventAdminSamplesGetByID,
	AuditEventAdminSamplesGetByIDFailed,
	AuditEventAdminSamplesCreate,
	AuditEventAdminSamplesCreateFailed,
	AuditEventAdminSamplesUpload,
	AuditEventAdminSamplesUploadFailed,
	AuditEventAdminSamplesUpdate,
	AuditEventAdminSamplesUpdateFailed,
	AuditEventAdminSamplesDelete,
	AuditEventAdminSamplesDeleteFailed,
	AuditEventAdminSampleSourcesGet,
	AuditEventAdminSampleSourcesGetFailed,
	AuditEventAdminSampleSourcesGetByID,
	AuditEventAdminSampleSourcesGetByIDFailed,
	AuditEventAdminSampleSourcesSearch,
	AuditEventAdminSampleSourcesSearchFailed,
	AuditEventAdminSampleSourcesCreate,
	AuditEventAdminSampleSourcesCreateFailed,
	AuditEventAdminSampleSourcesUpdate,
	AuditEventAdminSampleSourcesUpdateFailed,
	AuditEventAdminSampleSourcesDelete,
	AuditEventAdminSampleSourcesDeleteFailed,
	AuditEventAdminSequencersGet,
	AuditEventAdminSequencersGetFailed,
	AuditEventAdminSequencersGetByID,
	AuditEventAdminSequencersGetByIDFailed,
	AuditEventAdminSequencersSearch,
	AuditEventAdminSequencersSearchFailed,
	AuditEventAdminSequencersCreate,
	AuditEventAdminSequencersCreateFailed,
	AuditEventAdminSequencersUpdate,
	AuditEventAdminSequencersUpdateFailed,
	AuditEventAdminSequencersDelete,
	AuditEventAdminSequencersDeleteFailed,
	AuditEventAdminTicketsGet,
	AuditEventAdminTicketsGetFailed,
	AuditEventAdminTicketsGetByID,
	AuditEventAdminTicketsGetByIDFailed,
	AuditEventAdminTicketsAssign,
	AuditEventAdminTicketsAssignFailed,
	AuditEventAdminTicketsResolve,
	AuditEventAdminTicketsResolveFailed,
	AuditEventAdminTicketsDelete,
	AuditEventAdminTicketsDeleteFailed,
	AuditEventAdminUsersGet,
	AuditEventAdminUsersGetFailed,
	AuditEventAdminUsersGetByID,
	AuditEventAdminUsersGetByIDFailed,
	AuditEventAdminUsersCreate,
	AuditEventAdminUsersCreateFailed,
	AuditEventAdminUsersUpdate,
	AuditEventAdminUsersUpdateFailed,
	AuditEventAdminUsersActivate,
	AuditEventAdminUsersActivateFailed,
	AuditEventAdminUsersDeactivate,
	AuditEventAdminUsersDeactivateFailed,
	AuditEventAdminUsersDelete,
	AuditEventAdminUsersDeleteFailed,
}

func AuditEventSelectOptions() []SelectOption {
	opts := make([]SelectOption, len(AuditEvents))
	for i, event := range AuditEvents {
		opts[i] = SelectOption{Label: event, Value: event}
	}
	return opts
}

type Audit struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuidv7();primaryKey"`
	Event     string     `gorm:"not null;index:idx_audit_event_created,priority:1"`
	Source    string     `gorm:"not null"`
	Status    int        `gorm:"not null"`
	Metadata  string     `gorm:"type:jsonb"`
	CreatedAt time.Time  `gorm:"index:idx_audit_event_created,priority:2"`
	UserID    *uuid.UUID `gorm:"type:uuid;index"`
	User      *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL"`
}

type AuditResponse struct {
	ID        uuid.UUID `json:"id"`
	Event     string    `json:"event"`
	Source    string    `json:"source"`
	Status    int       `json:"status"`
	Metadata  string    `json:"metadata"`
	CreatedAt string    `json:"created_at"`
	Username  string    `json:"username"`
}

func (a *Audit) ToResponse() AuditResponse {
	var username string
	if a.User != nil {
		username = a.User.Username
	}

	return AuditResponse{
		ID:        a.ID,
		Event:     a.Event,
		Source:    a.Source,
		Status:    a.Status,
		Metadata:  a.Metadata,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		Username:  username,
	}
}

type AuditFilter struct {
	Event  string     `form:"event"`
	Source string     `form:"source"`
	Status int        `form:"status"`
	Date   *time.Time `form:"date" time_format:"2006-01-02" time_utc:"true"`
	UserID *uuid.UUID `form:"user,parser=encoding.TextUnmarshaler"`
}

type AuditInput struct {
	Event    string     `json:"event"`
	Source   string     `json:"source"`
	Status   int        `json:"status"`
	Metadata *string    `json:"metadata"`
	UserID   *uuid.UUID `json:"user_id"`
}
