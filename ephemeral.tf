ephemeral "aws_kms_secrets" "secret" {
  secret {
    name    = "secret"
    payload = "AQICAHhY4h6oV67xgokHQSgppqIoTXTPsMu3hBoLnmg4QWonqAGMzsK4nyoN+mee8rGu5Lm+AAAAZDBiBgkqhkiG9w0BBwagVTBTAgEAME4GCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQME6JCDiFXmgQ1KgRHAgEQgCGdTNbFT2fgotjzFi65w5TRmvP8zek9tZ9NFl61hGbUALg="
    key_id  = "alias/default"
  }
}
