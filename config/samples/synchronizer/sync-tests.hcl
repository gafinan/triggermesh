bridge "ibm_mq_sync_backend" {}

// ---- Event Producer ----

source github "github_repo" {
  event_types = [
    "issues",
    "issue_comment"
  ]
  owner_and_repository = "tzununbekov/tm"
  tokens = secret_name("github-secret")

  to = target.mq
}

// ---- Sync API ----

target synchronizer "mq" {
  request_key = "id[0:23]"
  response_correlation_key = "extensions.correlationid[0:23]"
  response_wait_timeout = 180
  // timeout_response_policy = "error"

  to = target.input_channel
}

// ---- IBM MQ Input target ----

target ibmmq "input_channel" {
  connection_name = "ibm-mq.tzununbekov.svc.cluster.local(1414)"
  credentials = secret_name("ibm-mq-secret")
  queue_manager = "QM1"
  queue_name = "DEV.QUEUE.1"
  channel_name = "DEV.APP.SVRCONN"
  reply_to {
    queue = "DEV.QUEUE.2"
  }

  discard_ce_context = true
}

// ---- IBM MQ Input source ----

target container "sync-replier" {
  image = "docker.io/tzununbekov/ibmmq-sync-replier"
  public = false
}

// ---- IBM MQ Output source ----

source ibmmq "output_channel" {
  connection_name = "ibm-mq.tzununbekov.svc.cluster.local(1414)"
  credentials = secret_name("ibm-mq-secret")
  queue_manager = "QM1"
  queue_name = "DEV.QUEUE.2"
  channel_name = "DEV.APP.SVRCONN"

  to = target.mq
}

// ---- Event Receiver ----

target container "sockeye" {
  image = "docker.io/n3wscott/sockeye:v0.7.0"
  public = true
}

