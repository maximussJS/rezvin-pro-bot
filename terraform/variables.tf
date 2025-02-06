variable "hcloud_token" {
  description = "Hetzner Cloud API Token"
  type        = string
  sensitive   = true
}

variable "hcloud_ssh_key_fingerprint" {
  description = "SSH key fingerprint to use for the server"
  type        = string
  sensitive   = true
}

variable "server_name" {
  description = "Server Name"
  type        = string
  default     = "rezvin"
}

variable "private_key_path" {
  description = "Path to the private SSH key for connecting to the server"
  type        = string
  default     = "~/.ssh/hetzner_id_rsa"
}

variable "private_key_content" {
  description = "Content of the private SSH key for connecting to the server"
  type        = string
  sensitive   = true
  default     = ""
}

variable "app_env" {
  description = "The application environment (e.g., development, production)"
  type        = string
  default     = "production"
}

variable "bots" {
  description = "Map of bot configurations"
  type = map(object({
    http_port = string
    container_name = string
    bot_token     = string
    alert_chat_id = string
    postgres_schema  = string
    webhook_secret_token = string
    image = string
    admin_name = string
  }))
}

variable "postgres_dsn" {
  description = "PostgreSQL Data Source Name containing connection details"
  type        = string
  sensitive   = true
}

variable "run_migrations" {
  description = "Flag to determine whether to run database migrations"
  type        = bool
  default     = true
}

variable "request_timeout_in_seconds" {
  description = "Timeout duration for requests in seconds"
  type        = number
  default     = 60
}

variable "ssl_cert_path" {
    description = "Path to the SSL certificate"
    type        = string
    default     = "/app/certs/cert.pem"
}

variable "ssl_key_path" {
    description = "Path to the SSL private key"
    type        = string
    default     = "/app/certs/priv.pem"
}
