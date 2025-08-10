package qbitorrent

// all, downloading, seeding, completed, paused, active, inactive, resumed, stalled, stalled_uploading, stalled_downloading, errored
const (
	TAG_ALL                 = "all"
	TAG_DOWNLOADING         = "downloading"
	TAG_SEEDING             = "seeding"
	TAG_COMPLETED           = "completed"
	TAG_PAUSED              = "paused"
	TAG_ACTIVE              = "active"
	TAG_INACTIVE            = "inactive"
	TAG_RESUMED             = "resumed"
	TAG_STALLED             = "stalled"
	TAG_STALLED_UPLOADING   = "stalled_uploading"
	TAG_STALLED_DOWNLOADING = "stalled_downloading"
	TAG_ERRORED             = "errored"
)

const (
	STATE_ERROR                = "error"              // Some error occurred, applies to paused torrents
	STATE_MISSING_FILES        = "missingFiles"       // Torrent data files is missing
	STATE_UPLOADING            = "uploading"          // Torrent is being seeded and data is being transferred
	STATE_PAUSED_UP            = "pausedUP"           // Torrent is paused and has finished downloading
	STATE_QUEUED_UP            = "queuedUP"           // Queuing is enabled and torrent is queued for upload
	STATE_STALLED_UP           = "stalledUP"          // Torrent is being seeded, but no connection were made
	STATE_CHECKING_UP          = "checkingUP"         // Torrent has finished downloading and is being checked
	STATE_FORCED_UP            = "forcedUP"           // Torrent is forced to uploading and ignore queue limit
	STATE_ALLOCATING           = "allocating"         // Torrent is allocating disk space for download
	STATE_DOWNLOADING          = "downloading"        // Torrent is being downloaded and data is being transferred
	STATE_META_DL              = "metaDL"             // Torrent has just started downloading and is fetching metadata
	STATE_PAUSED_DL            = "pausedDL"           // Torrent is paused and has NOT finished downloading
	STATE_QUEUED_DL            = "queuedDL"           // Queuing is enabled and torrent is queued for download
	STATE_STALLED_DL           = "stalledDL"          // Torrent is being downloaded, but no connection were made
	STATE_CHECKING_DL          = "checkingDL"         // Same as checkingUP, but torrent has NOT finished downloading
	STATE_FORCED_DL            = "forcedDL"           // Torrent is forced to downloading to ignore queue limit
	STATE_CHECKING_RESUME_DATA = "checkingResumeData" // Checking resume data on qBt startup
	STATE_MOVING               = "moving"             // Torrent is moving to another location
	STATE_UNKNOWN              = "unknown"            // Unknown status
)
