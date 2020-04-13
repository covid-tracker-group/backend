package config

import "time"

const AppName = "covid-tracker"

// ExpireDailyTracingTokensAfter states after how much time generated
// Daily Tracing Authentication Tokens should expire
const ExpireDailyTracingTokensAfter = time.Hour * 24 * 7
