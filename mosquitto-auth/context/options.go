package context

type ContextOptions struct {
    DbUri string
    DbName string
    LogLevel string
    TestPassword string
    MonPassword string
    PiotPassword string
}

func NewContextOptions() *ContextOptions {
    o := &ContextOptions{
        DbUri:          "piot",
        DbName:         "piot",
        LogLevel:       "INFO",
        TestPassword:   "",
        MonPassword:    "",
        PiotPassword:   "",
    }
    return o
}

