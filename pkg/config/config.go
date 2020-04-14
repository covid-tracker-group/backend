package config

import "time"

const AppName = "covid-tracker"

// ExpireDailyTracingTokensAfter states after how much time generated
// Daily Tracing Authentication Tokens must expire.
const ExpireDailyTracingTokensAfter = time.Hour * 24 * 7

// ExpireHealthAuthorisationTokensAfter status after how much time
// health authorisation tokens must expire.
const ExpireHealthAuthorisationTokensAfter = time.Hour * 24 * 2
